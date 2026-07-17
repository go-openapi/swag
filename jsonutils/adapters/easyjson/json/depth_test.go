// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// deepObject builds a JSON document nesting depth objects: {"a":{"a":...{}...}}.
func deepObject(depth int) []byte {
	return []byte(strings.Repeat(`{"a":`, depth) + `{}` + strings.Repeat(`}`, depth))
}

// deepArray builds a JSON document nesting depth arrays under one key: {"a":[[...[]...]]}.
func deepArray(depth int) []byte {
	return []byte(`{"a":` + strings.Repeat(`[`, depth) + strings.Repeat(`]`, depth) + `}`)
}

func TestMaxNestingDepthUnmarshal(t *testing.T) {
	t.Run("object nesting within the limit should unmarshal", func(t *testing.T) {
		var m MapSlice
		require.NoError(t, m.UnmarshalJSON(deepObject(100)))
	})

	t.Run("object nesting beyond the default limit should error, not crash", func(t *testing.T) {
		var m MapSlice
		err := m.UnmarshalJSON(deepObject(defaultMaxNestingDepth + 5))
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMaxNestingDepth)
	})

	t.Run("array nesting beyond the default limit should error, not crash", func(t *testing.T) {
		var m MapSlice
		err := m.UnmarshalJSON(deepArray(defaultMaxNestingDepth + 5))
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMaxNestingDepth)
	})

	t.Run("configurable limit via adapter option should be honored", func(t *testing.T) {
		a := NewAdapter(WithMaxNestingDepth(5))

		var okMap MapSlice
		require.NoError(t, a.OrderedUnmarshal(deepObject(4), &okMap))

		var badMap MapSlice
		err := a.OrderedUnmarshal(deepObject(10), &badMap)
		require.Error(t, err)
	})
}

func TestMaxNestingDepthMarshal(t *testing.T) {
	// deepMapSlice builds an in-memory MapSlice nested depth levels deep.
	deepMapSlice := func(depth int) MapSlice {
		m := MapSlice{{Key: "leaf", Value: "x"}}
		for range depth {
			m = MapSlice{{Key: "n", Value: m}}
		}
		return m
	}

	t.Run("adapter OrderedMarshal beyond the default limit should error, not crash", func(t *testing.T) {
		a := NewAdapter()
		_, err := a.OrderedMarshal(deepMapSlice(defaultMaxNestingDepth + 5))
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMaxNestingDepth)
	})

	t.Run("adapter OrderedMarshal within the limit should marshal", func(t *testing.T) {
		a := NewAdapter()
		_, err := a.OrderedMarshal(deepMapSlice(100))
		require.NoError(t, err)
	})

	t.Run("MarshalJSON path beyond the default limit should error, not crash", func(t *testing.T) {
		_, err := deepMapSlice(defaultMaxNestingDepth + 5).MarshalJSON()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMaxNestingDepth)
	})

	t.Run("configurable limit via adapter option should be honored", func(t *testing.T) {
		a := NewAdapter(WithMaxNestingDepth(5))
		_, err := a.OrderedMarshal(deepMapSlice(10))
		require.Error(t, err)
	})
}
