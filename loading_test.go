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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validUsername     = "fake-user"
	validPassword     = "correct-password"
	invalidPassword   = "incorrect-password"
	sharedHeaderKey   = "X-Myapp"
	sharedHeaderValue = "MySecretKey"
)

func TestLoadFromHTTP(t *testing.T) {

	_, err := LoadFromFileOrHTTP("httx://12394:abd")
	assert.Error(t, err)

	serv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer serv.Close()

	_, err = LoadFromFileOrHTTP(serv.URL)
	assert.Error(t, err)

	ts2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("the content"))
	}))
	defer ts2.Close()

	d, err := LoadFromFileOrHTTP(ts2.URL)
	assert.NoError(t, err)
	assert.Equal(t, []byte("the content"), d)

	ts3 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok && u == validUsername && p == validPassword {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
	}))
	defer ts3.Close()

	// no auth
	_, err = LoadFromFileOrHTTP(ts3.URL)
	assert.Error(t, err)

	// basic auth, invalide credentials
	LoadHTTPBasicAuthUsername = validUsername
	LoadHTTPBasicAuthPassword = invalidPassword

	_, err = LoadFromFileOrHTTP(ts3.URL)
	assert.Error(t, err)

	// basic auth, valid credentials
	LoadHTTPBasicAuthUsername = validUsername
	LoadHTTPBasicAuthPassword = validPassword

	_, err = LoadFromFileOrHTTP(ts3.URL)
	assert.NoError(t, err)

	ts4 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
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
	defer ts4.Close()

	_, err = LoadFromFileOrHTTP(ts4.URL)
	assert.Error(t, err)

	LoadHTTPCustomHeaders[sharedHeaderKey] = sharedHeaderValue

	_, err = LoadFromFileOrHTTP(ts4.URL)
	assert.NoError(t, err)

	// clean up for future tests
	LoadHTTPBasicAuthUsername = ""
	LoadHTTPBasicAuthPassword = ""
	LoadHTTPCustomHeaders = map[string]string{}
}
