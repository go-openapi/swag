// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

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

// serveYAMLPetStore is a http handler to serve the yamlPetStore doc.
func serveYAMLPetStore(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(yamlPetStore)
}

// serveJSONPetStore is a http handler to serve the jsonPetStore doc.
func serveJSONPetStore(rw http.ResponseWriter, _ *http.Request) {
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
