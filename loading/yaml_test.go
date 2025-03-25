package loading

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		serv := httptest.NewServer(http.HandlerFunc(serveYAMLPestore))
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
