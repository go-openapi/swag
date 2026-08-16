// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils

import (
	"io"
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func makeMapFSFixture(t *testing.T) *MapFS {
	t.Helper()

	mapFS, err := NewMapFS(FromRawMap(map[string][]byte{
		"top.txt":            []byte("content of top"),
		"dir/in-dir.txt":     []byte("content of in-dir"),
		"dir/sub/deep.txt":   []byte("content of deep"),
		"other/in-other.txt": []byte("content of in-other"),
	}))
	require.NoError(t, err)

	return mapFS
}

func TestMapFS(t *testing.T) {
	mapFS := makeMapFSFixture(t)

	t.Run("should conform to fs.FS", func(t *testing.T) {
		require.NoError(t,
			fstest.TestFS(mapFS, "top.txt", "dir/in-dir.txt", "dir/sub/deep.txt", "other/in-other.txt"),
		)
	})

	t.Run("should read a file", func(t *testing.T) {
		data, err := mapFS.ReadFile("dir/sub/deep.txt")
		require.NoError(t, err)
		assert.Equal(t, "content of deep", string(data))
	})

	t.Run("should list a directory, sorted", func(t *testing.T) {
		entries, err := mapFS.ReadDir(".")
		require.NoError(t, err)

		assert.Equal(t, []string{"dir", "other", "top.txt"}, entryNames(entries))
		assert.True(t, entries[0].IsDir())
		assert.False(t, entries[2].IsDir())
	})

	t.Run("should list a directory holding both a file and a directory", func(t *testing.T) {
		entries, err := mapFS.ReadDir("dir")
		require.NoError(t, err)

		assert.Equal(t, []string{"in-dir.txt", "sub"}, entryNames(entries))
	})

	t.Run("should walk every implied directory", func(t *testing.T) {
		var walked []string
		require.NoError(t, fs.WalkDir(mapFS, ".", func(name string, _ fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			walked = append(walked, name)

			return nil
		}))

		assert.Equal(t, []string{
			".", "dir", "dir/in-dir.txt", "dir/sub", "dir/sub/deep.txt",
			"other", "other/in-other.txt", "top.txt",
		}, walked)
	})

	t.Run("should stat a file", func(t *testing.T) {
		info, err := mapFS.Stat("top.txt")
		require.NoError(t, err)

		assert.Equal(t, "top.txt", info.Name())
		assert.Equal(t, int64(len("content of top")), info.Size())
		assert.Equal(t, DefaultFileMode, info.Mode())
		assert.False(t, info.IsDir())
	})

	t.Run("should stat a directory", func(t *testing.T) {
		info, err := mapFS.Stat("dir/sub")
		require.NoError(t, err)

		assert.Equal(t, "sub", info.Name())
		assert.Equal(t, DefaultDirMode, info.Mode())
		assert.True(t, info.IsDir())
		assert.Zero(t, info.Size())
		assert.True(t, info.ModTime().IsZero())
		assert.Nil(t, info.Sys(), "a directory is implied by the file names, so it carries no metadata")
	})
}

func TestMapFSNames(t *testing.T) {
	t.Run("should normalize names", func(t *testing.T) {
		mapFS, err := NewMapFS(FromRawMap(map[string][]byte{
			"/leading-slash.txt": []byte("A"),
			"./dot-slash.txt":    []byte("B"),
			"dir//redundant.txt": []byte("C"),
			"dir/../hoisted.txt": []byte("D"),
		}))
		require.NoError(t, err)

		for name, want := range map[string]string{
			"leading-slash.txt": "A",
			"dot-slash.txt":     "B",
			"dir/redundant.txt": "C",
			"hoisted.txt":       "D",
		} {
			data, err := mapFS.ReadFile(name)
			require.NoErrorf(t, err, "expected %q to resolve", name)
			assert.Equal(t, want, string(data))
		}
	})

	t.Run("should reject a name climbing above the root", func(t *testing.T) {
		_, err := NewMapFS(FromRawMap(map[string][]byte{"../escape.txt": []byte("nope")}))

		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrInvalid)
	})

	t.Run("should reject a name normalizing to the root", func(t *testing.T) {
		_, err := NewMapFS(FromRawMap(map[string][]byte{".": []byte("nope")}))

		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrInvalid)
	})

	t.Run("should reject two names normalizing to the same one", func(t *testing.T) {
		_, err := NewMapFS(FromRawMap(map[string][]byte{
			"dir/file.txt":  []byte("A"),
			"/dir/file.txt": []byte("B"),
		}))

		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrExist)
	})

	t.Run("should reject a name held as a file and as a directory", func(t *testing.T) {
		_, err := NewMapFS(FromRawMap(map[string][]byte{
			"dir":          []byte("a file"),
			"dir/file.txt": []byte("a file in a directory of the same name"),
		}))

		require.Error(t, err)
		assert.ErrorContains(t, err, "both as a file and as a directory")
	})
}

