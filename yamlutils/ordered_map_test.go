// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 the "License";
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

package yamlutils

import (
	"encoding/json"
	"testing"

	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

func TestOrderedMap(t *testing.T) {
	t.Parallel()

	harness := fixtures.NewHarness(t) // a test suite that is common to all JSON & YAML utilities
	harness.Init()

	for name, test := range harness.AllTests(fixtures.WithoutError(true)) {
		var checkNull bool
		if name == "with null value" { // extra assertions regarding the "null" case
			checkNull = true
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			t.Run("should unmarshal JSON", func(t *testing.T) {
				var data YAMLMapSlice
				require.NoError(t, json.Unmarshal(test.JSONBytes(), &data))
				if checkNull {
					require.Nil(t, data)
					require.Empty(t, data)
				}

				t.Run("should convert JSON to YAML", func(t *testing.T) {
					y, err := data.MarshalYAML()
					require.NoError(t, err)
					if checkNull {
						require.NotEmpty(t, y) // "null" token
					}
					b, ok := y.([]byte)
					require.True(t, ok)

					assert.Equal(t, test.YAMLPayload, string(b))
				})

				t.Run("should marshal back to JSON", func(t *testing.T) {
					jazon, err := json.Marshal(data)
					require.NoError(t, err)
					if checkNull {
						require.NotEmpty(t, jazon) // "null" token
					}
					// check an exact match of JSON tokens, so this is stricter than require.JSONEq
					fixtures.JSONEqualOrdered(t, test.JSONPayload, string(jazon))
				})
			})

			t.Run("should unmarshal YAML", func(t *testing.T) {
				var data YAMLMapSlice
				require.NoError(t, yaml.Unmarshal(test.JSONBytes(), &data))

				t.Run("should convert YAML to JSON", func(t *testing.T) {
					j, err := data.MarshalJSON()
					require.NoError(t, err)

					fixtures.JSONEqualOrdered(t, test.JSONPayload, string(j))
				})

				t.Run("should marshal back to YAML", func(t *testing.T) {
					y, err := json.Marshal(data)
					require.NoError(t, err)
					// check an exact match of YAML tokens, so this is stricter than require.YAMLEq
					fixtures.YAMLEqualOrdered(t, test.YAMLPayload, string(y))
				})
			})
		})
	}

	t.Run("with complete doc", func(t *testing.T) {
		t.Run("should convert bytes to YAML doc", func(t *testing.T) {
			ydoc, err := BytesToYAMLDoc(fixture2224)
			require.NoError(t, err)

			t.Run("should convert YAML doc to JSON", func(t *testing.T) {
				jazon, err := YAMLToJSON(ydoc)
				require.NoError(t, err)

				t.Run("should unmarshal JSON into YAMLMapSlice", func(t *testing.T) {
					var data YAMLMapSlice
					require.NoError(t, json.Unmarshal(jazon, &data))

					t.Run("should marshal YAMLMapSlice into the original doc", func(t *testing.T) {
						reconstructed, err := data.MarshalYAML()
						require.NoError(t, err)

						text, ok := reconstructed.([]byte)
						require.True(t, ok)

						assert.YAMLEq(t, string(fixture2224), string(text))
					})

					t.Run("should marshal back to JSON", func(t *testing.T) {
						jazon, err := json.Marshal(data)
						require.NoError(t, err)
						// check an exact match of JSON tokens, so this is stricter than require.JSONEqual
						fixtures.JSONEqualOrdered(t, string(jazon), string(jazon))
					})
				})
			})
		})
	})
}

