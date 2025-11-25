module github.com/go-openapi/swag/jsonutils

require (
	github.com/go-openapi/swag/conv v0.25.4
	github.com/go-openapi/swag/jsonutils/fixtures_test v0.25.4
	github.com/go-openapi/swag/typeutils v0.25.4
	github.com/go-openapi/testify/v2 v2.0.2
)

require (
	github.com/go-openapi/testify/enable/yaml/v2 v2.0.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
)

replace (
	github.com/go-openapi/swag/conv => ../conv
	github.com/go-openapi/swag/jsonutils/fixtures_test => ./fixtures_test
	github.com/go-openapi/swag/typeutils => ../typeutils
)

go 1.24.0
