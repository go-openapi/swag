// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loading

import (
	"embed"
	"fmt"
	"net/http"
)

// embedded test files

//go:embed petstore_fixture.yaml
var embeddedFixtures embed.FS

// YAMLPetStore embeds the classical pet store API swagger example.
var YAMLPetStore []byte

func init() {
	data, err := embeddedFixtures.ReadFile("petstore_fixture.yaml")
	if err != nil {
		panic(fmt.Errorf("wrong embedded FS configuration: %w", err))
	}

	YAMLPetStore = data
}

// test handlers

// serveYAMLPestore is a http handler to serve the YAMLPestore doc.
func serveYAMLPestore(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(YAMLPetStore)
}

func serveOK(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("the content"))
}

func serveKO(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

func serveBasicAuthFunc(user, password string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok && u == user && p == password {
			rw.WriteHeader(http.StatusOK)

			return
		}

		rw.WriteHeader(http.StatusForbidden)
	}
}

func serveRequireHeaderFunc(key, value string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		myHeaders := r.Header[key]
		ok := false
		for _, v := range myHeaders {
			if v == value {
				ok = true
				break
			}
		}
		if ok {
			rw.WriteHeader(http.StatusOK)

			return
		}

		rw.WriteHeader(http.StatusForbidden)
	}
}
