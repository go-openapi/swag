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
		assert.EqualT(t, "HelloSwagger", ToGoName(sample))
		assert.EqualT(t, "HTTPServer", ToGoName(sampleWithInitialism))
		assert.EqualT(t, "helloSwagger", ToVarName(sample))
		assert.EqualT(t, "httpServer", ToVarName(sampleWithInitialism))
		assert.EqualT(t, "hello_swagger", ToFileName(sample))
		assert.EqualT(t, "http_server", ToFileName(sampleWithInitialism))
		assert.EqualT(t, "hello-swagger", ToCommandName(sample))
		assert.EqualT(t, "http-server", ToCommandName(sampleWithInitialism))
		assert.EqualT(t, "hello swagger", ToHumanNameLower(sample))
		assert.EqualT(t, "Http server", ToHumanNameLower(sampleWithInitialism))
		assert.EqualT(t, "Hello Swagger", ToHumanNameTitle(sample))
		assert.EqualT(t, "Http Server", ToHumanNameTitle(sampleWithInitialism))

		assert.EqualT(t, "Swagger", Camelize("SWAGGER"))
		assert.EqualT(t, "helloSwagger", ToJSONName(sample))

		t.Run("with global config", func(t *testing.T) {
			AddInitialisms("ELB")                                   // adding non-default initialism
			GoNamePrefixFunc = func(_ string) string { return "Z" } // adding non-default prefix function

			assert.EqualT(t, "Z本HelloELB", ToGoName("本 hello Elb"))
		})
	})
}
