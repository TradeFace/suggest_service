package store

import (
	"reflect"
	"testing"

	"github.com/tradeface/jwt_service/pkg/authorization"
)

func TestNewAuthStore(t *testing.T) {
	tests := []struct {
		name string
		want *AuthStore
	}{
		{
			name: "InstanciateAuthStore",
			want: &AuthStore{
				authUser: map[string]*authorization.AuthUser{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthStore_GetAuthUser(t *testing.T) {
	type fields struct {
		authUser map[string]*authorization.AuthUser
	}
	type args struct {
		tokenString string
	}
	// lastSeen := time.Now().Add(15 * time.Minute)
	// expire := time.Now().Add(time.Hour)
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAuthUser *authorization.AuthUser
		wantErr      bool
	}{
		{
			name: "NoneMatchingToken",
			fields: fields{
				authUser: map[string]*authorization.AuthUser{},
			},
			args: args{
				tokenString: "abcd",
			},
			wantAuthUser: nil,
			wantErr:      true,
		},
		/*
			can't test, expire gets updated on function call
			{
				name: "MatchingToken",
				fields: fields{
					authUser: map[string]*authorization.AuthUser{
						"abcd": {
							TokenExpire:    expire,
							LastseenExpire: lastSeen,
						},
					},
				},
				args: args{
					tokenString: "abcd",
				},
				wantAuthUser: &authorization.AuthUser{
					TokenExpire:    expire,
					LastseenExpire: lastSeen,
				},
				wantErr: false,
			},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthStore{
				authUser: tt.fields.authUser,
			}
			gotAuthUser, err := ac.GetAuthUser(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthStore.GetAuthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAuthUser, tt.wantAuthUser) {
				t.Errorf("AuthStore.GetAuthUser() = %v, want %v", gotAuthUser, tt.wantAuthUser)
			}
		})
	}
}

func TestAuthStore_AddAuthUser(t *testing.T) {
	type fields struct {
		authUser map[string]*authorization.AuthUser
	}
	type args struct {
		tokenString string
		authUser    *authorization.AuthUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "AuthStoreNotInstancated",
			fields: fields{},
			args: args{
				tokenString: "abcd",
				authUser:    &authorization.AuthUser{},
			},
			wantErr: true,
		},
		/*
			can't test AuthUser.claim is private
			{
				name:   "AuthClaimWithoutExpireClaim",
				fields: fields{},
				args: args{
					tokenString: "abcd",
					authUser: &authorization.AuthUser{
						claims: &authorization.AuthClaims{
							Name:           "",
							Email:          "",
							CompanyId:      "",
							Roles:          []string{},
							StandardClaims: jwt.StandardClaims{},
						},
						TokenExpire:    time.Time{},
						LastseenExpire: time.Time{},
					},
				},
				wantErr: true,
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthStore{
				authUser: tt.fields.authUser,
			}
			if err := ac.AddAuthUser(tt.args.tokenString, tt.args.authUser); (err != nil) != tt.wantErr {
				t.Errorf("AuthStore.AddAuthUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthStore_forever(t *testing.T) {
	type fields struct {
		authUser map[string]*authorization.AuthUser
	}
	type args struct {
		done <-chan bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// How to test expires?
		{
			name: "RunForever",
			fields: fields{
				authUser: map[string]*authorization.AuthUser{},
			},
			args: args{
				done: make(<-chan bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AuthStore{
				authUser: tt.fields.authUser,
			}
			ac.forever(tt.args.done)
		})
	}
}
