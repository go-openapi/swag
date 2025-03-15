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
package swag

import (
	"time"

	"github.com/go-openapi/swag/loading"
)

var (
	// LoadHTTPTimeout the default timeout for load requests
	//
	// Deprecated: use [loading.WithTimeout] instead.
	LoadHTTPTimeout = 30 * time.Second

	// LoadHTTPBasicAuthUsername the username to use when load requests require basic auth
	//
	// Deprecated: use [loading.WithBasicAuth] instead.
	LoadHTTPBasicAuthUsername = ""

	// LoadHTTPBasicAuthPassword the password to use when load requests require basic auth
	//
	// Deprecated: use [loading.WithBasicAuth] instead.
	LoadHTTPBasicAuthPassword = ""

	// LoadHTTPCustomHeaders an optional collection of custom HTTP headers for load requests
	//
	// Deprecated: use [loading.WithCustomHeaders] instead.
	LoadHTTPCustomHeaders = map[string]string{}
)

// LoadFromFileOrHTTP loads the bytes from a file or a remote http server based on the path passed in
//
// See [loading.LoadFromFileOrHTTP].
func LoadFromFileOrHTTP(pth string, opts ...loading.Option) ([]byte, error) {
	return loading.LoadFromFileOrHTTP(pth, loadingOptionsWithDefaults(opts)...)
}

// LoadFromFileOrHTTPWithTimeout loads the bytes from a file or a remote http server based on the path passed in
// timeout arg allows for per request overriding of the request timeout.
//
// Deprecated: use [loading.LoadFileOrHTTP] with the [WithTimeout] option instead.
func LoadFromFileOrHTTPWithTimeout(pth string, timeout time.Duration, opts ...loading.Option) ([]byte, error) {
	opts = append(opts, loading.WithTimeout(timeout))

	return LoadFromFileOrHTTP(pth, opts...)
}

// LoadStrategy returns a loader function for a given path or URI.
//
// See [loading.LoadStrategy].
func LoadStrategy(pth string, local, remote func(string) ([]byte, error), opts ...loading.Option) func(string) ([]byte, error) {
	return loading.LoadStrategy(pth, local, remote, loadingOptionsWithDefaults(opts)...)
}

// loadingOptionsWithDefaults bridges deprecated default settings using package-level variables,
// with the recommended use of loading.Option.
func loadingOptionsWithDefaults(opts []loading.Option) []loading.Option {
	o := []loading.Option{
		loading.WithTimeout(LoadHTTPTimeout),
		loading.WithBasicAuth(LoadHTTPBasicAuthUsername, LoadHTTPBasicAuthUsername),
		loading.WithCustomHeaders(LoadHTTPCustomHeaders),
	}
	o = append(o, opts...)

	return o
}
