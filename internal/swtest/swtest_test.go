package swtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeYAML(t *testing.T) {
	serv := httptest.NewServer(http.HandlerFunc(ServeYAMLPestore))
	defer serv.Close()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, serv.URL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
