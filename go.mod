module github.com/go-openapi/swag

retract v0.24.0 // bad tagging of the main module: superseeded by v0.24.1

require (
	github.com/go-openapi/swag/cmdutils v0.25.4
	github.com/go-openapi/swag/conv v0.25.4
	github.com/go-openapi/swag/fileutils v0.25.4
	github.com/go-openapi/swag/jsonname v0.25.4
	github.com/go-openapi/swag/jsonutils v0.25.4
	github.com/go-openapi/swag/loading v0.25.4
	github.com/go-openapi/swag/mangling v0.25.4
	github.com/go-openapi/swag/netutils v0.25.4
	github.com/go-openapi/swag/stringutils v0.25.4
	github.com/go-openapi/swag/typeutils v0.25.4
	github.com/go-openapi/swag/yamlutils v0.25.4
	github.com/go-openapi/testify/v2 v2.2.0
)

require go.yaml.in/yaml/v3 v3.0.4 // indirect

replace (
	github.com/go-openapi/swag/cmdutils => ./cmdutils
	github.com/go-openapi/swag/conv => ./conv
	github.com/go-openapi/swag/fileutils => ./fileutils
	github.com/go-openapi/swag/jsonname => ./jsonname
	github.com/go-openapi/swag/jsonutils => ./jsonutils
	github.com/go-openapi/swag/jsonutils/fixtures_test => ./jsonutils/fixtures_test
	github.com/go-openapi/swag/loading => ./loading
	github.com/go-openapi/swag/mangling => ./mangling
	github.com/go-openapi/swag/netutils => ./netutils
	github.com/go-openapi/swag/stringutils => ./stringutils
	github.com/go-openapi/swag/typeutils => ./typeutils
	github.com/go-openapi/swag/yamlutils => ./yamlutils
)

go 1.24.0
