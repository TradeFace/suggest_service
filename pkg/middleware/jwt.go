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

func JWTWithConfig(stores *store.Provider, JWTSalt string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			tokenString, err := extractFromHeader(c, "Authorization", "Bearer")
			if err != nil {
				//no auth found, see what we can do without one
				log.Info().Msg("no auth found")
				return next(c)
			}

			//try to get it from stored auth
			authUser, err := stores.Auth.GetAuthUser(tokenString)
			if err == nil {
				//we have a valid authUser; returning
				log.Info().Msg("got auth user from cache")
				c.Set("authUser", authUser)
				return next(c)
			}

			//try to get auth by token
			token, err := ParseWithClaims(tokenString, JWTSalt, stores, &authorization.AuthClaims{})

			if !token.Valid || err != nil {
				//no valid auth found, see what we can do without one
				log.Info().Msg("no valid auth found")
				return next(c)
			}

			//we have a valid token, try to create authuser
			authClaims := token.Claims.(*authorization.AuthClaims)
			authUser, err = authorization.NewAuthUserWithClaims(authClaims)
			if err != nil {
				//failed to create authuser, see what we can do without one
				return next(c)
			}
			//we have a valid token, add it to cache
			stores.Auth.AddAuthUser(tokenString, authUser)

			c.Set("authUser", authUser)
			return next(c)
		}
	}
}

func ParseWithClaims(tokenString string, JWTSalt string, stores *store.Provider, authClaims *authorization.AuthClaims) (*jwt.Token, error) {

	return jwt.ParseWithClaims(tokenString, authClaims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(jwt.SigningMethod); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		expiresAt := token.Claims.(*authorization.AuthClaims).ExpiresAt
		if expiresAt == 0 {
			return nil, fmt.Errorf("expire claim missing")
		}

		Id := token.Claims.(*authorization.AuthClaims).Id
		if Id == "" {
			return nil, fmt.Errorf("id claim missing")
		}

		docUser, err := stores.User.GetWithId(Id)
		if err != nil {
			return nil, err
		}

		signingToken := authorization.GetSigningToken(JWTSalt, docUser[0].Password, expiresAt)

		return signingToken, nil
	})
}

// ExtractFromHeader attempt to get the JWT from the provided header
func extractFromHeader(c echo.Context, header string, authScheme string) (string, error) {

	auth := c.Request().Header.Get(header)
	l := len(authScheme)
	if len(auth) > l+1 && auth[:l] == authScheme {
		return auth[l+1:], nil
	}
	return "", ErrJWTMissing
}
