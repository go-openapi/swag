// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// errFS is a layer that always fails with an error which is not a "not found",
// so that [OverlayFS] stops resolving and reports it.
type errFS struct{}

func (errFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

// wrappingFS adds a layer of wrapping to the errors of the underlying file system,
// so that they are no longer directly a [fs.PathError].
type wrappingFS struct {
	inner fs.FS
}

func (f wrappingFS) Open(name string) (fs.File, error) {
	file, err := f.inner.Open(name)
	if err != nil {
		return nil, fmt.Errorf("wrappingFS: %w", err)
	}

	return file, nil
}

// statErrFS is like [errFS], but it also implements [fs.StatFS].
type statErrFS struct{ errFS }

func (statErrFS) Stat(name string) (fs.FileInfo, error) {
	return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrPermission}
}

// badStatFS is a layer whose files cannot be stat'ed.
type badStatFS struct{}

func (badStatFS) Open(string) (fs.File, error) { return badStatFile{}, nil }

type badStatFile struct{}

func (badStatFile) Read([]byte) (int, error) { return 0, io.EOF }

func (badStatFile) Close() error { return nil }

func (badStatFile) Stat() (fs.FileInfo, error) {
	return nil, &fs.PathError{Op: "stat", Path: "unknown", Err: fs.ErrPermission}
}

// unlistableDirFS is a layer that reports a directory for every name, but cannot list it.
type unlistableDirFS struct{}

func (unlistableDirFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

func (unlistableDirFS) Stat(name string) (fs.FileInfo, error) { return dirInfo{name: name}, nil }

// dirInfo describes a directory, and nothing else.
type dirInfo struct{ name string }

func (i dirInfo) Name() string     { return i.name }
func (dirInfo) Size() int64        { return 0 }
func (dirInfo) Mode() fs.FileMode  { return fs.ModeDir | 0o500 }
func (dirInfo) ModTime() time.Time { return time.Time{} }
func (dirInfo) IsDir() bool        { return true }
func (dirInfo) Sys() any           { return nil }

// entryNames maps directory entries to their names, preserving their order.
func entryNames(entries []fs.DirEntry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names
}

// dirNames returns the names listed by ReadDir, in the order reported by the file system.
func dirNames(tb testing.TB, fsys fs.ReadDirFS, name string) []string {
	tb.Helper()

	entries, err := fsys.ReadDir(name)
	require.NoError(tb, err)

	return entryNames(entries)
}

// overlayFixture returns a base file system and two overlays, with a few
// entries in common so that layer precedence can be observed.
func overlayFixture() (base, lower, upper fstest.MapFS) {
	base = fstest.MapFS{
		"shared.txt":       &fstest.MapFile{Data: []byte("from base")},
		"only-in-base.txt": &fstest.MapFile{Data: []byte("only in base")},
		"dir/in-base.txt":  &fstest.MapFile{Data: []byte("dir entry from base")},
		"dir/both.txt":     &fstest.MapFile{Data: []byte("dir entry from base")},
	}
	lower = fstest.MapFS{
		"shared.txt":        &fstest.MapFile{Data: []byte("from lower")},
		"only-in-lower.txt": &fstest.MapFile{Data: []byte("only in lower")},
	}
	upper = fstest.MapFS{
		"shared.txt":        &fstest.MapFile{Data: []byte("from upper")},
		"only-in-upper.txt": &fstest.MapFile{Data: []byte("only in upper")},
		"dir/in-upper.txt":  &fstest.MapFile{Data: []byte("dir entry from upper")},
		"dir/both.txt":      &fstest.MapFile{Data: []byte("dir entry from upper")},
	}

	return base, lower, upper
}

func TestOverlayFS(t *testing.T) {
	base, lower, upper := overlayFixture()

	t.Run("should implement the read-only file system interfaces, but not fs.GlobFS", func(t *testing.T) {
		overlayFS := NewOverlayFS(base)

		assert.Implements(t, new(fs.FS), overlayFS)
		assert.Implements(t, new(fs.ReadFileFS), overlayFS)
		assert.Implements(t, new(fs.StatFS), overlayFS)
		assert.Implements(t, new(fs.ReadDirFS), overlayFS)
		assert.NotImplements(t, new(fs.GlobFS), overlayFS)
	})

	t.Run("with no overlay", func(t *testing.T) {
		overlayFS := NewOverlayFS(base)

		t.Run("should resolve against the base", func(t *testing.T) {
			data, err := overlayFS.ReadFile("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "from base", string(data))
		})
	})

	t.Run("with a single overlay", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower)

		t.Run("should resolve the overlay before the base", func(t *testing.T) {
			data, err := overlayFS.ReadFile("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "from lower", string(data))
		})

		t.Run("should fall through to the base", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-base.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in base", string(data))
		})
	})

	t.Run("with stacked overlays", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		t.Run("should resolve the last overlay first", func(t *testing.T) {
			// the overlays are stacked in reverse order of declaration:
			// the last one provided sits on top.
			data, err := overlayFS.ReadFile("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "from upper", string(data))
		})

		t.Run("should resolve a file found in an intermediate layer", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-lower.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in lower", string(data))
		})

		t.Run("should resolve a file found in the base only", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-base.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in base", string(data))
		})
	})

}

