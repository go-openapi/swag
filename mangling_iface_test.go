// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
