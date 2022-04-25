package store

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type (
	mockMongoService struct {
	}
)

func (mmc *mockMongoService) GetAll(coll string, query bson.M, results interface{}) (err error) {

	r := results.(*[]*document.Domain)
	z := *r

	z = append(z, &document.Domain{

		ObjectID:           [12]byte{},
		Id:                 "",
		Active:             false,
		Catalogs:           []*string{},
		Host:               "testhost",
		MainClassification: "TC000001",
		Settings:           map[string]interface{}{},
		Modules:            []string{},
	})
	/*	results = make([]document.Domain, 0)
		result := document.Domain{
			ObjectID:           [12]byte{},
			Id:                 "",
			Active:             false,
			Catalogs:           []*string{},
			Host:               "testhost",
			MainClassification: "TC000001",
			Settings:           map[string]interface{}{},
			Modules:            []string{},
		}
		results = append(results.([]document.Domain), result)*/
	results = z
	fmt.Println("n mock getall ", query, results, z)
	return nil
}
func (mmc *mockMongoService) GetOne(coll string, query bson.M, result interface{}) error {
	r := result.(**document.Domain)
	z := **r
	z.Active = true //
	z.Host = "kaas"
	// result = &document.Domain{
	// 	ObjectID:           [12]byte{},
	// 	Id:                 "",
	// 	Active:             false,
	// 	Catalogs:           []*string{},
	// 	Host:               "testhost",
	// 	MainClassification: "TC000001",
	// 	Settings:           map[string]interface{}{},
	// 	Modules:            []string{},
	// }
	return nil
}
func (mmc *mockMongoService) GetMongoId(id string) (objID primitive.ObjectID, err error) {
	return primitive.ObjectIDFromHex(id)
}

func TestDomainStore_GetWithId(t *testing.T) {
	type fields struct {
		dbconn   MongoServiceInterface
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
			name: "TestRandomString",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args: args{
				id: "askadsldkasdk999",
			},
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "TestValidHex",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args: args{
				id: "6262ce0dafd1acb9dfbc4f87",
			},
			wantResults: []*document.Domain{
				0: {
					ObjectID:           [12]byte{},
					Id:                 "",
					Active:             false,
					Catalogs:           []*string{},
					Host:               "testhost",
					MainClassification: "TC000001",
					Settings:           map[string]interface{}{},
					Modules:            []string{},
				},
			},
			wantErr: false,
		},
		// {
		// 	name: "TestWithInterface",
		// 	fields: fields{
		// 		dbconn:   &mockMongoService{},
		// 		collName: "domain",
		// 	},
		// 	args: args{
		// 		id: "6262ce0dafd1acb9dfbc4f87",
		// 	},
		// 	wantResults: []*document.Domain{},
		// 	wantErr:     true,
		// },
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
		dbconn   MongoServiceInterface
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
		// {
		// 	name: "TestRandomString",
		// 	fields: fields{
		// 		dbconn:   &mockMongoService{},
		// 		collName: "",
		// 	},
		// 	args: args{
		// 		host: "nonExistingHost",
		// 	},
		// 	wantResults: nil,
		// 	wantErr:     true,
		// },
		// {
		// 	name:        "",
		// 	fields:      fields{
		// 		dbconn:   &mockMongoService{},
		// 		collName: "",
		// 	},
		// 	args:        args{},
		// 	wantResults: []*document.Domain{},
		// 	wantErr:     false,
		// },
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
		dbconn   MongoServiceInterface
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
		// {
		// 	name:   "TestRandomString",
		// 	fields: fields{},
		// 	args: args{
		// 		host: "askadsldkasdk999",
		// 	},
		// 	wantResult: nil,
		// 	wantErr:    true,
		// },
		{
			name: "",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args: args{
				host: "testHost",
			},
			wantResult: &document.Domain{},
			wantErr:    false,
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

func TestDomainStore_GetOne(t *testing.T) {
	type fields struct {
		dbconn   MongoServiceInterface
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
			name: "MongoNotInstanciated",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			if err := d.GetOne(tt.args.query, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.GetOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomainStore_GetAll(t *testing.T) {
	type fields struct {
		dbconn   MongoServiceInterface
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
			name: "MongoNotInstanciated",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			if err := d.GetAll(tt.args.query, tt.args.results); (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomainStore_getStringId(t *testing.T) {
	type fields struct {
		dbconn   MongoServiceInterface
		collName string
	}
	type args struct {
		result *document.Domain
	}
	objID, _ := primitive.ObjectIDFromHex("625a78c4afd1acb9dfbc4f86")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "EmptyInput",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args:    args{},
			want:    "",
			wantErr: true,
		},
		{
			name: "EmptyInput",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args: args{
				result: &document.Domain{
					ObjectID:           [12]byte{},
					Id:                 "",
					Active:             false,
					Catalogs:           []*string{},
					Host:               "",
					MainClassification: "",
					Settings:           map[string]interface{}{},
					Modules:            []string{},
				},
			},
			want:    "000000000000000000000000",
			wantErr: false,
		},
		{
			name: "",
			fields: fields{
				dbconn:   &mockMongoService{},
				collName: "",
			},
			args: args{
				result: &document.Domain{
					ObjectID:           objID,
					Id:                 "",
					Active:             false,
					Catalogs:           []*string{},
					Host:               "",
					MainClassification: "",
					Settings:           map[string]interface{}{},
					Modules:            []string{},
				},
			},
			want:    "625a78c4afd1acb9dfbc4f86",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &DomainStore{
				dbconn:   tt.fields.dbconn,
				collName: tt.fields.collName,
			}
			got, err := u.getStringId(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomainStore.getStringId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DomainStore.getStringId() = %v, want %v", got, tt.want)
			}
		})
	}
}
