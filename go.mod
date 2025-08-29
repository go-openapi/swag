module github.com/go-openapi/swag

require (
	github.com/go-openapi/swag/cmdutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/conv v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/fileutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/jsonname v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/jsonutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/loading v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/mangling v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/netutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/stringutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/typeutils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/swag/yamlutils v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-openapi/swag/cmdutils => ./cmdutils

replace github.com/go-openapi/swag/conv => ./conv

replace github.com/go-openapi/swag/fileutils => ./fileutils

replace github.com/go-openapi/swag/jsonname => ./jsonname

replace github.com/go-openapi/swag/jsonutils => ./jsonutils

replace github.com/go-openapi/swag/loading => ./loading

replace github.com/go-openapi/swag/mangling => ./mangling

replace github.com/go-openapi/swag/netutils => ./netutils

replace github.com/go-openapi/swag/stringutils => ./stringutils

replace github.com/go-openapi/swag/typeutils => ./typeutils

replace github.com/go-openapi/swag/yamlutils => ./yamlutils

go 1.20.0
