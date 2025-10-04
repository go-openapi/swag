// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestManglingIface(t *testing.T) {
	const (
		sample               = "hello_swagger"
		sampleWithInitialism = "Http_server"
	)

	t.Run("deprecated name mangling functions should work", func(t *testing.T) {
		assert.Equal(t, "HelloSwagger", ToGoName(sample))
		assert.Equal(t, "HTTPServer", ToGoName(sampleWithInitialism))
		assert.Equal(t, "helloSwagger", ToVarName(sample))
		assert.Equal(t, "httpServer", ToVarName(sampleWithInitialism))
		assert.Equal(t, "hello_swagger", ToFileName(sample))
		assert.Equal(t, "http_server", ToFileName(sampleWithInitialism))
		assert.Equal(t, "hello-swagger", ToCommandName(sample))
		assert.Equal(t, "http-server", ToCommandName(sampleWithInitialism))
		assert.Equal(t, "hello swagger", ToHumanNameLower(sample))
		assert.Equal(t, "Http server", ToHumanNameLower(sampleWithInitialism))
		assert.Equal(t, "Hello Swagger", ToHumanNameTitle(sample))
		assert.Equal(t, "Http Server", ToHumanNameTitle(sampleWithInitialism))

		assert.Equal(t, "Swagger", Camelize("SWAGGER"))
		assert.Equal(t, "helloSwagger", ToJSONName(sample))

		t.Run("with global config", func(t *testing.T) {
			AddInitialisms("ELB")                                   // adding non-default initialism
			GoNamePrefixFunc = func(_ string) string { return "Z" } // adding non-default prefix function

			assert.Equal(t, "Z本HelloELB", ToGoName("本 hello Elb"))
		})
	})
}
