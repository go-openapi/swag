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

package loading

import (
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromHTTP(t *testing.T) {
	t.Run("should load pet store API doc", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(serveYAMLPestore))
		defer ts.Close()

		content, err := LoadFromFileOrHTTP(ts.URL)
		require.NoError(t, err)

		assert.YAMLEq(t, string(yamlPetStore), string(content))
	})

	t.Run("should not load from invalid URL", func(t *testing.T) {
		_, err := LoadFromFileOrHTTP("httx://12394:abd")
		require.Error(t, err)
	})

	t.Run("should not load from remote URL with error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(serveKO))
		defer ts.Close()

		_, err := LoadFromFileOrHTTP(ts.URL)
		require.Error(t, err)
	})

	t.Run("should load from remote URL", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(serveOK))
		defer ts.Close()

		d, err := LoadFromFileOrHTTP(ts.URL)
		require.NoError(t, err)
		assert.Equal(t, []byte("the content"), d)
	})

	t.Run("with remote basic auth", func(t *testing.T) {
		const (
			validUsername   = "fake-user"
			validPassword   = "correct-password"
			invalidPassword = "incorrect-password"
		)

		ts := httptest.NewServer(http.HandlerFunc(serveBasicAuthFunc(validUsername, validPassword)))
		defer ts.Close()

		t.Run("should not load from remote URL unauthenticated", func(t *testing.T) {
			_, err := LoadFromFileOrHTTP(ts.URL) // no auth
			require.Error(t, err)
		})

		t.Run("using loading options", func(t *testing.T) {
			t.Run("should not load from remote URL with invalid credentials", func(t *testing.T) {
				_, err := LoadFromFileOrHTTP(ts.URL,
					WithBasicAuth(validUsername, invalidPassword),
				)
				require.Error(t, err)
			})

			t.Run("should load from remote URL with basic auth", func(t *testing.T) {
				_, err := LoadFromFileOrHTTP(ts.URL,
					WithBasicAuth(validUsername, validPassword), // basic auth, valid credentials
				)
				require.NoError(t, err)
			})
		})
	})

	t.Run("with remote API key auth", func(t *testing.T) {
		const (
			sharedHeaderKey   = "X-Myapp"
			sharedHeaderValue = "MySecretKey"
		)

		ts := httptest.NewServer(http.HandlerFunc(serveRequireHeaderFunc(sharedHeaderKey, sharedHeaderValue)))
		defer ts.Close()

		t.Run("using loading options", func(t *testing.T) {
			t.Run("should not load from remote URL with missing header", func(t *testing.T) {
				_, err := LoadFromFileOrHTTP(ts.URL)
				require.Error(t, err)
			})

			t.Run("should load from remote URL with API key header", func(t *testing.T) {
				_, err := LoadFromFileOrHTTP(ts.URL,
					WithCustomHeaders(map[string]string{sharedHeaderKey: sharedHeaderValue}),
				)
				require.NoError(t, err)
			})

			t.Run("with custom HTTP client mocking a remote", func(t *testing.T) {
				cwd, _ := os.Getwd()
				fixtureDir := filepath.Join(cwd, "fixtures")
				client := &http.Client{
					// intercepts calls to the server and serves local files instead
					Transport: http.NewFileTransport(http.Dir(fixtureDir)),
				}

				t.Run("should not load unknown path", func(t *testing.T) {
					_, err := LoadFromFileOrHTTP(ts.URL+"/unknown",
						WithCustomHeaders(map[string]string{sharedHeaderKey: sharedHeaderValue}),
						WithHTTPClient(client),
					)
					require.Error(t, err)
				})

				t.Run("should load from local path", func(t *testing.T) {
					petstore, err := LoadFromFileOrHTTP(ts.URL+"/petstore_fixture.yaml",
						WithCustomHeaders(map[string]string{sharedHeaderKey: sharedHeaderValue}),
						WithHTTPClient(client),
					)
					require.NoError(t, err)
					require.NotEmpty(t, petstore)
				})
			})
		})
	})

	t.Run("should not load when timeout", func(t *testing.T) {
		const (
			delay = 30 * time.Millisecond
			wait  = delay / 2
		)

		serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			time.Sleep(delay)
			rw.WriteHeader(http.StatusOK)
		}))
		defer serv.Close()

		t.Run("using loading options", func(t *testing.T) {
			_, err := LoadFromFileOrHTTP(serv.URL,
				WithTimeout(wait),
			)
			require.Error(t, err)
			require.ErrorIs(t, err, context.DeadlineExceeded)
		})

		t.Run("disabling timeout with loading options", func(t *testing.T) {
			_, err := LoadFromFileOrHTTP(serv.URL,
				WithTimeout(0),
			)
			require.NoError(t, err)
		})
	})

	t.Run("should load from local embedded file system (single file)", func(t *testing.T) {
		// using plain fs.FS
		rooted, err := fs.Sub(embeddedFixtures, "fixtures")
		require.NoError(t, err)
		b, err := LoadFromFileOrHTTP("petstore_fixture.yaml",
			WithFS(rooted),
		)
		require.NoError(t, err)
		assert.YAMLEq(t, string(yamlPetStore), string(b))
	})

	t.Run("should load from memory file system (single file)", func(t *testing.T) {
		mapfs := make(fstest.MapFS)
		mapfs["file"] = &fstest.MapFile{Data: []byte("content"), Mode: fs.ModePerm}
		// using fs.ReadFileFS
		b, err := LoadFromFileOrHTTP("file",
			WithFS(mapfs),
		)
		require.NoError(t, err)
		assert.Equal(t, string("content"), string(b))
	})

	t.Run("should load from local embedded file system (path)", func(t *testing.T) {
		// using plain fs.ReadFileFS
		// NOTE: this doesn't work on windows, because embed.FS uses / even on windows
		b, err := LoadFromFileOrHTTP("fixtures/petstore_fixture.yaml",
			WithFS(embeddedFixtures),
		)
		require.NoError(t, err)
		assert.YAMLEq(t, string(yamlPetStore), string(b))
	})
}

