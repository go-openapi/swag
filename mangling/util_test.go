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

package mangling

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