func TestMarshalYAML(t *testing.T) {
	t.Parallel()

	harness := fixtures.NewHarness(t) // a test suite that is common to all JSON & YAML utilities
	harness.Init()

	t.Run("marshalYAML should render nulls in values", func(t *testing.T) {
		fixture := harness.ShouldGet("with a null value")
		jazon := fixture.JSONPayload
		expected := fixture.YAMLPayload

		var data YAMLMapSlice
		require.NoError(t, json.Unmarshal([]byte(jazon), &data))
		ny, err := data.MarshalYAML()
		require.NoError(t, err)
		assert.Equal(t, expected, string(ny.([]byte)))
	})

	t.Run("marshalYAML should be deterministic", func(t *testing.T) {
		fixture := harness.ShouldGet("with numbers")
		jazon := fixture.JSONPayload
		expected := fixture.YAMLPayload

		const iterations = 10
		for range iterations {
			var data YAMLMapSlice
			require.NoError(t, json.Unmarshal([]byte(jazon), &data))
			ny, err := data.MarshalYAML()
			require.NoError(t, err)
			assert.Equal(t, expected, string(ny.([]byte)))
		}
	})

	t.Run("with only null", func(t *testing.T) {
		// the "null" token is reflected in this context as a "nil" go value, but as a non nil, empty slice.
		// The marshaling resorts to a "null" token and not to an empty string.
		fixture := harness.ShouldGet("with null value")
		input := fixture.JSONPayload
		expected := fixture.YAMLPayload // "null\n"

		t.Run("should unmarshal JSON", func(t *testing.T) {
			var data YAMLMapSlice
			require.NoError(t, json.Unmarshal([]byte(input), &data))
			require.Nil(t, data) // mutated by UnmarshalYAML

			t.Run("should convert JSON to YAML as an empty object", func(t *testing.T) {
				y, err := data.MarshalYAML()
				require.NoError(t, err)
				require.NotNil(t, y)
				b, ok := y.([]byte)
				require.True(t, ok)

				assert.Equal(t, expected, string(b))
			})

			t.Run("should marshal back to JSON", func(t *testing.T) {
				jazon, err := json.Marshal(data)
				require.NoError(t, err)
				// check an exact match of JSON tokens, so this is stricter than require.JSONEqual
				fixtures.JSONEqualOrdered(t, input, string(jazon))
			})
		})
	})

	t.Run("with maps", func(t *testing.T) {
		data := YAMLMapSlice{
			YAMLMapItem{
				Key: "a",
				Value: map[string]any{
					"x": 1,
					"y": 2,
				},
			},
		}

		t.Run("should MarshalYAML map, without ordering guarantee", func(t *testing.T) {
			const expected = `
a:
  x: 1
  y: 2
`

			y, err := data.MarshalYAML()
			require.NoError(t, err)
			require.NotNil(t, y)
			b, ok := y.([]byte)
			require.True(t, ok)

			assert.YAMLEq(t, expected, string(b))
		})
	})

	t.Run("with all numerical types", func(t *testing.T) {
		data := YAMLMapSlice{
			YAMLMapItem{
				Key: "signed",
				Value: YAMLMapSlice{
					YAMLMapItem{
						Key:   "a",
						Value: 1,
					},
					YAMLMapItem{
						Key:   "b",
						Value: int8(1),
					},
					YAMLMapItem{
						Key:   "c",
						Value: int16(1),
					},
					YAMLMapItem{
						Key:   "d",
						Value: int32(1),
					},
					YAMLMapItem{
						Key:   "e",
						Value: int64(1),
					},
				},
			},
			YAMLMapItem{
				Key: "unsigned",
				Value: YAMLMapSlice{
					YAMLMapItem{
						Key:   "a",
						Value: uint(1),
					},
					YAMLMapItem{
						Key:   "b",
						Value: uint8(1),
					},
					YAMLMapItem{
						Key:   "c",
						Value: uint16(1),
					},
					YAMLMapItem{
						Key:   "d",
						Value: uint32(1),
					},
					YAMLMapItem{
						Key:   "e",
						Value: uint64(1),
					},
				},
			},
			YAMLMapItem{
				Key: "float",
				Value: YAMLMapSlice{
					YAMLMapItem{
						Key:   "a",
						Value: float32(1.6),
					},
					YAMLMapItem{
						Key:   "b",
						Value: 1.6,
					},
				},
			},
		}

		t.Run("should MarshalYAML map, without ordering guarantee", func(t *testing.T) {
			const expected = `
signed:
  a: 1
  b: 1
  c: 1
  d: 1
  e: 1
unsigned:
  a: 1
  b: 1
  c: 1
  d: 1
  e: 1
float:
  a: 1.6
  b: 1.6
`
			y, err := data.MarshalYAML()
			require.NoError(t, err)
			require.NotNil(t, y)
			b, ok := y.([]byte)
			require.True(t, ok)

			fixtures.YAMLEqualOrdered(t, expected, string(b))
		})
	})
}

func TestUnmarshalYAML(t *testing.T) {
	t.Parallel()

	data := YAMLMapSlice{}
	t.Run("UnmarshalYAML of a nil node should just pass without error", func(t *testing.T) {
		require.NoError(t, data.UnmarshalYAML(nil))
	})
}

func TestSetOrdered(t *testing.T) {
	t.Parallel()
	data := YAMLMapSlice{} // can't be nil

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
		require.Equal(t, YAMLMapItem{Key: "a", Value: 1}, data[0])
		require.Equal(t, YAMLMapItem{Key: "b", Value: true}, data[1])
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
		require.Equal(t, YAMLMapItem{Key: "a", Value: 2}, data[0])    // merged
		require.Equal(t, YAMLMapItem{Key: "b", Value: true}, data[1]) // unchanged
		require.Equal(t, YAMLMapItem{Key: "c", Value: "x"}, data[2])  // appended
	})

	t.Run("with nil items should yield nil", func(t *testing.T) {
		data.SetOrderedItems(nil)
		require.Nil(t, data)

	})
}
