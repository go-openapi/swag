// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package loading

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestRootRelativeVolumeMismatch(t *testing.T) {
	// A path on a different volume cannot be expressed relative to the root. filepath.Rel
	// returns an error, which rootRelative must propagate so the read is rejected rather than
	// silently escaping the root.
	_, err := rootRelative(`C:\root`, `D:\secret.txt`)
	require.Error(t, err)
}
