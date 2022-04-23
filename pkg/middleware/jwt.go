package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/store"
)

const (
	// DefaultAuthScheme default authentication scheme for JWT tokens
	DefaultAuthScheme = "Bearer"
	// DefaultAuthHeaderName default header to load the JWT
	DefaultAuthHeaderName = "Authorization"

	// DefaultContextVar the variable in the context to store the JWT token after successful login
	DefaultContextVar = "user"
)

var (
	// ErrJWTMissing missing or malformed jwt
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")

	// ErrJWTValidation invalid jwt
	ErrJWTValidation = echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt")
)

// JWTConfig jwt middleware configuration
type JWTConfig struct {
	ProviderURL string
	ClientID    string
	AuthScheme  string
}

// JWTWithConfig middleware which validates tokens
func JWTWithConfig(config *JWTConfig, authStore *store.AuthStore) echo.MiddlewareFunc {

	if config.AuthScheme == "" {
		config.AuthScheme = DefaultAuthScheme
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Info().Msg("context updating")

			tokenString, err := extractFromHeader(c, "Authorization", "Bearer")
			if err != nil {
				//no auth found, see what we can do without one
				log.Info().Msg("no auth found")
				return next(c)
			}
			//try to get it from stored auth
			authUser, err := authStore.GetAuthUser(tokenString)
			if err == nil {
				//we have a valid authUser; returning
				log.Info().Msg("got auth user from cache")
				c.Set("authUser", authUser)
				return next(c)
			}

			token, err := jwt.ParseWithClaims(tokenString, &authorization.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(jwt.SigningMethod); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})

			if !token.Valid || err != nil {
				//no valid auth found, see what we can do without one
				log.Info().Msg("no valid auth found")
				return next(c)
			}

			authClaims := token.Claims.(*authorization.AuthClaims)
			authUser, err = authorization.NewAuthUserWithClaims(authClaims)
			if err != nil {
				return next(c)
			}
			authStore.AddAuthUser(tokenString, authUser)

			c.Set("authUser", authUser)
			return next(c)
		}
	}
}

// ExtractFromHeader attempt to get the JWT from the provided header
func extractFromHeader(c echo.Context, header, authScheme string) (string, error) {
	auth := c.Request().Header.Get(header)
	l := len(authScheme)
	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}
	return "", ErrJWTMissing
}
