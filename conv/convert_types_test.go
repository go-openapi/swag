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

package conv

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag/typeutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	wantsPointer = true
	wantsValue   = false
)

func TestSlice(t *testing.T) {
	t.Run("with nulls, should skip null map entries", func(t *testing.T) {
		require.Empty(t, ValueSlice[uint16](nil))

		require.Len(t, ValueSlice([]*uint16{Pointer(uint16(1)), nil, Pointer(uint16(2))}), 3)
	})

	t.Run("with PointerSlice on []string", func(t *testing.T) {
		testCasesStringSlice := [][]string{
			{"a", "b", "c", "d", "e"},
			{"a", "b", "", "", "e"},
		}

		for idx, in := range testCasesStringSlice {
			if in == nil {
				continue
			}
			out := PointerSlice(in)
			assertValues(t, in, out, wantsPointer, idx)

			out2 := ValueSlice(out)
			assertValues(t, in, out2, wantsValue, idx)
		}
	})

	t.Run("with ValueSlice on []string", func(t *testing.T) {
		testCasesStringValueSlice := [][]*string{
			{Pointer("a"), Pointer("b"), nil, Pointer("c")},
		}

		for idx, in := range testCasesStringValueSlice {
			if in == nil {
				continue
			}
			out := ValueSlice(in)
			assertValues(t, in, out, wantsValue, idx)

			out2 := PointerSlice(out)
			assertValues(t, in, out2, wantsPointer, idx)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("with nulls", func(t *testing.T) {
		require.Empty(t, ValueMap[string, uint16](nil))

		require.Len(t, ValueMap(map[string]*int{"a": Pointer(1), "b": nil, "c": Pointer(2)}), 2)
	})

	t.Run("with PointerMap on map[string]string", func(t *testing.T) {
		testCasesStringMap := []map[string]string{
			{"a": "1", "b": "2", "c": "3"},
		}

		for idx, in := range testCasesStringMap {
			if in == nil {
				continue
			}
			out := PointerMap(in)
			assertValues(t, in, out, wantsPointer, idx)

			out2 := ValueMap(out)
			assertValues(t, in, out2, wantsValue, idx)
		}
	})

	t.Run("with ValueMap on map[string]bool", func(t *testing.T) {
		testCasesBoolMap := []map[string]bool{
			{"a": true, "b": false, "c": true},
		}

		for idx, in := range testCasesBoolMap {
			if in == nil {
				continue
			}
			out := PointerMap(in)
			assertValues(t, in, out, wantsPointer, idx)

			out2 := ValueMap(out)
			assertValues(t, in, out2, wantsValue, idx)
		}
	})
}

func TestPointer(t *testing.T) {
	t.Run("with Pointer on string", func(t *testing.T) {
		testCasesString := []string{"a", "b", "c", "d", "e", ""}

		for idx, in := range testCasesString {
			out := Pointer(in)
			assertValues(t, in, out, wantsPointer, idx)

			out2 := Value(out)
			assertValues(t, in, out2, wantsValue, idx)
		}
		assert.Zerof(t, Value[string](nil), "expected conversion from nil to return zero value")
	})

	t.Run("with Value on bool", func(t *testing.T) {
		testCasesBool := []bool{true, false}

		for idx, in := range testCasesBool {
			out := Pointer(in)
			assertValues(t, in, out, wantsPointer, idx)

			out2 := Value(out)
			assertValues(t, in, out2, wantsValue, idx)
		}
		assert.Zerof(t, Value[bool](nil), "expected conversion from nil to return zero value")
	})
}

func assertSingleValue(t *testing.T, inElem, elem reflect.Value, expectPointer bool, idx int) {
	require.Equalf(t,
		expectPointer, (elem.Kind() == reflect.Ptr),
		"unexpected expectPointer=%t value type %T at idx %d", expectPointer, elem, idx,
	)

	if inElem.Kind() == reflect.Ptr && !inElem.IsNil() {
		inElem = reflect.Indirect(inElem)
	}

	if elem.Kind() == reflect.Ptr && !elem.IsNil() {
		elem = reflect.Indirect(elem)
	}

	require.Truef(t,
		(elem.Kind() == reflect.Ptr && elem.IsNil()) ||
			typeutils.IsZero(elem.Interface()) == (inElem.Kind() == reflect.Ptr && inElem.IsNil()) ||
			typeutils.IsZero(inElem.Interface()),
		"unexpected nil pointer at idx %d", idx,
	)

	if !((elem.Kind() == reflect.Ptr && elem.IsNil()) || typeutils.IsZero(elem.Interface())) {
		require.IsTypef(t, inElem.Interface(), elem.Interface(),
			"expected in/out to match types at idx %d", idx,
		)
		assert.EqualValuesf(t, inElem.Interface(), elem.Interface(),
			"unexpected value at idx %d: %v", idx, elem.Interface(),
		)
	}
}

// assertValues checks equivalent representation pointer vs values for single var, slices and maps
func assertValues(t *testing.T, in, out interface{}, expectPointer bool, idx int) {
	vin := reflect.ValueOf(in)
	vout := reflect.ValueOf(out)

	switch vin.Kind() { //nolint:exhaustive
	case reflect.Slice, reflect.Map:
		require.Equalf(t, vin.Kind(), vout.Kind(),
			"unexpected output type at idx %d", idx,
		)
		require.Equalf(t, vin.Len(), vout.Len(),
			"unexpected len at idx %d", idx,
		)

		var elem, inElem reflect.Value
		for i := 0; i < vin.Len(); i++ {
			switch vin.Kind() { //nolint:exhaustive
			case reflect.Slice:
				elem = vout.Index(i)
				inElem = vin.Index(i)
			case reflect.Map:
				keys := vin.MapKeys()
				elem = vout.MapIndex(keys[i])
				inElem = vout.MapIndex(keys[i])
			default:
			}

			assertSingleValue(t, inElem, elem, expectPointer, idx)
		}

	default:
		inElem := vin
		elem := vout

		assertSingleValue(t, inElem, elem, expectPointer, idx)
	}
}
