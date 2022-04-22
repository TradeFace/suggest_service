package authorization

import (
	"github.com/golang-jwt/jwt"
)

type AuthChecker struct {
	claims AuthClaims
}
type AuthClaims struct {
	Name      string `json:"name"`
	CompanyId string `json:"companyid"`
	Admin     bool   `json:"admin"`
	jwt.StandardClaims
}

func (ac *AuthChecker) GetAuthClaims() *AuthClaims {
	return &AuthClaims{}
}

func (ac *AuthChecker) SetTokenClaims(claims *AuthClaims) {
	ac.claims = *claims
}

//TODO: add function to check authorization
//TODO: add function to  get other claims (like companyid)
