// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// opaqueFixture returns a base file system, and an overlay meant to entirely own "cfg".
func opaqueFixture() (base, overlay fstest.MapFS) {
	base = fstest.MapFS{
		"root.txt":          &fstest.MapFile{Data: []byte("root from base")},
		"cfg/a.txt":         &fstest.MapFile{Data: []byte("cfg a from base")},
		"cfg/b.txt":         &fstest.MapFile{Data: []byte("cfg b from base")},
		"cfg/sub/deep.txt":  &fstest.MapFile{Data: []byte("deep from base")},
		"other/keep-me.txt": &fstest.MapFile{Data: []byte("other from base")},
	}
	overlay = fstest.MapFS{
		"cfg/a.txt":   &fstest.MapFile{Data: []byte("cfg a from overlay")},
		"other.txt":   &fstest.MapFile{Data: []byte("other from overlay")},
		"cfg/new.txt": &fstest.MapFile{Data: []byte("cfg new from overlay")},
	}

	return base, overlay
}

func TestOpaqueFS(t *testing.T) {
	_, overlay := opaqueFixture()

	t.Run("should implement OpaqueDirFS", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay, "cfg")

		assert.Implements(t, new(fs.FS), opaqueFS)
		assert.Implements(t, new(OpaqueDirFS), opaqueFS)
		assert.Implements(t, new(fs.ReadFileFS), opaqueFS)
		assert.Implements(t, new(fs.StatFS), opaqueFS)
		assert.Implements(t, new(fs.ReadDirFS), opaqueFS)
	})

	t.Run("should report the declared directories", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay, "cfg")

		assert.TrueT(t, opaqueFS.IsOpaqueDir("cfg"))
		assert.FalseT(t, opaqueFS.IsOpaqueDir("other"))
		assert.FalseT(t, opaqueFS.IsOpaqueDir("cfg/sub"))
	})

	t.Run("should clean the declared directories", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay, "./cfg/", "a/../b")

		assert.TrueT(t, opaqueFS.IsOpaqueDir("cfg"))
		assert.TrueT(t, opaqueFS.IsOpaqueDir("b"))
	})

	t.Run("should declare every directory but the root with OpaqueDirsAll", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay, OpaqueDirsAll)

		assert.TrueT(t, opaqueFS.IsOpaqueDir("cfg"))
		assert.TrueT(t, opaqueFS.IsOpaqueDir("cfg/sub"))
		assert.TrueT(t, opaqueFS.IsOpaqueDir("anything"))
		assert.FalseT(t, opaqueFS.IsOpaqueDir("."))
	})

	t.Run("should declare the root alongside OpaqueDirsAll", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay, OpaqueDirsAll, ".")

		assert.TrueT(t, opaqueFS.IsOpaqueDir("."))
		assert.TrueT(t, opaqueFS.IsOpaqueDir("cfg"))
	})

	t.Run("should declare nothing by default", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay)

		assert.FalseT(t, opaqueFS.IsOpaqueDir("cfg"))
		assert.FalseT(t, opaqueFS.IsOpaqueDir("."))
	})

	t.Run("should delegate to the wrapped file system", func(t *testing.T) {
		opaqueFS := NewOpaqueFS(overlay, "cfg")

		data, err := opaqueFS.ReadFile("cfg/a.txt")
		require.NoError(t, err)
		assert.EqualT(t, "cfg a from overlay", string(data))

		info, err := opaqueFS.Stat("cfg")
		require.NoError(t, err)
		assert.TrueT(t, info.IsDir())

		assert.SliceEqualT(t, []string{"a.txt", "new.txt"}, dirNames(t, opaqueFS, "cfg"))

		file, err := opaqueFS.Open("other.txt")
		require.NoError(t, err)
		require.NoError(t, file.Close())
	})
}

