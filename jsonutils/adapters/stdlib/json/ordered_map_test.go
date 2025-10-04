// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	stdjson "encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

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

		fixtures.JSONEqualOrderedBytes(t, []byte(`{"a":1,"c":"y","b":2}`), jazon)
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

				fixtures.JSONEqualOrderedBytes(t, test.JSONBytes(), jazon)
			})

			t.Run("should keep the order of keys", func(t *testing.T) {
				fixture := harness.ShouldGet("with numbers")
				input := fixture.JSONBytes()

				const iterations = 10
				for range iterations {
					var data MapSlice
					require.NoError(t, stdjson.Unmarshal(input, &data))
					jazon, err := stdjson.Marshal(data)
					require.NoError(t, err)

					fixtures.JSONEqualOrderedBytes(t, input, jazon) // specifically check the same order, not require.JSONEq()
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

func TestReproDataRace(t *testing.T) {
	t.Parallel()
	const parallelRoutines = 1000

	// NOTE: with go1.25, use synctest.Test
	var wg sync.WaitGroup

	for range parallelRoutines {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			toks := make([]token, 0, 4)
			buf := []byte(`{"test":"data"}`)
			l := poolOfLexers.Borrow(buf)

			for tok := l.NextToken(); tok != eofToken; tok = l.NextToken() {
				toks = append(toks, tok)
			}
			assert.Len(t, toks, 4)
			fmt.Fprintf(io.Discard, "%d", len(toks))
			defer func() {
				poolOfLexers.Redeem(l)
			}()
		}()
	}

	wg.Wait()
}
