package authorization

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tradeface/suggest_service/pkg/document"
)

func MakeJwt(usr *document.User) error {

	// Set custom claims
	claims := &AuthClaims{
		usr.Name,
		usr.Email,
		usr.CompanyId,
		usr.Roles,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5000).Unix(),
			Id:        usr.Id,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	usr.Token = t
	return nil
}
