// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestJSONUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		t.Run("ReadJSON and WriteJSON back", func(t *testing.T) {
			var b any

			jazon := []byte(`{"a": 1}`)
			require.NoError(t, ReadJSON(jazon, &b))

			buf, err := WriteJSON(b)
			require.NoError(t, err)
			require.JSONEqBytes(t, jazon, buf)
		})

		t.Run("ConcatJSON merge 2 objects", func(t *testing.T) {
			buf := ConcatJSON([]byte(`{"a": 1}`), []byte(`{"b": 2}`))
			require.JSONEqBytes(t, []byte(`{"a": 1, "b": 2}`), buf)
		})

		t.Run("with go struct", func(t *testing.T) {
			var a struct {
				A int
			}
			a.A = 1

			t.Run("FromDynamicJSON into a map", func(t *testing.T) {
				var b any

				require.NoError(t, FromDynamicJSON(a, &b))

				_, isMap := b.(map[string]any)
				assert.True(t, isMap)
			})

			t.Run("ToDynamicJSON into a map", func(t *testing.T) {
				c := ToDynamicJSON(a)
				_, isMap := c.(map[string]any)
				assert.True(t, isMap)

				t.Run("DynamicJSONToStruct back to struct", func(t *testing.T) {
					a.A = 0
					require.NoError(t, DynamicJSONToStruct(c, &a))
					assert.Equalf(t, 1, a.A, "expected to restore original value")
				})
			})
		})
	})
}
