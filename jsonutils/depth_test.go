// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonutils_test

import (
	"strings"
	"testing"

	"github.com/go-openapi/swag/jsonutils"
	"github.com/go-openapi/testify/v2/require"
)

// TestReadJSONDeepNestingDoesNotCrash exercises the advisory scenario end-to-end:
// a deeply nested document must return an error instead of driving the runtime to a
// non-recoverable stack overflow (CWE-674, unchecked recursion).
func TestReadJSONDeepNestingDoesNotCrash(t *testing.T) {
	t.Run("deeply nested arrays should error, not crash", func(t *testing.T) {
		const depth = 20000 // well beyond the default 10000 limit
		payload := []byte(`{"a":` + strings.Repeat("[", depth) + strings.Repeat("]", depth) + `}`)

		var v jsonutils.JSONMapSlice
		require.Error(t, jsonutils.ReadJSON(payload, &v))
	})

	t.Run("deeply nested objects should error, not crash", func(t *testing.T) {
		const depth = 20000
		payload := []byte(strings.Repeat(`{"a":`, depth) + `{}` + strings.Repeat(`}`, depth))

		var v jsonutils.JSONMapSlice
		require.Error(t, jsonutils.ReadJSON(payload, &v))
	})

	t.Run("moderately nested document should round-trip cleanly", func(t *testing.T) {
		const depth = 200
		payload := []byte(strings.Repeat(`{"a":`, depth) + `{}` + strings.Repeat(`}`, depth))

		var v jsonutils.JSONMapSlice
		require.NoError(t, jsonutils.ReadJSON(payload, &v))
	})
}

// TestWriteJSONDeepNestingDoesNotCrash covers the marshal path: a deep in-memory
// JSONMapSlice must error rather than overflow the stack.
func TestWriteJSONDeepNestingDoesNotCrash(t *testing.T) {
	t.Run("deeply nested value should error, not crash", func(t *testing.T) {
		const depth = 20000
		v := jsonutils.JSONMapSlice{{Key: "leaf", Value: "x"}}
		for range depth {
			v = jsonutils.JSONMapSlice{{Key: "n", Value: v}}
		}

		_, err := jsonutils.WriteJSON(v)
		require.Error(t, err)
	})
}
