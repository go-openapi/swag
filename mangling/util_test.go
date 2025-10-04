// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package mangling

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

type translationSample struct {
	str, out string
}

func titleize(s string) string { return strings.ToTitle(s[:1]) + lower(s[1:]) }

func TestIsEqualFoldIgnoreSpace(t *testing.T) {
	t.Run("should find equal", func(t *testing.T) {
		require.True(t, isEqualFoldIgnoreSpace([]rune(""), ""))
		require.True(t, isEqualFoldIgnoreSpace([]rune(""), "  "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), " a"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "a "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), " a "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "\ta\t"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "a"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("A"), "\u00A0a\u00A0"))

		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), "ab "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), "AB "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB "))

		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), " à"))
		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), "à "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), " à "))
		require.True(t, isEqualFoldIgnoreSpace([]rune("À"), " à "))
	})

	t.Run("should find different", func(t *testing.T) {
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))

		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " A B "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " a b "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB \u00A0\u00A0x"))
		require.False(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB \u00A0\u00A0é"))

		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))

		require.False(t, isEqualFoldIgnoreSpace([]rune("A"), " à"))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), " bà"))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), "àb "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), " a "))
		require.False(t, isEqualFoldIgnoreSpace([]rune("À"), "Á"))
	})
}
