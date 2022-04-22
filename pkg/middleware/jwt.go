package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/pkg/authorization"
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

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// JWTWithConfig middleware which validates tokens
func JWTWithConfig(config *JWTConfig, auth *authorization.AuthChecker) echo.MiddlewareFunc {

	if config.AuthScheme == "" {
		config.AuthScheme = DefaultAuthScheme
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Info().Msg("context updating")

			tokenString, err := extractFromHeader(c, "Authorization", config.AuthScheme)
			if err != nil {
				//no auth found, see what we can do without one
				log.Info().Msg("no auth found")
				return next(c)
			}

			token, err := jwt.ParseWithClaims(tokenString, &authorization.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(jwt.SigningMethod); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})

			if token.Valid != true {
				//no valid auth found, see what we can do without one
				log.Info().Msg("no valid auth found")
				return next(c)
			}

			claims := token.Claims.(*authorization.AuthClaims)
			auth.SetTokenClaims(claims)

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

// ValidateToken and return an error if it fails
// func validateToken(ctx context.Context, providerURL, token string) (*jwt.JwtPayload, error) {
// 	payload, err := jwt.Validate(ctx, providerURL, token)
// 	if err != nil {
// 		log.Warn().Err(err).Msg("failed to validate header")
// 		return nil, ErrJWTValidation
// 	}
// 	return payload, nil
// }
