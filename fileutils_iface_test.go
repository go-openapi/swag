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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		t.Run("with test package path", func(t *testing.T) {
			const tgt = "testpath"

			td, err := os.MkdirTemp("", tgt) //nolint:usetesting // as t.TempDir in testing not yet fully working (on windows)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = os.RemoveAll(td)
			})

			realPath := filepath.Join(td, "src", "foo", "bar")
			require.NoError(t, os.MkdirAll(realPath, os.ModePerm))

			assert.NotEmpty(t, FindInSearchPath(td, "foo/bar"))
		})

		// The following functions are environment-dependant and difficult to test.
		// Deferred to in-package unit testing.
		assert.NotPanics(t, func() {
			_ = FullGoSearchPath()
		})
		assert.NotPanics(t, func() {
			_ = FindInGoSearchPath("github.com/go-openapi/swag")
		})
	})
}
