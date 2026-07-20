// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package loading

import (
	"path/filepath"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestRootRelativeVolumeMismatch(t *testing.T) {
	// A path on a different volume cannot be expressed relative to the root. filepath.Rel
	// returns an error, which rootRelative must propagate so the read is rejected rather than
	// silently escaping the root.
	_, err := rootRelative(`C:\root`, `D:\secret.txt`)
	require.Error(t, err)
}

func TestRootRelativeWindowsURLPath(t *testing.T) {
	// go-openapi/spec normalizes Windows paths to a file-URL form with a leading slash before the
	// drive letter, e.g. "/D:/dir/child.json". WithRoot must rebase such an in-root path onto the
	// root rather than reject it for not looking absolute.
	root := t.TempDir()
	urlForm := "/" + filepath.ToSlash(filepath.Join(root, "child.json")) // "/D:/.../child.json"

	rel, err := rootRelative(root, urlForm)
	require.NoError(t, err)
	assert.EqualT(t, "child.json", rel)
}