func TestMapFSMetadata(t *testing.T) {
	modTime := time.Date(2026, time.August, 16, 10, 0, 0, 0, time.UTC)

	mapFS, err := NewMapFS(map[string]MapFile{
		"preset.txt":   {Data: []byte("A")},
		"explicit.txt": {Data: []byte("B"), Mode: 0o600, ModTime: modTime, Sys: "meta"},
	})
	require.NoError(t, err)

	t.Run("should preset the mode of a file that declares none", func(t *testing.T) {
		info, err := mapFS.Stat("preset.txt")
		require.NoError(t, err)

		assert.Equal(t, DefaultFileMode, info.Mode())
		assert.True(t, info.ModTime().IsZero())
		assert.Nil(t, info.Sys())
	})

	t.Run("should report the metadata a file declares", func(t *testing.T) {
		info, err := mapFS.Stat("explicit.txt")
		require.NoError(t, err)

		assert.Equal(t, fs.FileMode(0o600), info.Mode())
		assert.Equal(t, modTime, info.ModTime())
		assert.Equal(t, "meta", info.Sys())
	})
}

func TestMapFSIsolation(t *testing.T) {
	mapFS := makeMapFSFixture(t)

	t.Run("should not let a caller mutate the content it read", func(t *testing.T) {
		data, err := mapFS.ReadFile("top.txt")
		require.NoError(t, err)
		data[0] = 'X'

		again, err := mapFS.ReadFile("top.txt")
		require.NoError(t, err)
		assert.Equal(t, "content of top", string(again))
	})

	t.Run("should not let a caller mutate the directory index", func(t *testing.T) {
		entries, err := mapFS.ReadDir(".")
		require.NoError(t, err)
		entries[0] = nil

		again, err := mapFS.ReadDir(".")
		require.NoError(t, err)
		assert.Equal(t, []string{"dir", "other", "top.txt"}, entryNames(again))
	})
}

