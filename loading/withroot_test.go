// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package loading

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestWithRoot(t *testing.T) {
	const (
		inside = "inside the root"
		nested = "nested in the root"
		secret = "this is a secret outside the root"
	)

	// layout:
	//   <parent>/secret.txt        (outside the root)
	//   <parent>/root/api.yaml
	//   <parent>/root/sub/api.yaml
	parent := t.TempDir()
	root := filepath.Join(parent, "root")
	require.NoError(t, os.MkdirAll(filepath.Join(root, "sub"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(parent, "secret.txt"), []byte(secret), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "api.yaml"), []byte(inside), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "sub", "api.yaml"), []byte(nested), 0o600))

	t.Run("should load paths confined to the root", func(t *testing.T) {
		for _, pth := range []string{
			"api.yaml",
			"./api.yaml",
			"file://api.yaml",
			"sub/../api.yaml",
		} {
			t.Run(pth, func(t *testing.T) {
				b, err := LoadFromFileOrHTTP(pth, WithRoot(root))
				require.NoError(t, err)
				assert.EqualT(t, inside, string(b))
			})
		}

		t.Run("nested path", func(t *testing.T) {
			b, err := LoadFromFileOrHTTP("sub/api.yaml", WithRoot(root))
			require.NoError(t, err)
			assert.EqualT(t, nested, string(b))
		})
	})

	t.Run("should reject paths escaping the root", func(t *testing.T) {
		for _, pth := range []string{
			"file:///etc/passwd",                // absolute via file:// URI
			filepath.Join(parent, "secret.txt"), // absolute path to an existing sibling file
			"../secret.txt",                     // traversal
			"file://../secret.txt",              // traversal via file:// URI
		} {
			t.Run(pth, func(t *testing.T) {
				b, err := LoadFromFileOrHTTP(pth, WithRoot(root))
				require.Error(t, err)
				// the rejected read must not leak any bytes
				assert.Empty(t, b)
			})
		}
	})

	t.Run("should reject a symlink escaping the root", func(t *testing.T) {
		// <root>/escape.yaml -> <parent>/secret.txt  (escapes the root)
		if err := os.Symlink(filepath.Join(parent, "secret.txt"), filepath.Join(root, "escape.yaml")); err != nil {
			t.Skipf("symlinks not supported on this platform/filesystem: %v", err)
		}

		b, err := LoadFromFileOrHTTP("escape.yaml", WithRoot(root))
		require.Error(t, err)
		assert.Empty(t, b)
		// this is exactly the case os.DirFS would NOT block
		assert.NotContains(t, string(b), secret)
	})

	t.Run("should surface an error for a missing root", func(t *testing.T) {
		b, err := LoadFromFileOrHTTP("api.yaml", WithRoot(filepath.Join(parent, "does-not-exist")))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrLoader)
		assert.Empty(t, b)
	})

	t.Run("WithRoot and WithFS are mutually exclusive (last wins)", func(t *testing.T) {
		mapfs := fstest.MapFS{"api.yaml": &fstest.MapFile{Data: []byte("from map fs"), Mode: fs.ModePerm}}

		t.Run("WithRoot after WithFS uses the root", func(t *testing.T) {
			b, err := LoadFromFileOrHTTP("api.yaml", WithFS(mapfs), WithRoot(root))
			require.NoError(t, err)
			assert.EqualT(t, inside, string(b))
		})

		t.Run("WithFS after WithRoot uses the fs", func(t *testing.T) {
			b, err := LoadFromFileOrHTTP("api.yaml", WithRoot(root), WithFS(mapfs))
			require.NoError(t, err)
			assert.EqualT(t, "from map fs", string(b))
		})
	})

	t.Run("default loader (no WithRoot) is unchanged and reads outside any root", func(t *testing.T) {
		// regression guard: the fix is opt-in and must not change default behavior
		b, err := LoadFromFileOrHTTP(filepath.Join(parent, "secret.txt"))
		require.NoError(t, err)
		assert.EqualT(t, secret, string(b))
	})
}
