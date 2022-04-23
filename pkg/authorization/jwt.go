package authorization

import (
	"errors"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/helpers"
)

type AuthUser struct {
	claims         AuthClaims
	roles          *helpers.Set
	TokenExpire    time.Time
	LastseenExpire time.Time
}

type AuthClaims struct {
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	CompanyId string   `json:"companyid"`
	Roles     []string `json:"roles"`
	jwt.StandardClaims
}

const TOKEN_VALID_MIN = 60

func NewAuthUserWithClaims(authClaims *AuthClaims) (authUser *AuthUser, err error) {

	roles, err := getRolesSet(*authClaims)
	if err != nil {
		return authUser, err
	}

	authUser = &AuthUser{
		claims: *authClaims,
		roles:  roles,
	}
	return authUser, nil
}

func (au *AuthUser) HasRole(role string) bool {

	return au.roles.Contains(role)
}

func (au *AuthUser) GetClaim(claim string) (interface{}, error) {

	v := reflect.ValueOf(au.claims)
	if _, ok := v.Type().FieldByName(claim); ok {
		return v.FieldByName(claim).Interface(), nil
	}
	return nil, errors.New("not a claim")
}

func NewJwtWithUser(docUser *document.User) (token string, err error) {

	// Set custom claims
	authClaims := &AuthClaims{
		docUser.Name,
		docUser.Email,
		docUser.CompanyId,
		docUser.Roles,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * TOKEN_VALID_MIN).Unix(),
			Id:        docUser.Id,
		},
	}

	// Create token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)

	// Generate encoded token and send it as response.
	return jwtToken.SignedString([]byte("secret"))
}

// Private helpers
func getRolesSet(authClaims AuthClaims) (roles *helpers.Set, err error) {

	if authClaims.Roles != nil {
		roles = helpers.NewSet(authClaims.Roles)
	}
	return roles, nil
}
