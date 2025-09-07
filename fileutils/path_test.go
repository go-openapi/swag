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

package fileutils

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
