package authorization

import (
	"errors"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tradeface/suggest_service/pkg/helpers"
)

type AuthChecker struct {
	user map[string]*AuthUser
}

type AuthUser struct {
	claims AuthClaims
	roles  *helpers.Set
	expire time.Time
}

type AuthClaims struct {
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	CompanyId string   `json:"companyid"`
	Roles     []string `json:"roles"`
	jwt.StandardClaims
}

func NewAuthChecker() *AuthChecker {
	return &AuthChecker{
		user: make(map[string]*AuthUser, 0),
	}
}

func (ac *AuthChecker) GetAuthUser(userId string) (user *AuthUser, err error) {

	if user, ok := ac.user[userId]; ok {
		return user, nil
	}
	return user, errors.New("user not found")
}

func (ac *AuthChecker) SetTokenClaims(claims *AuthClaims) (userClaims *AuthUser, err error) {

	roles, err := ac.getRolesSet(*claims)
	if err != nil {
		return userClaims, err
	}

	userClaims = &AuthUser{
		claims: *claims,
		roles:  roles,
	}
	// userId, err := userClaims.GetClaim("Id")
	// if err != nil {
	// 	return userClaims, err
	// }

	//ac.user[userId.(string)] = userClaims
	return userClaims, nil
}

func (ac *AuthChecker) getRolesSet(claims AuthClaims) (roles *helpers.Set, err error) {

	if claims.Roles != nil {
		roles = helpers.NewSet(claims.Roles)
	}
	return roles, nil
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

//TODO: add function to check authorization
//TODO: add function to  get other claims (like companyid)