func TestOverlayFSMethods(t *testing.T) {
	base, lower, upper := overlayFixture()

	t.Run("with Open", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		t.Run("should open a file from the topmost layer that holds it", func(t *testing.T) {
			file, err := overlayFS.Open("shared.txt")
			require.NoError(t, err)
			t.Cleanup(func() { _ = file.Close() })

			data, err := io.ReadAll(file)
			require.NoError(t, err)
			assert.EqualT(t, "from upper", string(data))
		})

		t.Run("should open a file from the base", func(t *testing.T) {
			file, err := overlayFS.Open("only-in-base.txt")
			require.NoError(t, err)
			require.NoError(t, file.Close())
		})

		t.Run("should not open a file absent from all layers", func(t *testing.T) {
			_, err := overlayFS.Open("nowhere.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

	t.Run("with Stat", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		t.Run("should stat a file from the topmost layer that holds it", func(t *testing.T) {
			info, err := overlayFS.Stat("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "shared.txt", info.Name())
			assert.EqualT(t, int64(len("from upper")), info.Size())
		})

		t.Run("should stat a directory", func(t *testing.T) {
			info, err := overlayFS.Stat("dir")
			require.NoError(t, err)
			assert.TrueT(t, info.IsDir())
		})

		t.Run("should not stat a file absent from all layers", func(t *testing.T) {
			_, err := overlayFS.Stat("nowhere.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

	t.Run("with ReadDir", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		t.Run("should merge the entries of every layer holding the directory", func(t *testing.T) {
			// "dir" is held by both the base and the upper overlay
			assert.SliceEqualT(t,
				[]string{"both.txt", "in-base.txt", "in-upper.txt"},
				dirNames(t, overlayFS, "dir"),
			)
		})

		t.Run("should not list a directory absent from all layers", func(t *testing.T) {
			_, err := overlayFS.ReadDir("nowhere")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

}

func TestOverlayFSLayerKinds(t *testing.T) {
	base, lower, upper := overlayFixture()

	t.Run("with layers that only implement fs.FS", func(t *testing.T) {
		// exercises the Open-based fallback used when a layer does not provide Stat
		overlayFS := NewOverlayFS(
			openOnlyFS{inner: base},
			openOnlyFS{inner: lower}, openOnlyFS{inner: upper},
		)

		t.Run("should resolve with ReadFile", func(t *testing.T) {
			data, err := overlayFS.ReadFile("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "from upper", string(data))
		})

		t.Run("should resolve with Stat", func(t *testing.T) {
			info, err := overlayFS.Stat("only-in-lower.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only-in-lower.txt", info.Name())
		})

		t.Run("should resolve with ReadDir", func(t *testing.T) {
			assert.SliceEqualT(t,
				[]string{"both.txt", "in-base.txt", "in-upper.txt"},
				dirNames(t, overlayFS, "dir"),
			)
		})

		t.Run("should fall through to the base", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-base.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in base", string(data))
		})
	})

	t.Run("with a layer that wraps its errors", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, wrappingFS{inner: upper})

		t.Run("should resolve the wrapping layer", func(t *testing.T) {
			data, err := overlayFS.ReadFile("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "from upper", string(data))
		})

		t.Run("should fall through to the base", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-base.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in base", string(data))
		})
	})

	t.Run("with a layer that is itself an OverlayFS", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, NewOverlayFS(lower, upper))

		t.Run("should resolve the topmost layer of the nested overlay", func(t *testing.T) {
			data, err := overlayFS.ReadFile("shared.txt")
			require.NoError(t, err)
			assert.EqualT(t, "from upper", string(data))
		})

		t.Run("should resolve the base of the nested overlay", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-lower.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in lower", string(data))
		})

		t.Run("should fall through to the outer base", func(t *testing.T) {
			data, err := overlayFS.ReadFile("only-in-base.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in base", string(data))
		})

		t.Run("should not resolve a file absent from all layers", func(t *testing.T) {
			_, err := overlayFS.ReadFile("nowhere.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrNotExist)
		})
	})

}

func TestOverlayFSLayerInError(t *testing.T) {
	base, _, upper := overlayFixture()

	t.Run("with a layer in error", func(t *testing.T) {
		t.Run("should report an Open error rather than fall through", func(t *testing.T) {
			overlayFS := NewOverlayFS(base, errFS{})

			_, err := overlayFS.Open("only-in-base.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrPermission)
		})

		t.Run("should report a Stat error rather than fall through", func(t *testing.T) {
			overlayFS := NewOverlayFS(base, statErrFS{})

			_, err := overlayFS.Stat("only-in-base.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrPermission)

			_, err = overlayFS.ReadFile("only-in-base.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrPermission)

			_, err = overlayFS.ReadDir("dir")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrPermission)
		})

		t.Run("should report an error raised by a layer below the resolved one", func(t *testing.T) {
			// the failing layer is never reached, since the upper one resolves first
			overlayFS := NewOverlayFS(errFS{}, upper)

			data, err := overlayFS.ReadFile("only-in-upper.txt")
			require.NoError(t, err)
			assert.EqualT(t, "only in upper", string(data))

			_, err = overlayFS.ReadFile("only-in-base.txt")
			require.Error(t, err)
			assert.ErrorIs(t, err, fs.ErrPermission)
		})
	})
}

func TestOverlayFSMergeDirs(t *testing.T) {
	base, lower, upper := overlayFixture()

	t.Run("should merge directories by default", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		t.Run("should report every layer, sorted by file name", func(t *testing.T) {
			assert.SliceEqualT(t,
				[]string{"both.txt", "in-base.txt", "in-upper.txt"},
				dirNames(t, overlayFS, "dir"),
			)
		})

		t.Run("should report a name held by several layers only once", func(t *testing.T) {
			data, err := overlayFS.ReadFile("dir/both.txt")
			require.NoError(t, err)
			assert.EqualT(t, "dir entry from upper", string(data),
				"expected the topmost layer to win",
			)
		})

		t.Run("should merge the root directory", func(t *testing.T) {
			assert.SliceEqualT(t, []string{
				"dir", "only-in-base.txt", "only-in-lower.txt", "only-in-upper.txt", "shared.txt",
			}, dirNames(t, overlayFS, "."))
		})
	})

	t.Run("should open a directory as a fs.ReadDirFile", func(t *testing.T) {
		for _, toPin := range []struct {
			name     string
			overlays []fs.FS
		}{
			{name: "merged", overlays: []fs.FS{lower, upper}},
			{name: "opaque", overlays: []fs.FS{lower, NewOpaqueFS(upper, "dir")}},
		} {
			tc := toPin

			t.Run("when "+tc.name, func(t *testing.T) {
				overlayFS := NewOverlayFS(base, tc.overlays...)

				file, err := overlayFS.Open("dir")
				require.NoError(t, err)
				t.Cleanup(func() { _ = file.Close() })

				dir, isDir := file.(fs.ReadDirFile)
				require.TrueT(t, isDir)

				entries, err := dir.ReadDir(-1)
				require.NoError(t, err)

				t.Run("should report the same entries as ReadDir", func(t *testing.T) {
					assert.SliceEqualT(t, dirNames(t, overlayFS, "dir"), entryNames(entries))
				})
			})
		}
	})

	t.Run("with a merged directory read in chunks", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		file, err := overlayFS.Open("dir")
		require.NoError(t, err)
		t.Cleanup(func() { _ = file.Close() })

		dir, isDir := file.(fs.ReadDirFile)
		require.TrueT(t, isDir)

		t.Run("should report the entries chunk by chunk", func(t *testing.T) {
			var names []string
			for {
				entries, err := dir.ReadDir(2)
				if errors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				names = append(names, entryNames(entries)...)
			}

			assert.SliceEqualT(t, []string{"both.txt", "in-base.txt", "in-upper.txt"}, names)
		})

		t.Run("should not read a directory as a regular file", func(t *testing.T) {
			_, err := file.Read(make([]byte, 1))
			require.Error(t, err)

			var pathErr *fs.PathError
			require.ErrorAs(t, err, &pathErr)
			assert.EqualT(t, "dir", pathErr.Path)
		})
	})

	t.Run("should let a file shadow a directory held by a lower layer", func(t *testing.T) {
		// "dir" is a regular file in the topmost layer, so it hides the directory below
		shadowing := fstest.MapFS{"dir": &fstest.MapFile{Data: []byte("not a directory")}}
		overlayFS := NewOverlayFS(base, shadowing)

		data, err := overlayFS.ReadFile("dir")
		require.NoError(t, err)
		assert.EqualT(t, "not a directory", string(data))

		_, err = overlayFS.ReadDir("dir")
		require.Error(t, err)
	})

	t.Run("should stop merging below a file held by an intermediate layer", func(t *testing.T) {
		// the topmost layer holds "dir" as a directory, an intermediate layer holds it
		// as a regular file, so the base contributes nothing to the merge
		shadowing := fstest.MapFS{"dir": &fstest.MapFile{Data: []byte("not a directory")}}
		overlayFS := NewOverlayFS(base, shadowing, upper)

		assert.SliceEqualT(t, []string{"both.txt", "in-upper.txt"}, dirNames(t, overlayFS, "dir"))
	})

	t.Run("should not merge a directory absent from all layers", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		_, err := overlayFS.ReadDir("nowhere")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("should report an error raised by a layer while merging", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, statErrFS{})

		_, err := overlayFS.ReadDir("dir")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrPermission)
	})
}

func TestOverlayFSMergeDirsErrors(t *testing.T) {
	base, _, upper := overlayFixture()

	t.Run("should report a file that cannot be stat'ed when opening", func(t *testing.T) {
		overlayFS := NewOverlayFS(badStatFS{})

		_, err := overlayFS.Open("anything.txt")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrPermission)
	})

	t.Run("should report a layer in error when opening a directory", func(t *testing.T) {
		// the topmost layer resolves "dir", the one below it fails while merging
		overlayFS := NewOverlayFS(statErrFS{}, upper)

		_, err := overlayFS.Open("dir")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrPermission)
	})

	t.Run("should report a directory that cannot be listed", func(t *testing.T) {
		overlayFS := NewOverlayFS(unlistableDirFS{})

		_, err := overlayFS.ReadDir("dir")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrPermission)
	})

	t.Run("should report a directory owned by a layer that cannot list it", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, NewOpaqueFS(unlistableDirFS{}, "dir"))

		_, err := overlayFS.ReadDir("dir")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrPermission)
	})
}

