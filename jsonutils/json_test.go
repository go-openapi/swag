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
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type SharedCounters struct {
	Counter1 int64 `json:"counter1,omitempty"`
	Counter2 int64 `json:"counter2:,omitempty"` // the ":" in the json field name is left on-purpose for this test
}

type AggregationObject struct {
	SharedCounters

	Count int64 `json:"count,omitempty"`
}

func (m *AggregationObject) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 SharedCounters
	if err := ReadJSON(raw, &aO0); err != nil {
		return err
	}

	m.SharedCounters = aO0

	// now for regular properties
	var propsAggregationObject struct {
		Count int64 `json:"count,omitempty"`
	}
	if err := ReadJSON(raw, &propsAggregationObject); err != nil {
		return err
	}

	m.Count = propsAggregationObject.Count

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m AggregationObject) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 1)

	aO0, err := WriteJSON(m.SharedCounters)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	// now for regular properties
	var propsAggregationObject struct {
		Count int64 `json:"count,omitempty"`
	}
	propsAggregationObject.Count = m.Count

	jsonDataPropsAggregationObject, errAggregationObject := WriteJSON(propsAggregationObject)
	if errAggregationObject != nil {
		return nil, errAggregationObject
	}
	_parts = append(_parts, jsonDataPropsAggregationObject)

	return ConcatJSON(_parts...), nil
}

func TestReadWriteJSON(t *testing.T) {
	obj := AggregationObject{Count: 290, SharedCounters: SharedCounters{Counter1: 304, Counter2: 948}}

	t.Run("with default adapter", func(t *testing.T) {
		t.Run("should WriteJSON from struct", func(t *testing.T) {
			rtjson, err := WriteJSON(obj)
			require.NoError(t, err)

			t.Run("should MarshalJSON using WriteJSON from this type", func(t *testing.T) {
				otjson, err := obj.MarshalJSON()
				require.NoError(t, err)

				t.Run("both marshaling methods should be equivalent", func(t *testing.T) {
					require.JSONEq(t, string(rtjson), string(otjson))
				})
			})

			t.Run("should MarshalJSON using the standard library", func(t *testing.T) {
				otjson, err := json.Marshal(obj)
				require.NoError(t, err)

				t.Run("both marshaling methods should be equivalent", func(t *testing.T) {
					require.JSONEq(t, string(rtjson), string(otjson))
				})
			})

			t.Run("should ReadJSON into new struct", func(t *testing.T) {
				var obj1 AggregationObject
				require.NoError(t, ReadJSON(rtjson, &obj1))

				t.Run("this should copy the object", func(t *testing.T) {
					require.Equal(t, obj, obj1)
				})
			})

			t.Run("should UnmarshalJSON using ReadJSON into new struct", func(t *testing.T) {
				var obj11 AggregationObject
				require.NoError(t, obj11.UnmarshalJSON(rtjson))

				t.Run("this should copy the object", func(t *testing.T) {
					require.Equal(t, obj, obj11)
				})
			})

			t.Run("should UnmarshalJSON using the standard library", func(t *testing.T) {
				var obj11 AggregationObject
				require.NoError(t, json.Unmarshal(rtjson, &obj11))

				t.Run("this should copy the object", func(t *testing.T) {
					require.Equal(t, obj, obj11)
				})
			})
		})

		t.Run("with counters", func(t *testing.T) {
			t.Run("should ReadJSON into struct", func(t *testing.T) {
				jsons := `{"counter1":123,"counter2:":456,"count":999}`
				var obj2 AggregationObject

				require.NoError(t, ReadJSON([]byte(jsons), &obj2))
				require.Equal(t, AggregationObject{SharedCounters: SharedCounters{Counter1: 123, Counter2: 456}, Count: 999}, obj2)
			})
		})
		t.Run("using FromDynamicJSON", func(t *testing.T) {
			const epsilon = 1e-6
			var obj2 any

			require.NoError(t, FromDynamicJSON(obj, &obj2))
			asMap, ok := obj2.(map[string]any)
			require.True(t, ok)
			assert.Len(t, asMap, 3) // 3 fields in struct
			c1, ok := asMap["counter1"]
			require.True(t, ok)
			assert.InDelta(t, float64(304), c1, epsilon)

			c2, ok := asMap["counter2:"]
			require.True(t, ok)
			assert.InDelta(t, float64(948), c2, epsilon)

			c, ok := asMap["count"]
			require.True(t, ok)
			assert.InDelta(t, float64(290), c, epsilon)
		})

		t.Run("with error edge cases", func(t *testing.T) {
			t.Run("should not unmarshal non pointer, nil interface", func(t *testing.T) {
				var obj2 any

				err := FromDynamicJSON(obj, obj2)
				require.Error(t, err)
				require.ErrorContains(t, err, "Unmarshal(nil)")
			})

			t.Run("should not unmarshal non pointer struct", func(t *testing.T) {
				var obj2 struct{}

				err := FromDynamicJSON(obj, obj2)
				require.Error(t, err)
				require.ErrorContains(t, err, "Unmarshal(non-pointer struct {})")
			})

			t.Run("should not unmarshal non-serializable exported field", func(t *testing.T) {
				var obj2 any
				var source struct {
					A int `json:"a"`
					B func()
				}
				err := FromDynamicJSON(source, obj2)
				require.Error(t, err)
				require.ErrorContains(t, err, "unsupported type: func()")
			})
		})
	})
}
