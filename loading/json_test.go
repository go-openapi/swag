// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package loading

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestJSONMatcher(t *testing.T) {
	t.Run("should recognize a json file", func(t *testing.T) {
		assert.TrueT(t, JSONMatcher("local.json"))
		assert.TrueT(t, JSONMatcher("local.jso"))
		assert.TrueT(t, JSONMatcher("local.jsn"))
		assert.FalseT(t, JSONMatcher("local.yml"))
	})
}

func TestJSONDoc(t *testing.T) {
	t.Run("should retrieve pet store API as JSON", func(t *testing.T) {
		serv := httptest.NewServer(http.HandlerFunc(serveJSONPetStore))

		defer serv.Close()

		s, err := JSONDoc(serv.URL)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.JSONEqBytes(t, jsonPetStore, s)
	})

	t.Run("should not retrieve any doc", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write([]byte("\n"))
		}))
		defer ts.Close()

		_, err := JSONDoc(ts.URL)
		require.Error(t, err)
	})
}
