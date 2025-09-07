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

package jsonutils

import (
	"encoding/json"
	"iter"
	"testing"

	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONMapSlice(t *testing.T) {
	harness := fixtures.NewHarness(t)
	harness.Init()

	for name, test := range harness.AllTests(fixtures.WithoutError(true)) {
		t.Run(name, func(t *testing.T) {
			t.Run("should unmarshal JSON", func(t *testing.T) {
				input := test.JSONPayload

				var data JSONMapSlice
				require.NoError(t, json.Unmarshal([]byte(input), &data))

				t.Run("should marshal JSON", func(t *testing.T) {
					jazon, err := json.Marshal(data)
					require.NoError(t, err)

					fixtures.JSONEqualOrdered(t, input, string(jazon))
				})
			})
		})
	}

	t.Run("should keep the order of keys", func(t *testing.T) {
		fixture := harness.ShouldGet("with numbers")
		input := fixture.JSONPayload

		const iterations = 10
		for range iterations {
			var data JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(input), &data))
			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			fixtures.JSONEqualOrdered(t, input, string(jazon)) // specifically check the same order, not require.JSONEq()
		}
	})

	t.Run("key ordering doesn't have to be stable with Read/Write JSON", func(t *testing.T) {
		fixture := harness.ShouldGet("with numbers")
		input := fixture.JSONPayload

		var data JSONMapSlice
		require.NoError(t, json.Unmarshal([]byte(input), &data))

		var obj any
		require.NoError(t, FromDynamicJSON(data, &obj))

		asMap, ok := obj.(map[string]any)
		require.True(t, ok)
		assert.Len(t, asMap, 3) // 3 fields in struct

		var target JSONMapSlice
		require.NoError(t, FromDynamicJSON(obj, &target))

		// the order of keys may be altered, since the intermediary representation is a map[string]any
		jazon, err := WriteJSON(target)
		require.NoError(t, err)
		require.JSONEq(t, input, string(jazon))
	})

	t.Run("key ordering is maintained with nested ifaces.Ordered types", func(t *testing.T) {
		fixture := harness.ShouldGet("with numbers")
		input := fixture.JSONPayload

		var data JSONMapSlice
		require.NoError(t, json.Unmarshal([]byte(input), &data))
		require.Len(t, data, 3) // 3 fields

		custom := makeCustomOrdered(
			data...,
		)

		const iterations = 10
		for range iterations {
			var obj customOrdered
			require.NoError(t, FromDynamicJSON(custom, &obj))
			assert.Len(t, obj.elems, 3) // 3 fields in struct

			var target JSONMapSlice
			require.NoError(t, FromDynamicJSON(obj, &target))
			// the order of keys may is maintained
			jazon, err := WriteJSON(target)
			require.NoError(t, err)
			fixtures.JSONEqualOrdered(t, input, string(jazon))
		}
	})

	t.Run("UnmarshalJSON with error cases", func(t *testing.T) {
		// test directly this endpoint, as the json standard library
		// performs a preventive early check for well-formed JSON.
		for name, test := range harness.AllTests(fixtures.WithError(true)) {
			t.Run(name, func(t *testing.T) {
				t.Run("should yield an error", func(t *testing.T) {
					var data JSONMapSlice
					require.Error(t, json.Unmarshal(test.JSONBytes(), &data))
				})
			})
		}
	})
}

func TestSetOrdered(t *testing.T) {
	t.Parallel()
	data := JSONMapSlice{} // can't be nil

	t.Run("should insert keys", func(t *testing.T) {
		kv := []struct {
			k string
			v any
		}{
			{k: "a", v: 1},
			{k: "b", v: true},
		}

		data.SetOrderedItems(func(yield func(string, any) bool) {
			for _, e := range kv {
				if !yield(e.k, e.v) {
					return
				}
			}
		})

		require.Len(t, data, len(kv))
		require.Equal(t, JSONMapItem{Key: "a", Value: 1}, data[0])
		require.Equal(t, JSONMapItem{Key: "b", Value: true}, data[1])
	})

	t.Run("should merge keys", func(t *testing.T) {
		kv := []struct {
			k string
			v any
		}{
			{k: "a", v: 2},
			{k: "c", v: "x"},
		}

		data.SetOrderedItems(func(yield func(string, any) bool) {
			for _, e := range kv {
				if !yield(e.k, e.v) {
					return
				}
			}
		})

		require.Len(t, data, len(kv)+1)
		require.Equal(t, JSONMapItem{Key: "a", Value: 2}, data[0])    // merged
		require.Equal(t, JSONMapItem{Key: "b", Value: true}, data[1]) // unchanged
		require.Equal(t, JSONMapItem{Key: "c", Value: "x"}, data[2])  // appended
	})

	t.Run("with nil items should yield nil", func(t *testing.T) {
		data.SetOrderedItems(nil)
		require.Nil(t, data)

	})
}

// customOrdered implements ifaces.Ordered and ifaces.SetOrdered: ReadJSON and WriteJSON
// should recognize this and honor the ordering of keys.
//
// Technically, this illustrates an alternate implementation to JSONMapSlice, with which
// retrieving keys is a constant-time operation.
type customOrdered struct {
	elems []JSONMapItem
	idx   map[string]int
}

var (
	_ ifaces.Ordered    = customOrdered{}
	_ ifaces.SetOrdered = &customOrdered{}
)

func makeCustomOrdered(items ...JSONMapItem) customOrdered {
	o := customOrdered{
		elems: make([]JSONMapItem, len(items)),
		idx:   make(map[string]int, len(items)),
	}

	for i, item := range items {
		o.elems[i] = item
		o.idx[item.Key] = i
	}

	return o
}

func (o customOrdered) Get(key string) (any, bool) {
	idx, ok := o.idx[key]
	if !ok {
		return nil, false
	}

	return o.elems[idx].Value, true
}

func (o *customOrdered) Set(key string, value any) bool {
	idx, ok := o.idx[key]
	if !ok {
		o.elems = append(o.elems, JSONMapItem{Key: key, Value: value})
		o.idx[key] = len(o.elems)

		return false
	}

	o.elems[idx].Value = value

	return true
}

func (o customOrdered) OrderedItems() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for _, item := range o.elems {
			if !yield(item.Key, item.Value) {
				return
			}
		}
	}
}

func (o *customOrdered) SetOrderedItems(items iter.Seq2[string, any]) {
	if items == nil {
		o.elems = nil
		clear(o.idx)

		return
	}
	if o.idx == nil {
		o.idx = make(map[string]int, 0)
	}

	for k, v := range items {
		_ = o.Set(k, v)
	}
}