func TestOverlayFSNotFoundError(t *testing.T) {
	base, _, _ := overlayFixture()
	overlayFS := NewOverlayFS(base)

	for _, toPin := range []struct {
		method string
		op     string
		call   func(string) error
	}{
		{
			method: "Open",
			op:     "open",
			call:   func(name string) error { _, err := overlayFS.Open(name); return err },
		},
		{
			method: "ReadFile",
			op:     "open",
			call:   func(name string) error { _, err := overlayFS.ReadFile(name); return err },
		},
		{
			method: "Stat",
			op:     "stat",
			call:   func(name string) error { _, err := overlayFS.Stat(name); return err },
		},
		{
			method: "ReadDir",
			op:     "readdir",
			call:   func(name string) error { _, err := overlayFS.ReadDir(name); return err },
		},
	} {
		tc := toPin

		t.Run("with "+tc.method, func(t *testing.T) {
			t.Run("should report the missing name as a fs.PathError", func(t *testing.T) {
				err := tc.call("nowhere.txt")
				require.Error(t, err)
				assert.ErrorIs(t, err, fs.ErrNotExist)

				var pathErr *fs.PathError
				require.ErrorAs(t, err, &pathErr)
				assert.EqualT(t, "nowhere.txt", pathErr.Path)
				assert.EqualT(t, tc.op, pathErr.Op)
			})
		})
	}
}

