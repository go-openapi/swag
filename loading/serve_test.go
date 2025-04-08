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
	"os"
	"path"
	"testing"
)

// embedded test files

//go:embed fixtures/*
var embeddedFixtures embed.FS

// yamlPetStore embeds the classical pet store API swagger example.
var yamlPetStore []byte
var jsonPetStore []byte

func TestMain(m *testing.M) {
	yamlPetStore = mustLoadFixture("petstore_fixture.yaml")
	jsonPetStore = mustLoadFixture("petstore_fixture.json")

	os.Exit(m.Run())
}

// test handlers

// serveYAMLPestore is a http handler to serve the YAMLPestore doc.
func serveYAMLPestore(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(yamlPetStore)
}

// serveJSONPestore is a http handler to serve the jsonPestore doc.
func serveJSONPestore(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(jsonPetStore)
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

func mustLoadFixture(name string) []byte {
	const msg = "wrong embedded FS configuration: %w"
	data, err := embeddedFixtures.ReadFile(path.Join("fixtures", name))
	if err != nil {
		panic(fmt.Errorf(msg, err))
	}

	return data
}
