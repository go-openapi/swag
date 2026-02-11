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
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune(""), ""))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune(""), "  "))

		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("A"), " a"))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("A"), "a "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("A"), " a "))

		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("A"), "\ta\t"))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("A"), "a"))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("A"), "\u00A0a\u00A0"))

		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab"))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), "ab "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), " ab "))

		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB"))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), "AB "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB "))

		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("À"), " à"))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("À"), "à "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("À"), " à "))
		require.TrueT(t, isEqualFoldIgnoreSpace([]rune("À"), " à "))
	})

	t.Run("should find different", func(t *testing.T) {
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))

		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("AB"), " A B "))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("AB"), " a b "))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB \u00A0\u00A0x"))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("AB"), " AB \u00A0\u00A0é"))

		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), ""))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), " b "))

		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("A"), " à"))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("À"), " bà"))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("À"), "àb "))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("À"), " a "))
		require.FalseT(t, isEqualFoldIgnoreSpace([]rune("À"), "Á"))
	})
}