func TestMapFSEmpty(t *testing.T) {
	mapFS, err := NewMapFS(nil)
	require.NoError(t, err)

	t.Run("should hold a root", func(t *testing.T) {
		entries, err := mapFS.ReadDir(".")
		require.NoError(t, err)
		assert.Empty(t, entries)

		info, err := mapFS.Stat(".")
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("should conform to fs.FS", func(t *testing.T) {
		require.NoError(t, fstest.TestFS(mapFS))
	})
}

func TestMapFSOpen(t *testing.T) {
	mapFS := makeMapFSFixture(t)

	t.Run("should read a file through Open", func(t *testing.T) {
		file, err := mapFS.Open("dir/in-dir.txt")
		require.NoError(t, err)
		t.Cleanup(func() { _ = file.Close() })

		data, err := io.ReadAll(file)
		require.NoError(t, err)
		assert.Equal(t, "content of in-dir", string(data))
	})

	t.Run("should read a directory in pages", func(t *testing.T) {
		file, err := mapFS.Open(".")
		require.NoError(t, err)
		t.Cleanup(func() { _ = file.Close() })

		dir, isDirFile := file.(fs.ReadDirFile)
		require.True(t, isDirFile)

		first, err := dir.ReadDir(2)
		require.NoError(t, err)
		assert.Equal(t, []string{"dir", "other"}, entryNames(first))

		second, err := dir.ReadDir(2)
		require.NoError(t, err)
		assert.Equal(t, []string{"top.txt"}, entryNames(second))

		_, err = dir.ReadDir(2)
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("should refuse to read a directory as a file", func(t *testing.T) {
		file, err := mapFS.Open("dir")
		require.NoError(t, err)
		t.Cleanup(func() { _ = file.Close() })

		_, err = file.Read(make([]byte, 1))
		assert.ErrorContains(t, err, "is a directory")
	})
}

func TestMapFSErrors(t *testing.T) {
	mapFS := makeMapFSFixture(t)

	t.Run("should report a missing name", func(t *testing.T) {
		for _, probe := range []func() error{
			func() error { _, err := mapFS.Open("nowhere.txt"); return err },
			func() error { _, err := mapFS.ReadFile("nowhere.txt"); return err },
			func() error { _, err := mapFS.Stat("nowhere.txt"); return err },
			func() error { _, err := mapFS.ReadDir("nowhere"); return err },
		} {
			err := probe()
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)

			var pathError *fs.PathError
			require.ErrorAs(t, err, &pathError)
			assert.Contains(t, pathError.Path, "nowhere")
		}
	})

	t.Run("should report an invalid name", func(t *testing.T) {
		for _, probe := range []func() error{
			func() error { _, err := mapFS.Open("/absolute"); return err },
			func() error { _, err := mapFS.ReadFile("/absolute"); return err },
			func() error { _, err := mapFS.Stat("/absolute"); return err },
			func() error { _, err := mapFS.ReadDir("/absolute"); return err },
		} {
			assert.ErrorIs(t, probe(), fs.ErrInvalid)
		}
	})

	t.Run("should report a directory read as a file", func(t *testing.T) {
		_, err := mapFS.ReadFile("dir")

		require.Error(t, err)
		assert.ErrorContains(t, err, "is a directory")
	})

	t.Run("should report a file listed as a directory", func(t *testing.T) {
		_, err := mapFS.ReadDir("top.txt")

		require.Error(t, err)
		assert.ErrorContains(t, err, "not a directory")
	})
}

// TestMapFSAsOverlayLayer covers the use case the type is built for: a handful of files held in
// memory, stacked on top of a file system that holds everything else.
func TestMapFSAsOverlayLayer(t *testing.T) {
	base := makeMapFSFixture(t)

	patch, err := NewMapFS(FromRawMap(map[string][]byte{
		"dir/in-dir.txt": []byte("patched"),
		"dir/added.txt":  []byte("added by the patch"),
	}))
	require.NoError(t, err)

	overlayFS := NewOverlayFS(base, patch)

	t.Run("should override a file of the base", func(t *testing.T) {
		data, err := overlayFS.ReadFile("dir/in-dir.txt")
		require.NoError(t, err)
		assert.Equal(t, "patched", string(data))
	})

	t.Run("should leave the rest of the base reachable", func(t *testing.T) {
		data, err := overlayFS.ReadFile("dir/sub/deep.txt")
		require.NoError(t, err)
		assert.Equal(t, "content of deep", string(data))
	})

	t.Run("should merge the sparse layer into the directories of the base", func(t *testing.T) {
		entries, err := overlayFS.ReadDir("dir")
		require.NoError(t, err)

		assert.Equal(t, []string{"added.txt", "in-dir.txt", "sub"}, entryNames(entries))
	})
}

func TestFromRawMap(t *testing.T) {
	files := FromRawMap(map[string][]byte{"a.txt": []byte("A")})

	require.Len(t, files, 1)
	assert.Equal(t, "A", string(files["a.txt"].Data))
	assert.Equal(t, fs.FileMode(0), files["a.txt"].Mode, "the mode is left to NewMapFS to preset")

	assert.Empty(t, FromRawMap(nil))
}

func TestMapFSAliasing(t *testing.T) {
	t.Run("should copy the map, so a later entry is not served", func(t *testing.T) {
		files := map[string]MapFile{"a.txt": {Data: []byte("a")}}

		mapFS, err := NewMapFS(files)
		require.NoError(t, err)

		files["b.txt"] = MapFile{Data: []byte("b")}
		delete(files, "a.txt")

		_, err = mapFS.ReadFile("b.txt")
		require.ErrorIs(t, err, fs.ErrNotExist, "an entry added afterwards is not served")

		content, err := mapFS.ReadFile("a.txt")
		require.NoError(t, err, "an entry removed afterwards is still served")
		assert.Equal(t, "a", string(content))
	})

	t.Run("should serve the content a caller still holds", func(t *testing.T) {
		// the contents are not copied, which the documentation states
		data := []byte("original")

		mapFS, err := NewMapFS(map[string]MapFile{"a.txt": {Data: data}})
		require.NoError(t, err)

		data[0] = 'X'

		content, err := mapFS.ReadFile("a.txt")
		require.NoError(t, err)
		assert.Equal(t, "Xriginal", string(content))
	})

	t.Run("should not let a reader corrupt what the next one gets", func(t *testing.T) {
		mapFS, err := NewMapFS(FromRawMap(map[string][]byte{"a.txt": []byte("a")}))
		require.NoError(t, err)

		first, err := mapFS.ReadFile("a.txt")
		require.NoError(t, err)
		first[0] = 'X'

		second, err := mapFS.ReadFile("a.txt")
		require.NoError(t, err)
		assert.Equal(t, "a", string(second))
	})
}
