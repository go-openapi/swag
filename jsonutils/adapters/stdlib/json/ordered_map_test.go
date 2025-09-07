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
	stdjson "encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
)

func TestSetOrdered(t *testing.T) {
	t.Parallel()

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

func TestMapSlice(t *testing.T) {
	t.Parallel()

	harness := fixtures.NewHarness(t)
	harness.Init()

	for name, test := range harness.AllTests() {
		// in this testcase, "null" renders a nil as expected.
		// Notice the difference in how we declared the target:
		//
		// 1.  var data MapSlice => will be set to nil
		// 2.  data := make(MapSlice,0,10) => will be set to empty
		t.Run(name, func(t *testing.T) {
			t.Run("should unmarshal and marshal MapSlice", func(t *testing.T) {
				var data MapSlice
				if test.ExpectError() {
					require.Error(t, stdjson.Unmarshal(test.JSONBytes(), &data))
					return
				}

				require.NoError(t, stdjson.Unmarshal(test.JSONBytes(), &data))

				jazon, err := stdjson.Marshal(data)
				require.NoError(t, err)

				fixtures.JSONEqualOrdered(t, test.JSONPayload, string(jazon))
			})

			t.Run("should keep the order of keys", func(t *testing.T) {
				fixture := harness.ShouldGet("with numbers")
				input := fixture.JSONPayload

				const iterations = 10
				for range iterations {
					var data MapSlice
					require.NoError(t, stdjson.Unmarshal([]byte(input), &data))
					jazon, err := stdjson.Marshal(data)
					require.NoError(t, err)

					fixtures.JSONEqualOrdered(t, input, string(jazon)) // specifically check the same order, not require.JSONEq()
				}
			})
		})
	}
}

func TestLexerErrors(t *testing.T) {
	t.Parallel()

	harness := fixtures.NewHarness(t)
	harness.Init()

	for name, test := range harness.AllTests(fixtures.WithError(true)) {
		t.Run(name, func(t *testing.T) {
			t.Run("should raise a lexer error", func(t *testing.T) {
				// test directly this endpoint, as the json standard library
				// performs a preventive early check for well-formed JSON.
				data := make(MapSlice, 0)
				l := newLexer(test.JSONBytes())
				data.unmarshalObject(l)
				err := l.Error()
				require.ErrorIs(t, err, ErrStdlib)
			})
		})
	}
}
