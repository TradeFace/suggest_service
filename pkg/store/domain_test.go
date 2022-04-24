package store

import (
	"reflect"
	"testing"

	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewDomainStore(t *testing.T) {
	type args struct {
		dbconn *service.MongoService
	}
	tests := []struct {
		name string
		args args
		want *DomainStore
	}{
		{
			name: "InstanciateDomainStore",
			args: args{
				dbconn: &service.MongoService{},
			},
			want: &DomainStore{
				dbconn:   &service.MongoService{},
				collName: "domain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDomainStore(tt.args.dbconn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDomainStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainStore_GetWithId(t *testing.T) {
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
		wantResults []*document.Domain
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
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			gotResults, err := d.GetWithId(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.GetWithId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("DomainStore.GetWithId() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestDomainStore_GetWithHost(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		host string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults []*document.Domain
		wantErr     bool
	}{
		{
			name:   "TestRandomString",
			fields: fields{},
			args: args{
				host: "askadsldkasdk999",
			},
			wantResults: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			gotResults, err := d.GetWithHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.GetWithHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("DomainStore.GetWithHost() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestDomainStore_GetOneWithHost(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		host string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult *document.Domain
		wantErr    bool
	}{
		{
			name:   "TestRandomString",
			fields: fields{},
			args: args{
				host: "askadsldkasdk999",
			},
			wantResult: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			gotResult, err := d.GetOneWithHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.GetOneWithHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DomainStore.GetOneWithHost() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestDomainStore_getOne(t *testing.T) {
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
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			if err := d.getOne(tt.args.query, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.getOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomainStore_getAll(t *testing.T) {
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
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			if err := d.getAll(tt.args.query, tt.args.results); (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.getAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomainStore_setStringId(t *testing.T) {
	type fields struct {
		dbconn   *service.MongoService
		collName string
	}
	type args struct {
		result *document.Domain
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			d.setStringId(tt.args.result)
		})
	}
}
