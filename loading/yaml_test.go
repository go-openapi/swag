// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package loading

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-openapi/testify/enable/yaml/v2" // enable YAMLEq in testify
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestYAMLMatcher(t *testing.T) {
	t.Run("should recognize a yaml file", func(t *testing.T) {
		assert.True(t, YAMLMatcher("local.yml"))
		assert.True(t, YAMLMatcher("local.yaml"))
		assert.False(t, YAMLMatcher("local.json"))
	})
}

func TestYAMLDoc(t *testing.T) {
	t.Run("should retrieve pet store API as YAML", func(t *testing.T) {
		serv := httptest.NewServer(http.HandlerFunc(serveYAMLPetStore))
		defer serv.Close()

		s, err := YAMLDoc(serv.URL)
		require.NoError(t, err)
		require.NotNil(t, s)
	})

	t.Run("should not retrieve any doc", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write([]byte("\n"))
		}))
		defer ts.Close()

		_, err := YAMLDoc(ts.URL)
		require.Error(t, err)
	})
}
