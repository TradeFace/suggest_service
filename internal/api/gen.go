// Package api contains the REST interfaces.
package api

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -generate client,types,server,spec -package api -o suggest.gen.go suggest.yml
//go:generate gofmt -s -w suggest.gen.go
/*

docker run --rm \
    -v $PWD:/local openapitools/openapi-generator-cli generate \
    -i /local/suggest.yml \
    -g go-server \
    -o /local/suggest.gotest
*/
