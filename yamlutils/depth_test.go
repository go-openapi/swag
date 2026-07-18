// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package yamlutils

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	yaml "go.yaml.in/yaml/v3"
)

func TestMaxNestingDepth(t *testing.T) {
	t.Run("YAMLToJSON on deeply nested slices should error, not crash", func(t *testing.T) {
		var v any = "leaf"
		for range defaultMaxNestingDepth + 5 {
			v = []any{v}
		}

		_, err := YAMLToJSON(v)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrYAML)
	})

	t.Run("YAMLToJSON on deeply nested maps should error, not crash", func(t *testing.T) {
		var v any = "leaf"
		for range defaultMaxNestingDepth + 5 {
			v = map[any]any{"n": v}
		}

		_, err := YAMLToJSON(v)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrYAML)
	})

	t.Run("MarshalYAML on a deeply nested YAMLMapSlice should error, not crash", func(t *testing.T) {
		var v any = "leaf"
		for range defaultMaxNestingDepth + 5 {
			v = YAMLMapSlice{{Key: "n", Value: v}}
		}
		top, ok := v.(YAMLMapSlice)
		require.True(t, ok)

		_, err := yaml.Marshal(top)
		require.Error(t, err)
	})

	t.Run("BytesToYAMLDoc on deeply nested YAML should error, not crash", func(t *testing.T) {
		const depth = defaultMaxNestingDepth + 5
		data := []byte("a: " + strings.Repeat("[", depth) + strings.Repeat("]", depth))

		_, err := BytesToYAMLDoc(data)
		require.Error(t, err)
	})

	t.Run("moderately nested document should round-trip cleanly", func(t *testing.T) {
		const depth = 200
		data := []byte("a: " + strings.Repeat("[", depth) + strings.Repeat("]", depth))

		doc, err := BytesToYAMLDoc(data)
		require.NoError(t, err)

		_, err = YAMLToJSON(doc)
		require.NoError(t, err)
	})
}
