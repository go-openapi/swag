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
		assert.True(t, JSONMatcher("local.json"))
		assert.True(t, JSONMatcher("local.jso"))
		assert.True(t, JSONMatcher("local.jsn"))
		assert.False(t, JSONMatcher("local.yml"))
	})
}

func TestJSONDoc(t *testing.T) {
	t.Run("should retrieve pet store API as JSON", func(t *testing.T) {
		serv := httptest.NewServer(http.HandlerFunc(serveJSONPestore))

		defer serv.Close()

		s, err := JSONDoc(serv.URL)
		require.NoError(t, err)
		require.NotNil(t, s)
		require.JSONEq(t, string(jsonPetStore), string(s))
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