func TestOverlayFSOpaqueDir(t *testing.T) {
	base, overlay := opaqueFixture()
	overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, "cfg"))

	t.Run("should replace the content of the owned directory", func(t *testing.T) {
		assert.SliceEqualT(t, []string{"a.txt", "new.txt"}, dirNames(t, overlayFS, "cfg"))
	})

	t.Run("should keep merging the directories that are not owned", func(t *testing.T) {
		assert.SliceEqualT(t,
			[]string{"cfg", "other", "other.txt", "root.txt"},
			dirNames(t, overlayFS, "."),
		)
	})

	t.Run("should resolve a file of the owned directory", func(t *testing.T) {
		data, err := overlayFS.ReadFile("cfg/a.txt")
		require.NoError(t, err)
		assert.EqualT(t, "cfg a from overlay", string(data))
	})

	t.Run("should hide a file that the base holds under the owned directory", func(t *testing.T) {
		// the whole point: "cfg/b.txt" is neither listed nor readable
		for _, name := range []string{"cfg/b.txt", "cfg/sub", "cfg/sub/deep.txt"} {
			_, err := overlayFS.ReadFile(name)
			require.Error(t, err, name)
			assert.ErrorIs(t, err, fs.ErrNotExist, name)

			_, err = overlayFS.Open(name)
			require.Error(t, err, name)
			assert.ErrorIs(t, err, fs.ErrNotExist, name)

			_, err = overlayFS.Stat(name)
			require.Error(t, err, name)
			assert.ErrorIs(t, err, fs.ErrNotExist, name)
		}
	})

	t.Run("should hide a directory that the base holds under the owned directory", func(t *testing.T) {
		_, err := overlayFS.ReadDir("cfg/sub")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("should keep resolving outside of the owned directory", func(t *testing.T) {
		data, err := overlayFS.ReadFile("other/keep-me.txt")
		require.NoError(t, err)
		assert.EqualT(t, "other from base", string(data))

		data, err = overlayFS.ReadFile("root.txt")
		require.NoError(t, err)
		assert.EqualT(t, "root from base", string(data))
	})

	t.Run("should walk the owned directory without its lower layers", func(t *testing.T) {
		// WalkDir descends "other" before visiting "other.txt": the order is
		// lexical within each directory, not over the whole path
		assert.SliceEqualT(t, []string{
			"cfg/a.txt",
			"cfg/new.txt",
			"other/keep-me.txt",
			"other.txt",
			"root.txt",
		}, walkFiles(t, overlayFS))
	})

	t.Run("should conform to fs.FS", func(t *testing.T) {
		require.NoError(t, fstest.TestFS(overlayFS,
			"root.txt", "other.txt", "other/keep-me.txt", "cfg/a.txt", "cfg/new.txt",
		))
	})
}

func TestOverlayFSOpaqueDirEdgeCases(t *testing.T) {
	base, overlay := opaqueFixture()

	t.Run("should ignore a declaration that the layer doesn't hold", func(t *testing.T) {
		// the overlay claims "other", which it does not hold, so the base still resolves
		overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, "other"))

		assert.SliceEqualT(t, []string{"keep-me.txt"}, dirNames(t, overlayFS, "other"))

		data, err := overlayFS.ReadFile("other/keep-me.txt")
		require.NoError(t, err)
		assert.EqualT(t, "other from base", string(data))
	})

	t.Run("should ignore a declaration naming a regular file", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, "other.txt"))

		data, err := overlayFS.ReadFile("root.txt")
		require.NoError(t, err)
		assert.EqualT(t, "root from base", string(data))
	})

	t.Run("should own the whole file system when the root is declared", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, "."))

		assert.SliceEqualT(t, []string{"cfg", "other.txt"}, dirNames(t, overlayFS, "."))

		// owning the root owns every directory below it, not just its own entries
		assert.SliceEqualT(t, []string{"a.txt", "new.txt"}, dirNames(t, overlayFS, "cfg"))

		for _, name := range []string{"root.txt", "cfg/b.txt", "other/keep-me.txt"} {
			_, err := overlayFS.ReadFile(name)
			require.Error(t, err, name)
			assert.ErrorIs(t, err, fs.ErrNotExist, name)
		}

		require.NoError(t, fstest.TestFS(overlayFS, "other.txt", "cfg/a.txt", "cfg/new.txt"))
	})

	t.Run("should let a layer above the owner keep contributing", func(t *testing.T) {
		// the owner sits in the middle: what is stacked above it still merges in
		top := fstest.MapFS{"cfg/top.txt": &fstest.MapFile{Data: []byte("cfg top")}}
		overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, "cfg"), top)

		assert.SliceEqualT(t, []string{"a.txt", "new.txt", "top.txt"}, dirNames(t, overlayFS, "cfg"))

		// the base is still hidden by the owner below the top layer
		_, err := overlayFS.ReadFile("cfg/b.txt")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("should own a nested directory alone", func(t *testing.T) {
		nested := fstest.MapFS{"cfg/sub/only.txt": &fstest.MapFile{Data: []byte("only")}}
		overlayFS := NewOverlayFS(base, NewOpaqueFS(nested, "cfg/sub"))

		assert.SliceEqualT(t, []string{"only.txt"}, dirNames(t, overlayFS, "cfg/sub"))
		// "cfg" itself is not owned, so it still merges
		assert.SliceEqualT(t, []string{"a.txt", "b.txt", "sub"}, dirNames(t, overlayFS, "cfg"))
	})

	t.Run("with OpaqueDirsAll", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, OpaqueDirsAll))

		t.Run("should still merge the root", func(t *testing.T) {
			assert.SliceEqualT(t,
				[]string{"cfg", "other", "other.txt", "root.txt"},
				dirNames(t, overlayFS, "."),
			)
		})

		t.Run("should own every directory that the layer holds", func(t *testing.T) {
			assert.SliceEqualT(t, []string{"a.txt", "new.txt"}, dirNames(t, overlayFS, "cfg"))

			_, err := overlayFS.ReadFile("cfg/b.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})

		t.Run("should leave the directories that the layer doesn't hold alone", func(t *testing.T) {
			assert.SliceEqualT(t, []string{"keep-me.txt"}, dirNames(t, overlayFS, "other"))

			data, err := overlayFS.ReadFile("other/keep-me.txt")
			require.NoError(t, err)
			assert.EqualT(t, "other from base", string(data))
		})

		t.Run("should keep resolving the files held beside the root", func(t *testing.T) {
			data, err := overlayFS.ReadFile("root.txt")
			require.NoError(t, err)
			assert.EqualT(t, "root from base", string(data))
		})

		t.Run("should conform to fs.FS", func(t *testing.T) {
			require.NoError(t, fstest.TestFS(overlayFS,
				"root.txt", "other.txt", "other/keep-me.txt", "cfg/a.txt", "cfg/new.txt",
			))
		})
	})

	t.Run("should terminate on a malformed name", func(t *testing.T) {
		// "/" is its own parent, so walking up the directories of an absolute
		// name must not loop forever
		overlayFS := NewOverlayFS(base, NewOpaqueFS(overlay, "cfg"))

		for _, name := range []string{"/cfg", "/", "./cfg", "cfg/", ""} {
			_, err := overlayFS.Open(name)
			require.Error(t, err, name)
		}
	})
}

func TestOpaqueDirNotHeld(t *testing.T) {
	base := fstest.MapFS{
		"dir/a.txt": {Data: []byte("base")},
		"top.txt":   {Data: []byte("top")},
	}
	// the upper layer declares "dir" opaque, and holds nothing under it
	upper := NewOpaqueFS(fstest.MapFS{"other.txt": {Data: []byte("upper")}}, "dir")

	t.Run("should report the declaration, held or not", func(t *testing.T) {
		assert.True(t, upper.IsOpaqueDir("dir"))
	})

	t.Run("should leave the layers below resolving that directory", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, upper)

		entries, err := fs.ReadDir(overlayFS, "dir")
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, "a.txt", entries[0].Name())

		content, err := fs.ReadFile(overlayFS, "dir/a.txt")
		require.NoError(t, err)
		assert.Equal(t, "base", string(content))
	})
}
