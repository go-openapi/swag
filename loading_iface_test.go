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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromHTTP(t *testing.T) {
	// Check backward-compatible global config.
	// More comprehensive testing is carried out in package loading.

	t.Run("with remote basic auth", func(t *testing.T) {
		const (
			validUsername   = "fake-user"
			validPassword   = "correct-password"
			invalidPassword = "incorrect-password"
		)

		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if ok && u == validUsername && p == validPassword {
				rw.WriteHeader(http.StatusOK)

				return
			}

			rw.WriteHeader(http.StatusForbidden)
		}))
		defer ts.Close()

		t.Run("using global config", func(t *testing.T) {
			t.Cleanup(func() {
				LoadHTTPBasicAuthUsername = ""
				LoadHTTPBasicAuthPassword = ""
			})

			t.Run("should load from remote URL with basic auth", func(t *testing.T) {
				// basic auth, valid credentials
				LoadHTTPBasicAuthUsername = validUsername
				LoadHTTPBasicAuthPassword = validPassword

				_, err := LoadFromFileOrHTTP(ts.URL)
				require.NoError(t, err)
			})
		})
	})

	t.Run("with remote API key auth", func(t *testing.T) {
		const (
			sharedHeaderKey   = "X-Myapp"
			sharedHeaderValue = "MySecretKey"
		)

		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			myHeaders := r.Header[sharedHeaderKey]
			ok := false
			for _, v := range myHeaders {
				if v == sharedHeaderValue {
					ok = true
					break
				}
			}
			if ok {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.WriteHeader(http.StatusForbidden)
			}
		}))
		defer ts.Close()

		t.Run("using global config", func(t *testing.T) {
			t.Cleanup(func() {
				LoadHTTPCustomHeaders = map[string]string{}
			})

			t.Run("should load from remote URL with API key header", func(t *testing.T) {
				LoadHTTPCustomHeaders[sharedHeaderKey] = sharedHeaderValue

				_, err := LoadFromFileOrHTTP(ts.URL)
				require.NoError(t, err)
			})
		})
	})

	t.Run("should not load when timeout", func(t *testing.T) {
		const (
			delay = 30 * time.Millisecond
			wait  = delay / 2
		)

		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			time.Sleep(delay)
			rw.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		t.Run("using global configuration", func(t *testing.T) {
			original := LoadHTTPTimeout
			t.Cleanup(func() {
				LoadHTTPTimeout = original
			})
			LoadHTTPTimeout = wait

			_, err := LoadFromFileOrHTTP(ts.URL)
			require.Error(t, err)
			require.ErrorIs(t, err, context.DeadlineExceeded)
		})

		t.Run("using deprecated method", func(t *testing.T) {
			_, err := LoadFromFileOrHTTPWithTimeout(ts.URL, wait)
			require.Error(t, err)
			require.ErrorIs(t, err, context.DeadlineExceeded)
		})

		t.Run("should serve local strategy", func(t *testing.T) {
			loader := func(_ string) ([]byte, error) {
				return []byte("local"), nil
			}
			remLoader := func(_ string) ([]byte, error) {
				return []byte("remote"), nil
			}
			ldr := LoadStrategy("not_http", loader, remLoader)
			b, _ := ldr("")
			assert.Equal(t, "local", string(b))
		})
	})
}

func TestYAMLDoc(t *testing.T) {
	t.Run("deprecated loading YAML functions should work", func(t *testing.T) {
		require.True(t, YAMLMatcher("a.yml"))

		ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(http.StatusOK)
			const ydoc = "x:\n  a: one\n  b: two\n"
			_, _ = rw.Write([]byte(ydoc))
		}))
		defer ts.Close()

		b, err := YAMLDoc(ts.URL)
		require.NoError(t, err)
		require.NotEmpty(t, b)

		doc, err := YAMLData(ts.URL)
		require.NoError(t, err)
		require.NotEmpty(t, doc)
	})
}
