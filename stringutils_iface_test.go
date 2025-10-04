// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestStringUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		assert.True(t, ContainsStrings([]string{"a", "b"}, "a"))
		assert.True(t, ContainsStringsCI([]string{"a", "b"}, "A"))
		require.Len(t, JoinByFormat([]string{"a", "b"}, "pipes"), 1)
		require.Len(t, SplitByFormat("a|b", "pipes"), 2)
	})
}
