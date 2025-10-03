module github.com/go-openapi/swag/loading

require (
	github.com/go-openapi/swag/yamlutils v0.25.1
	github.com/go-openapi/testify/enable/yaml/v2 v2.0.2
	github.com/go-openapi/testify/v2 v2.0.2
)

require (
	github.com/go-openapi/swag/conv v0.25.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.1 // indirect
	github.com/go-openapi/swag/typeutils v0.25.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace (
	github.com/go-openapi/swag/conv => ../conv
	github.com/go-openapi/swag/jsonutils => ../jsonutils
	github.com/go-openapi/swag/jsonutils/fixtures_test => ../jsonutils/fixtures_test
	github.com/go-openapi/swag/typeutils => ../typeutils
	github.com/go-openapi/swag/yamlutils => ../yamlutils
)

go 1.24.0
