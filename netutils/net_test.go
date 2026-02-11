// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package netutils

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestSplitHostPort(t *testing.T) {
	data := []struct {
		Input string
		Host  string
		Port  int
		Err   bool
	}{
		{"localhost:3933", "localhost", 3933, false},
		{"localhost:yellow", "", -1, true},
		{"localhost", "", -1, true},
		{"localhost:", "", -1, true},
		{"localhost:3933", "localhost", 3933, false},
	}

	for _, e := range data {
		h, p, err := SplitHostPort(e.Input)
		if !e.Err {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}

		assert.EqualT(t, e.Host, h)
		assert.EqualT(t, e.Port, p)
	}
}
