package store

import (
	"reflect"
	"testing"

	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewUserStore(t *testing.T) {
	type args struct {
		dbconn *service.MongoService
	}
	tests := []struct {
		name string
		args args
		want *UserStore
	}{
		{
			name: "InstanciateUserStore",
			args: args{
				dbconn: &service.MongoService{},
			},
			want: &UserStore{
				dbconn:   &service.MongoService{},
				collName: "user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserStore(tt.args.dbconn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserStore_GetWithEmail(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		email string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults []*document.User
		wantErr     bool
	}{
		{
			name:   "TestRandomString",
			fields: fields{},
			args: args{
				email: "askadsldkasdk999",
			},
			wantResults: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			gotResults, err := u.GetWithEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetWithEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("UserStore.GetWithEmail() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestUserStore_GetWithId(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults []*document.User
		wantErr     bool
	}{
		{
			name:   "TestRandomString",
			fields: fields{},
			args: args{
				id: "askadsldkasdk999",
			},
			wantResults: nil,
			wantErr:     true,
		},
		{
			name:   "TestValidHex",
			fields: fields{},
			args: args{
				id: "6262ce0dafd1acb9dfbc4f87",
			},
			wantResults: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			gotResults, err := u.GetWithId(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetWithId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("UserStore.GetWithId() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestUserStore_getOne(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		query  bson.M
		result interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "MongoNotInstanciated",
			fields:  fields{},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			if err := u.getOne(tt.args.query, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("UserStore.getOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserStore_getAll(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		query   bson.M
		results interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "MongoNotInstanciated",
			fields:  fields{},
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			if err := u.getAll(tt.args.query, tt.args.results); (err != nil) != tt.wantErr {
				t.Errorf("UserStore.getAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserStore_setStringId(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		result *document.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "EmptyUserDocument",
			fields: fields{
				dbconn:   &service.MongoService{},
				collName: "",
			},
			args: args{
				result: &document.User{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			u.setStringId(tt.args.result)
		})
	}
}
