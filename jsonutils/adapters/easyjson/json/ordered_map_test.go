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

package json

import (
	"testing"

	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
	"github.com/go-openapi/testify/v2/require"
)

func TestSetOrdered(t *testing.T) {
	t.Run("should merge keys", func(t *testing.T) {
		m := MapSlice{}
		const initial = `{"a":"x","c":"y"}`
		require.NoError(t, m.UnmarshalJSON([]byte(initial)))

		appender := func(yield func(string, any) bool) {
			elements := MapSlice{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
			}

			for _, elem := range elements {
				if !yield(elem.Key, elem.Value) {
					return
				}
			}
		}

		m.SetOrderedItems(appender)

		jazon, err := m.MarshalJSON()
		require.NoError(t, err)

		fixtures.JSONEqualOrdered(t, `{"a":1,"c":"y","b":2}`, string(jazon))
	})

	t.Run("should reset keys", func(t *testing.T) {
		m := MapSlice{}
		const initial = `{"a":"x","c":"y"}`
		require.NoError(t, m.UnmarshalJSON([]byte(initial)))
		m.SetOrderedItems(nil)
		require.Nil(t, m)
	})
}
