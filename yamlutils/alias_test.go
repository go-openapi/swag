// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package yamlutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// aliasBomb builds a "billion laughs" YAML document: each anchor references the previous
// one `fan` times, so the expanded size grows like fan^levels while the text stays tiny.
func aliasBomb(levels, fan int) []byte {
	var b strings.Builder
	b.WriteString("a0: &a0 \"lol\"\n")
	for i := 1; i <= levels; i++ {
		fmt.Fprintf(&b, "a%d: &a%d [", i, i)
		for j := range fan {
			if j > 0 {
				b.WriteByte(',')
			}
			fmt.Fprintf(&b, "*a%d", i-1)
		}
		b.WriteString("]\n")
	}
	return []byte(b.String())
}

func TestAliasExpansion(t *testing.T) {
	t.Run("an alias bomb should be rejected, not expanded", func(t *testing.T) {
		doc, err := BytesToYAMLDoc(aliasBomb(20, 2))
		require.NoError(t, err) // the parse step keeps aliases unexpanded

		_, err = YAMLToJSON(doc)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrYAML)
		assert.Contains(t, err.Error(), "excessive aliasing")
	})

	t.Run("a self-referential anchor should error cleanly as a cycle", func(t *testing.T) {
		doc, err := BytesToYAMLDoc([]byte("a: &a\n  b: *a\n"))
		require.NoError(t, err)

		_, err = YAMLToJSON(doc)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrYAML)
		assert.Contains(t, err.Error(), "contains itself")
		// the cycle is caught explicitly, not by exhausting the depth guard
		assert.NotContains(t, err.Error(), "maximum nesting depth")
	})

	t.Run("legitimate, bounded alias use should still resolve", func(t *testing.T) {
		doc, err := BytesToYAMLDoc([]byte("defs: &d [1, 2, 3]\na: *d\nb: *d\n"))
		require.NoError(t, err)

		out, err := YAMLToJSON(doc)
		require.NoError(t, err)
		assert.JSONEq(t, `{"defs":[1,2,3],"a":[1,2,3],"b":[1,2,3]}`, string(out))
	})
}
