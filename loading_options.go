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
	"io/fs"
	"net/http"
	"os"
	"time"
)

type (
	// LoadingOption provides options for loading a file over HTTP or from a file.
	LoadingOption func(*loadingOptions)

	httpOptions struct {
		httpTimeout       time.Duration
		basicAuthUsername string
		basicAuthPassword string
		customHeaders     map[string]string
		client            *http.Client
	}

	fileOptions struct {
		fs fs.ReadFileFS
	}

	loadingOptions struct {
		httpOptions
		fileOptions
	}
)

func (fo fileOptions) ReadFileFunc() func(string) ([]byte, error) {
	if fo.fs == nil {
		return os.ReadFile
	}

	return fo.fs.ReadFile
}

// LoadingWithTimeout sets a timeout for the remote file loader.
func LoadingWithTimeout(timeout time.Duration) LoadingOption {
	return func(o *loadingOptions) {
		o.httpTimeout = timeout
	}
}

// LoadingWithBasicAuth sets a basic authentication scheme for the remote file loader.
func LoadingWithBasicAuth(username, password string) LoadingOption {
	return func(o *loadingOptions) {
		o.basicAuthUsername = username
		o.basicAuthPassword = password
	}
}

// LoadingWithCustomHeaders sets custom headers for the remote file loader.
func LoadingWithCustomHeaders(headers map[string]string) LoadingOption {
	return func(o *loadingOptions) {
		if o.customHeaders == nil {
			o.customHeaders = make(map[string]string, len(headers))
		}

		for header, value := range headers {
			o.customHeaders[header] = value
		}
	}
}

// LoadingWithHTTClient overrides the default HTTP client used to fetch a remote file.
//
// By default, [http.DefaultClient] is used.
func LoadingWithHTTPClient(client *http.Client) LoadingOption {
	return func(o *loadingOptions) {
		o.client = client
	}
}

// LoadingWithFS sets a file system for the local file loader.
//
// By default, the file system is the one provided by the os package.
//
// For example, this may be set to consume from an embedded file system, or a rooted FS.
func LoadingWithFS(fs fs.ReadFileFS) LoadingOption {
	return func(o *loadingOptions) {
		o.fs = fs
	}
}

func loadingOptionsWithDefaults(opts []LoadingOption) loadingOptions {
	o := loadingOptions{
		// package level defaults
		httpOptions: httpOptions{
			httpTimeout:       LoadHTTPTimeout,
			customHeaders:     LoadHTTPCustomHeaders,
			basicAuthUsername: LoadHTTPBasicAuthUsername,
			basicAuthPassword: LoadHTTPBasicAuthPassword,
			client:            http.DefaultClient,
		},
	}

	for _, apply := range opts {
		apply(&o)
	}

	return o
}
