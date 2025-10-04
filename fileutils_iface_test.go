// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestFileUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		t.Run("with test package path", func(t *testing.T) {
			td := t.TempDir()

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
