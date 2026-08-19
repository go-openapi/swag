// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// openOnlyFS hides all the extension interfaces of the wrapped file system,
// leaving only [fs.FS].
//
// It is used to exercise the fallback paths that apply when a layer implements
// nothing more than [fs.FS.Open].
type openOnlyFS struct {
	inner fs.FS
}

func (f openOnlyFS) Open(name string) (fs.File, error) { return f.inner.Open(name) }

// makeOsFSFixture populates a temporary directory with a few files and returns
// its path with forward slashes, so that paths built by [path.Join] and results
// returned by [fs.Glob] compare equal on all platforms.
func makeOsFSFixture(tb testing.TB) string {
	tb.Helper()

	root := tb.TempDir()
	require.NoError(tb, os.MkdirAll(filepath.Join(root, "sub"), 0o750))
	require.NoError(tb, os.WriteFile(filepath.Join(root, "a.txt"), []byte("content of a"), 0o600))
	require.NoError(tb, os.WriteFile(filepath.Join(root, "b.txt"), []byte("content of b"), 0o600))
	require.NoError(tb, os.WriteFile(filepath.Join(root, "sub", "c.json"), []byte(`{"c": true}`), 0o600))

	return filepath.ToSlash(root)
}

func TestOsFS(t *testing.T) {
	root := makeOsFSFixture(t)
	osFS := NewReadOnlyOsFS()

	t.Run("should implement the read-only file system interfaces", func(t *testing.T) {
		assert.Implements(t, new(fs.FS), osFS)
		assert.Implements(t, new(fs.ReadFileFS), osFS)
		assert.Implements(t, new(fs.ReadDirFS), osFS)
	})

	t.Run("with Open", func(t *testing.T) {
		t.Run("should open a file from an absolute path", func(t *testing.T) {
			file, err := osFS.Open(path.Join(root, "a.txt"))
			require.NoError(t, err)
			t.Cleanup(func() { _ = file.Close() })

			data, err := io.ReadAll(file)
			require.NoError(t, err)
			assert.EqualT(t, "content of a", string(data))
		})

		t.Run("should open a file from a relative path", func(t *testing.T) {
			// unlike [os.DirFS], [OsFS] is not rooted: relative paths resolve
			// against the current working directory.
			file, err := osFS.Open("fs.go")
			require.NoError(t, err)
			require.NoError(t, file.Close())
		})

		t.Run("should not open a file that does not exist", func(t *testing.T) {
			_, err := osFS.Open(path.Join(root, "nowhere.txt"))
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

	t.Run("with ReadFile", func(t *testing.T) {
		t.Run("should read a file", func(t *testing.T) {
			data, err := osFS.ReadFile(path.Join(root, "sub", "c.json"))
			require.NoError(t, err)
			assert.EqualT(t, `{"c": true}`, string(data))
		})

		t.Run("should not read a file that does not exist", func(t *testing.T) {
			_, err := osFS.ReadFile(path.Join(root, "nowhere.txt"))
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

	t.Run("with ReadDir", func(t *testing.T) {
		t.Run("should list a directory, sorted by file name", func(t *testing.T) {
			entries, err := osFS.ReadDir(root)
			require.NoError(t, err)

			names := make([]string, 0, len(entries))
			for _, entry := range entries {
				names = append(names, entry.Name())
			}
			assert.SliceEqualT(t, []string{"a.txt", "b.txt", "sub"}, names)
		})

		t.Run("should not list a directory that does not exist", func(t *testing.T) {
			_, err := osFS.ReadDir(path.Join(root, "nowhere"))
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

	t.Run("should work with the io/fs helpers", func(t *testing.T) {
		t.Run("with fs.ReadFile", func(t *testing.T) {
			data, err := fs.ReadFile(osFS, path.Join(root, "b.txt"))
			require.NoError(t, err)
			assert.EqualT(t, "content of b", string(data))
		})

		t.Run("with fs.Stat", func(t *testing.T) {
			// [OsFS] does not implement [fs.StatFS]: this exercises the
			// open-then-stat fallback in [fs.Stat].
			info, err := fs.Stat(osFS, path.Join(root, "sub"))
			require.NoError(t, err)
			assert.TrueT(t, info.IsDir())
		})
	})
}

func TestGlobOsFS(t *testing.T) {
	root := makeOsFSFixture(t)
	globFS := NewGlobOsFS()

	t.Run("should implement the read-only file system interfaces, plus fs.GlobFS", func(t *testing.T) {
		assert.Implements(t, new(fs.FS), globFS)
		assert.Implements(t, new(fs.ReadFileFS), globFS)
		assert.Implements(t, new(fs.ReadDirFS), globFS)
		assert.Implements(t, new(fs.GlobFS), globFS)
	})

	t.Run("should inherit the OsFS methods", func(t *testing.T) {
		data, err := globFS.ReadFile(path.Join(root, "a.txt"))
		require.NoError(t, err)
		assert.EqualT(t, "content of a", string(data))
	})

	t.Run("with Glob", func(t *testing.T) {
		t.Run("should match a pattern", func(t *testing.T) {
			matches, err := globFS.Glob(path.Join(root, "*.txt"))
			require.NoError(t, err)
			assert.SliceEqualT(t, []string{
				path.Join(root, "a.txt"),
				path.Join(root, "b.txt"),
			}, matches)
		})

		t.Run("should match a pattern in a sub-directory", func(t *testing.T) {
			matches, err := globFS.Glob(path.Join(root, "*", "*.json"))
			require.NoError(t, err)
			assert.SliceEqualT(t, []string{path.Join(root, "sub", "c.json")}, matches)
		})

		t.Run("should match a pattern without any meta character", func(t *testing.T) {
			matches, err := globFS.Glob(path.Join(root, "a.txt"))
			require.NoError(t, err)
			assert.SliceEqualT(t, []string{path.Join(root, "a.txt")}, matches)
		})

		t.Run("should return no match", func(t *testing.T) {
			matches, err := globFS.Glob(path.Join(root, "*.md"))
			require.NoError(t, err)
			assert.Empty(t, matches)
		})

		t.Run("should not match an unreadable directory", func(t *testing.T) {
			matches, err := globFS.Glob(path.Join(root, "nowhere", "*.txt"))
			require.NoError(t, err)
			assert.Empty(t, matches)
		})

		t.Run("should error on a malformed pattern", func(t *testing.T) {
			_, err := globFS.Glob(path.Join(root, "[a-"))
			require.Error(t, err)
			assert.ErrorIs(t, err, path.ErrBadPattern)
		})
	})
}

func TestFileReaderFS(t *testing.T) {
	// a base file system that provides Open and nothing else
	base := openOnlyFS{
		inner: fstest.MapFS{
			"a.txt":      &fstest.MapFile{Data: []byte("content of a")},
			"empty.txt":  &fstest.MapFile{Data: []byte{}},
			"sub/b.json": &fstest.MapFile{Data: []byte(`{"b": true}`)},
			"sub/c.yaml": &fstest.MapFile{Data: []byte("c: true")},
		},
	}
	readerFS := NewFileReaderFS(base)

	t.Run("should turn a fs.FS into a fs.ReadFileFS", func(t *testing.T) {
		assert.NotImplements(t, new(fs.ReadFileFS), base)
		assert.Implements(t, new(fs.FS), readerFS)
		assert.Implements(t, new(fs.ReadFileFS), readerFS)
	})

	t.Run("should read a file", func(t *testing.T) {
		data, err := readerFS.ReadFile("sub/b.json")
		require.NoError(t, err)
		assert.EqualT(t, `{"b": true}`, string(data))
	})

	t.Run("should read an empty file", func(t *testing.T) {
		data, err := readerFS.ReadFile("empty.txt")
		require.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("should not read a file that does not exist", func(t *testing.T) {
		_, err := readerFS.ReadFile("nowhere.txt")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("should delegate Open to the base file system", func(t *testing.T) {
		file, err := readerFS.Open("a.txt")
		require.NoError(t, err)
		t.Cleanup(func() { _ = file.Close() })

		data, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.EqualT(t, "content of a", string(data))
	})
}

func TestMustSub(t *testing.T) {
	base, err := NewMapFS(FromRawMap(map[string][]byte{
		"templates/index.html":              []byte("<h1>index</h1>"),
		"templates/contrib/mine/index.html": []byte("<h1>mine</h1>"),
	}))
	require.NoError(t, err)

	t.Run("should re-root a file system", func(t *testing.T) {
		subFS := MustSub(base, "templates")

		data, err := fs.ReadFile(subFS, "index.html")
		require.NoError(t, err)
		assert.Equal(t, "<h1>index</h1>", string(data))
	})

	t.Run("should compose as overlay layers sharing a path space", func(t *testing.T) {
		overlayFS := NewOverlayFS(
			MustSub(base, "templates"),
			MustSub(base, "templates/contrib/mine"),
		)

		data, err := overlayFS.ReadFile("index.html")
		require.NoError(t, err)
		assert.Equal(t, "<h1>mine</h1>", string(data))
	})

	t.Run("should panic on a directory that is not a valid path", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = MustSub(base, "/absolute")
		})
	})
}
