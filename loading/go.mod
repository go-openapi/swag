module github.com/go-openapi/swag/loading

require (
	github.com/go-openapi/swag/yamlutils v0.25.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/swag/conv v0.25.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.0 // indirect
	github.com/go-openapi/swag/typeutils v0.25.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/go-openapi/swag/conv => ../conv
	github.com/go-openapi/swag/jsonutils => ../jsonutils
	github.com/go-openapi/swag/jsonutils/fixtures_test => ../jsonutils/fixtures_test
	github.com/go-openapi/swag/typeutils => ../typeutils
	github.com/go-openapi/swag/yamlutils => ../yamlutils
)

go 1.24.0