func TestOverlayFSConformance(t *testing.T) {
	base, lower, upper := overlayFixture()

	t.Run("should conform to fs.FS when there is no overlay", func(t *testing.T) {
		require.NoError(t,
			fstest.TestFS(NewOverlayFS(base), "shared.txt", "only-in-base.txt", "dir/in-base.txt"),
		)
	})

	t.Run("should conform to fs.FS with merged directories", func(t *testing.T) {
		require.NoError(t,
			fstest.TestFS(NewOverlayFS(base, lower, upper),
				"shared.txt", "only-in-base.txt", "only-in-lower.txt", "only-in-upper.txt",
				"dir/in-base.txt", "dir/in-upper.txt", "dir/both.txt",
			),
		)
	})

	t.Run("should conform to fs.FS with an opaque directory", func(t *testing.T) {
		require.NoError(t,
			fstest.TestFS(NewOverlayFS(base, lower, NewOpaqueFS(upper, "dir")),
				"shared.txt", "only-in-base.txt", "only-in-lower.txt", "only-in-upper.txt",
				"dir/in-upper.txt", "dir/both.txt",
			),
		)
	})

	t.Run("should walk the entries of every layer", func(t *testing.T) {
		overlayFS := NewOverlayFS(base, lower, upper)

		assert.SliceEqualT(t, []string{
			"dir/both.txt",
			"dir/in-base.txt",
			"dir/in-upper.txt",
			"only-in-base.txt",
			"only-in-lower.txt",
			"only-in-upper.txt",
			"shared.txt",
		}, walkFiles(t, overlayFS))
	})

	t.Run("should walk an opaque directory without its lower layers", func(t *testing.T) {
		// only "dir" is claimed: the rest of the tree still merges
		overlayFS := NewOverlayFS(base, lower, NewOpaqueFS(upper, "dir"))

		assert.SliceEqualT(t, []string{
			"dir/both.txt",
			"dir/in-upper.txt",
			"only-in-base.txt",
			"only-in-lower.txt",
			"only-in-upper.txt",
			"shared.txt",
		}, walkFiles(t, overlayFS))
	})
}

// walkFiles returns the regular files reachable from the root of a file system.
func walkFiles(tb testing.TB, fsys fs.FS) []string {
	tb.Helper()

	var walked []string
	require.NoError(tb, fs.WalkDir(fsys, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			walked = append(walked, name)
		}

		return nil
	}))

	return walked
}

func TestOverlayFSOnDisk(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "a.txt"), []byte("content of a"), 0o600))

	overlayFS := NewOverlayFS(os.DirFS(root), fstest.MapFS{
		"b.txt": &fstest.MapFile{Data: []byte("content of b")},
	})

	t.Run("should resolve a file on disk", func(t *testing.T) {
		data, err := overlayFS.ReadFile("a.txt")
		require.NoError(t, err)
		assert.EqualT(t, "content of a", string(data))
	})

	t.Run("should resolve a file in the overlay", func(t *testing.T) {
		data, err := overlayFS.ReadFile("b.txt")
		require.NoError(t, err)
		assert.EqualT(t, "content of b", string(data))
	})

	t.Run("should treat a path traversing a regular file as not found", func(t *testing.T) {
		// the on-disk layer reports ENOTDIR here, which resolves as a "not found"
		_, err := overlayFS.Open("a.txt/nested.txt")
		require.Error(t, err)
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})
}
