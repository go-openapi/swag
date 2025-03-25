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

package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		t.Run("ReadJSON and WriteJSON back", func(t *testing.T) {
			var b any

			const jazon = `{"a": 1}`
			require.NoError(t, ReadJSON([]byte(jazon), &b))

			buf, err := WriteJSON(b)
			require.NoError(t, err)
			require.JSONEq(t, jazon, string(buf))
		})

		t.Run("ConcatJSON merge 2 objects", func(t *testing.T) {
			buf := ConcatJSON([]byte(`{"a": 1}`), []byte(`{"b": 2}`))
			require.JSONEq(t, `{"a": 1, "b": 2}`, string(buf))
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