func TestLoadStrategy(t *testing.T) {
	const thisIsNotIt = "not it"
	loader := func(_ string) ([]byte, error) {
		return yamlPetStore, nil
	}
	remLoader := func(_ string) ([]byte, error) {
		return []byte(thisIsNotIt), nil
	}

	t.Run("should serve local strategy", func(t *testing.T) {
		ldr := LoadStrategy("blah", loader, remLoader)
		b, _ := ldr("")
		assert.YAMLEq(t, string(yamlPetStore), string(b))
	})

	t.Run("should serve remote strategy with http", func(t *testing.T) {
		ldr := LoadStrategy("http://blah", loader, remLoader)
		b, _ := ldr("")
		assert.Equal(t, thisIsNotIt, string(b))
	})

	t.Run("should serve remote strategy with https", func(t *testing.T) {
		ldr := LoadStrategy("https://blah", loader, remLoader)
		b, _ := ldr("")
		assert.Equal(t, thisIsNotIt, string(b))
	})
}

func TestLoadStrategyFile(t *testing.T) {
	const (
		thisIsIt    = "thisIsIt"
		thisIsNotIt = "not it"
	)

	type strategyTest struct {
		Title           string
		Path            string
		Expected        string
		ExpectedWindows string
		ExpectError     bool
	}

	t.Run("with local file strategy", func(t *testing.T) {
		loader := func(called *bool, pth *string) func(string) ([]byte, error) {
			return func(p string) ([]byte, error) {
				*called = true
				*pth = p
				return []byte(thisIsIt), nil
			}
		}

		remLoader := func(_ string) ([]byte, error) {
			return []byte(thisIsNotIt), nil
		}

		for _, toPin := range []strategyTest{
			{
				Title:           "valid fully qualified local URI, with rooted path",
				Path:            "file:///a/c/myfile.yaml",
				Expected:        "/a/c/myfile.yaml",
				ExpectedWindows: `\a\c\myfile.yaml`,
			},
			{
				Title:           "local URI with scheme, with host segment before path",
				Path:            "file://a/c/myfile.yaml",
				Expected:        "a/c/myfile.yaml",
				ExpectedWindows: `\\a\c\myfile.yaml`, // UNC host
			},
			{
				Title:           "local URI with scheme, with escaped characters",
				Path:            "file://a/c/myfile%20%28x86%29.yaml",
				Expected:        "a/c/myfile (x86).yaml",
				ExpectedWindows: `\\a\c\myfile (x86).yaml`,
			},
			{
				Title:           "local URI with scheme, rooted, with escaped characters",
				Path:            "file:///a/c/myfile%20%28x86%29.yaml",
				Expected:        "/a/c/myfile (x86).yaml",
				ExpectedWindows: `\a\c\myfile (x86).yaml`,
			},
			{
				Title:           "local URI with scheme, unescaped, with host",
				Path:            "file://a/c/myfile (x86).yaml",
				Expected:        "a/c/myfile (x86).yaml",
				ExpectedWindows: `\\a\c\myfile (x86).yaml`,
			},
			{
				Title:           "local URI with scheme, rooted, unescaped",
				Path:            "file:///a/c/myfile (x86).yaml",
				Expected:        "/a/c/myfile (x86).yaml",
				ExpectedWindows: `\a\c\myfile (x86).yaml`,
			},
			{
				Title:    "file URI with drive letter and backslashes, as a relative Windows path",
				Path:     `file://C:\a\c\myfile.yaml`,
				Expected: `C:\a\c\myfile.yaml`, // outcome on all platforms, not only windows
			},
			{
				Title:           "file URI with drive letter and backslashes, as a rooted Windows path",
				Path:            `file:///C:\a\c\myfile.yaml`,
				Expected:        `/C:\a\c\myfile.yaml`, // on non-windows, this results most likely in a wrong path
				ExpectedWindows: `C:\a\c\myfile.yaml`,  // on windows, we know that C: is a drive letter, so /C: becomes C:
			},
			{
				Title:    "file URI with escaped backslashes",
				Path:     `file://C%3A%5Ca%5Cc%5Cmyfile.yaml`,
				Expected: `C:\a\c\myfile.yaml`, // outcome on all platforms, not only windows
			},
			{
				Title:           "file URI with escaped backslashes, rooted",
				Path:            `file:///C%3A%5Ca%5Cc%5Cmyfile.yaml`,
				Expected:        `/C:\a\c\myfile.yaml`, // outcome on non-windows (most likely not a desired path)
				ExpectedWindows: `C:\a\c\myfile.yaml`,  // outcome on windows
			},
			{
				Title:           "URI with the file scheme, host omitted: relative path with extra dots",
				Path:            `file://./a/c/d/../myfile.yaml`,
				Expected:        `./a/c/d/../myfile.yaml`,
				ExpectedWindows: `a\c\myfile.yaml`, // on windows, extra processing cleans the path
			},
			{
				Title:           "relative URI without the file scheme, rooted path",
				Path:            `/a/c/myfile.yaml`,
				Expected:        `/a/c/myfile.yaml`,
				ExpectedWindows: `\a\c\myfile.yaml`, // there is no drive letter, this would probably result in a wrong path on Windows
			},
			{
				Title:           "relative URI without the file scheme, relative path",
				Path:            `a/c/myfile.yaml`,
				Expected:        `a/c/myfile.yaml`,
				ExpectedWindows: `a\c\myfile.yaml`,
			},
			{
				Title:           "relative URI without the file scheme, relative path with dots",
				Path:            `./a/c/myfile.yaml`,
				Expected:        `./a/c/myfile.yaml`,
				ExpectedWindows: `.\a\c\myfile.yaml`,
			},
			{
				Title:           "relative URI without the file scheme, relative path with extra dots",
				Path:            `./a/c/../myfile.yaml`,
				Expected:        `./a/c/../myfile.yaml`,
				ExpectedWindows: `.\a\c\..\myfile.yaml`,
			},
			{
				Title:           "relative URI without the file scheme, windows slashed-path with drive letter",
				Path:            `A:/a/c/myfile.yaml`,
				Expected:        `A:/a/c/myfile.yaml`, // on non-windows, this results most likely in a wrong path
				ExpectedWindows: `A:\a\c\myfile.yaml`, // on windows, slashes are converted
			},
			{
				Title:           "relative URI without the file scheme, windows backslashed-path with drive letter",
				Path:            `A:\a\c\myfile.yaml`,
				Expected:        `A:\a\c\myfile.yaml`, // on non-windows, this results most likely in a wrong path
				ExpectedWindows: `A:\a\c\myfile.yaml`,
			},
			{
				Title:           "URI with file scheme, host as Windows UNC name",
				Path:            `file://host/share/folder/myfile.yaml`,
				Expected:        `host/share/folder/myfile.yaml`,   // there is no host component accounted for
				ExpectedWindows: `\\host\share\folder\myfile.yaml`, // on windows, the host is interpreted as an UNC host for a file share
			},
			{
				Title:       "invalid URL encoding",
				Path:        `/folder%GF/myfile.yaml`,
				ExpectError: true,
			},
		} {
			tc := toPin
			t.Run(tc.Title, func(t *testing.T) {
				var (
					called bool
					pth    string
				)

				loader := LoadStrategy("local", loader(&called, &pth), remLoader)
				b, err := loader(tc.Path)
				if tc.ExpectError {
					require.Error(t, err)
					assert.False(t, called)

					return
				}

				require.NoError(t, err)
				assert.True(t, called)
				assert.Equal(t, []byte(thisIsIt), b)

				if tc.ExpectedWindows != "" && runtime.GOOS == "windows" {
					assert.Equalf(t, tc.ExpectedWindows, pth,
						"expected local LoadStrategy(%q) to open: %q (windows)",
						tc.Path, tc.ExpectedWindows,
					)

					return
				}

				assert.Equalf(t, tc.Expected, pth,
					"expected local LoadStrategy(%q) to open: %q (any OS)",
					tc.Path, tc.Expected,
				)
			})
		}
	})
}
