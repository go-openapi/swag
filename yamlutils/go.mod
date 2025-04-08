module github.com/go-openapi/swag/yamlutils

require (
	github.com/go-openapi/swag/jsonutils v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

replace github.com/go-openapi/swag/jsonutils => ../jsonutils

go 1.20
