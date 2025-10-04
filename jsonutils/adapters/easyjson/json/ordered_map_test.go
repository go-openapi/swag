// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

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
