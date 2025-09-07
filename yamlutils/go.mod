module github.com/go-openapi/swag/yamlutils

require (
	github.com/go-openapi/swag/conv v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/jsonutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/jsonutils/fixtures_test v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/typeutils v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

replace (
	github.com/go-openapi/swag/conv => ../conv
	github.com/go-openapi/swag/jsonutils => ../jsonutils
	github.com/go-openapi/swag/jsonutils/fixtures_test => ../jsonutils/fixtures_test
	github.com/go-openapi/swag/typeutils => ../typeutils
)

go 1.24.0
