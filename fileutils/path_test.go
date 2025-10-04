// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func makeDirStructure(tb testing.TB, tgt string) (string, string) {
	_ = tgt

	tb.Helper()

	td := tb.TempDir()
	realPath := filepath.Join(td, "src", "foo", "bar")
	err := os.MkdirAll(realPath, os.ModePerm)
	require.NoError(tb, err)
	linkPathBase := filepath.Join(td, "src", "baz")
	err = os.MkdirAll(linkPathBase, os.ModePerm)
	require.NoError(tb, err)
	linkPath := filepath.Join(linkPathBase, "das")
	err = os.Symlink(realPath, linkPath)
	require.NoError(tb, err)

	td2 := tb.TempDir()
	realPath = filepath.Join(td2, "src", "fuu", "bir")
	err = os.MkdirAll(realPath, os.ModePerm)
	require.NoError(tb, err)
	linkPathBase = filepath.Join(td2, "src", "biz")
	err = os.MkdirAll(linkPathBase, os.ModePerm)
	require.NoError(tb, err)
	linkPath = filepath.Join(linkPathBase, "dis")
	err = os.Symlink(realPath, linkPath)
	require.NoError(tb, err)
	return td, td2
}

func TestFindPackage(t *testing.T) {
	pth, pth2 := makeDirStructure(t, "")

	searchPath := pth + string(filepath.ListSeparator) + pth2
	// finds package when real name mentioned
	pkg := FindInSearchPath(searchPath, "foo/bar")
	assert.NotEmpty(t, pkg)
	assertPath(t, path.Join(pth, "src", "foo", "bar"), pkg)
	// finds package when real name is mentioned in secondary
	pkg = FindInSearchPath(searchPath, "fuu/bir")
	assert.NotEmpty(t, pkg)
	assertPath(t, path.Join(pth2, "src", "fuu", "bir"), pkg)
	// finds package when symlinked
	pkg = FindInSearchPath(searchPath, "baz/das")
	assert.NotEmpty(t, pkg)
	assertPath(t, path.Join(pth, "src", "foo", "bar"), pkg)
	// finds package when symlinked in secondary
	pkg = FindInSearchPath(searchPath, "biz/dis")
	assert.NotEmpty(t, pkg)
	assertPath(t, path.Join(pth2, "src", "fuu", "bir"), pkg)
	// return empty string when nothing is found
	pkg = FindInSearchPath(searchPath, "not/there")
	assert.Empty(t, pkg)
}

//nolint:unparam
func assertPath(t testing.TB, expected, actual string) bool {
	fp, err := filepath.EvalSymlinks(expected)
	require.NoError(t, err)

	return assert.Equal(t, fp, actual)
}

func TestFullGOPATH(t *testing.T) {
	os.Unsetenv(GOPATHKey)
	ngp := "/some/where:/other/place"
	t.Setenv(GOPATHKey, ngp)

	expected := ngp + ":" + runtime.GOROOT() //nolint: staticcheck // this is a deprecated function
	assert.Equal(t, expected, FullGoSearchPath())
}
