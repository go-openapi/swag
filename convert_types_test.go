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
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringSlice(t *testing.T) {
	testCasesStringSlice := [][]string{
		{"a", "b", "c", "d", "e"},
		{"a", "b", "", "", "e"},
	}

	for idx, in := range testCasesStringSlice {
		if in == nil {
			continue
		}
		out := StringSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := StringValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestStringValueSlice(t *testing.T) {
	testCasesStringValueSlice := [][]*string{
		{String("a"), String("b"), nil, String("c")},
	}

	for idx, in := range testCasesStringValueSlice {
		if in == nil {
			continue
		}
		out := StringValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := StringSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestStringMap(t *testing.T) {
	testCasesStringMap := []map[string]string{
		{"a": "1", "b": "2", "c": "3"},
	}

	for idx, in := range testCasesStringMap {
		if in == nil {
			continue
		}
		out := StringMap(in)
		assertValues(t, in, out, true, idx)

		out2 := StringValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestBoolSlice(t *testing.T) {
	testCasesBoolSlice := [][]bool{
		{true, true, false, false},
	}

	for idx, in := range testCasesBoolSlice {
		if in == nil {
			continue
		}
		out := BoolSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := BoolValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestBoolValueSlice(t *testing.T) {
	testCasesBoolValueSlice := [][]*bool{
		{Bool(true), Bool(true), Bool(false), Bool(false)},
	}

	for idx, in := range testCasesBoolValueSlice {
		if in == nil {
			continue
		}
		out := BoolValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := BoolSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestBoolMap(t *testing.T) {
	testCasesBoolMap := []map[string]bool{
		{"a": true, "b": false, "c": true},
	}

	for idx, in := range testCasesBoolMap {
		if in == nil {
			continue
		}
		out := BoolMap(in)
		assertValues(t, in, out, true, idx)

		out2 := BoolValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestIntSlice(t *testing.T) {
	testCasesIntSlice := [][]int{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesIntSlice {
		if in == nil {
			continue
		}
		out := IntSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := IntValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestIntValueSlice(t *testing.T) {
	testCasesIntValueSlice := [][]*int{
		{Int(1), Int(2), Int(3), Int(4)},
	}

	for idx, in := range testCasesIntValueSlice {
		if in == nil {
			continue
		}
		out := IntValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := IntSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestIntMap(t *testing.T) {
	testCasesIntMap := []map[string]int{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesIntMap {
		if in == nil {
			continue
		}
		out := IntMap(in)
		assertValues(t, in, out, true, idx)

		out2 := IntValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestInt64Slice(t *testing.T) {
	testCasesInt64Slice := [][]int64{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesInt64Slice {
		if in == nil {
			continue
		}
		out := Int64Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Int64ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestInt64ValueSlice(t *testing.T) {
	testCasesInt64ValueSlice := [][]*int64{
		{Int64(1), Int64(2), Int64(3), Int64(4)},
	}

	for idx, in := range testCasesInt64ValueSlice {
		if in == nil {
			continue
		}
		out := Int64ValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := Int64Slice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestInt64Map(t *testing.T) {
	testCasesInt64Map := []map[string]int64{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesInt64Map {
		if in == nil {
			continue
		}
		out := Int64Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Int64ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestFloat32Slice(t *testing.T) {
	testCasesFloat32Slice := [][]float32{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesFloat32Slice {
		if in == nil {
			continue
		}

		out := Float32Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Float32ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestFloat64Slice(t *testing.T) {
	testCasesFloat64Slice := [][]float64{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesFloat64Slice {
		if in == nil {
			continue
		}
		out := Float64Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUintSlice(t *testing.T) {
	testCasesUintSlice := [][]uint{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesUintSlice {
		if in == nil {
			continue
		}
		out := UintSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := UintValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUintValueSlice(t *testing.T) {
	testCasesUintValueSlice := [][]*uint{}

	for idx, in := range testCasesUintValueSlice {
		if in == nil {
			continue
		}
		out := UintValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := UintSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUintMap(t *testing.T) {
	testCasesUintMap := []map[string]uint{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesUintMap {
		if in == nil {
			continue
		}
		out := UintMap(in)
		assertValues(t, in, out, true, idx)

		out2 := UintValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint16Slice(t *testing.T) {
	testCasesUint16Slice := [][]uint16{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesUint16Slice {
		if in == nil {
			continue
		}

		out := Uint16Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint16ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint16ValueSlice(t *testing.T) {
	testCasesUint16ValueSlice := [][]*uint16{}

	for idx, in := range testCasesUint16ValueSlice {
		if in == nil {
			continue
		}

		out := Uint16ValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint16Slice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint16Map(t *testing.T) {
	testCasesUint16Map := []map[string]uint16{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesUint16Map {
		if in == nil {
			continue
		}

		out := Uint16Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint16ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint64Slice(t *testing.T) {
	testCasesUint64Slice := [][]uint64{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesUint64Slice {
		if in == nil {
			continue
		}
		out := Uint64Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint64ValueSlice(t *testing.T) {
	testCasesUint64ValueSlice := [][]*uint64{}

	for idx, in := range testCasesUint64ValueSlice {
		if in == nil {
			continue
		}
		out := Uint64ValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64Slice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint64Map(t *testing.T) {
	testCasesUint64Map := []map[string]uint64{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesUint64Map {
		if in == nil {
			continue
		}
		out := Uint64Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestFloat32ValueSlice(t *testing.T) {
	testCasesFloat32ValueSlice := [][]*float32{}

	for idx, in := range testCasesFloat32ValueSlice {
		if in == nil {
			continue
		}

		out := Float32ValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := Float32Slice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestFloat32Map(t *testing.T) {
	testCasesFloat32Map := []map[string]float32{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesFloat32Map {
		if in == nil {
			continue
		}

		out := Float32Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Float32ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestFloat64ValueSlice(t *testing.T) {
	testCasesFloat64ValueSlice := [][]*float64{}

	for idx, in := range testCasesFloat64ValueSlice {
		if in == nil {
			continue
		}
		out := Float64ValueSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64Slice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestFloat64Map(t *testing.T) {
	testCasesFloat64Map := []map[string]float64{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesFloat64Map {
		if in == nil {
			continue
		}
		out := Float64Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestTimeSlice(t *testing.T) {
	testCasesTimeSlice := [][]time.Time{
		{time.Now(), time.Now().AddDate(100, 0, 0)},
	}

	for idx, in := range testCasesTimeSlice {
		if in == nil {
			continue
		}
		out := TimeSlice(in)
		assertValues(t, in, out, true, idx)

		out2 := TimeValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestTimeValueSlice(t *testing.T) {
	testCasesTimeValueSlice := [][]*time.Time{
		{Time(time.Now()), Time(time.Now().AddDate(100, 0, 0))},
	}

	for idx, in := range testCasesTimeValueSlice {
		if in == nil {
			continue
		}
		out := TimeValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := TimeSlice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestTimeMap(t *testing.T) {
	testCasesTimeMap := []map[string]time.Time{
		{"a": time.Now().AddDate(-100, 0, 0), "b": time.Now()},
	}

	for idx, in := range testCasesTimeMap {
		if in == nil {
			continue
		}
		out := TimeMap(in)
		assertValues(t, in, out, true, idx)

		out2 := TimeValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestInt32Slice(t *testing.T) {
	testCasesInt32Slice := [][]int32{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesInt32Slice {
		if in == nil {
			continue
		}
		out := Int32Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Int32ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestInt32ValueSlice(t *testing.T) {
	testCasesInt32ValueSlice := [][]*int32{
		{Int32(1), Int32(2), Int32(3), Int32(4)},
	}

	for idx, in := range testCasesInt32ValueSlice {
		if in == nil {
			continue
		}
		out := Int32ValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := Int32Slice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestInt32Map(t *testing.T) {
	testCasesInt32Map := []map[string]int32{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesInt32Map {
		if in == nil {
			continue
		}
		out := Int32Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Int32ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint32Slice(t *testing.T) {
	testCasesUint32Slice := [][]uint32{
		{1, 2, 3, 4},
	}

	for idx, in := range testCasesUint32Slice {
		if in == nil {
			continue
		}
		out := Uint32Slice(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint32ValueSlice(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestUint32ValueSlice(t *testing.T) {
	testCasesUint32ValueSlice := [][]*uint32{
		{Uint32(1), Uint32(2), Uint32(3), Uint32(4)},
	}

	for idx, in := range testCasesUint32ValueSlice {
		if in == nil {
			continue
		}
		out := Uint32ValueSlice(in)
		assertValues(t, in, out, false, idx)

		out2 := Uint32Slice(out)
		assertValues(t, in, out2, true, idx)
	}
}

func TestUint32Map(t *testing.T) {
	testCasesUint32Map := []map[string]uint32{
		{"a": 3, "b": 2, "c": 1},
	}

	for idx, in := range testCasesUint32Map {
		if in == nil {
			continue
		}
		out := Uint32Map(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint32ValueMap(out)
		assertValues(t, in, out2, false, idx)
	}
}

func TestStringValue(t *testing.T) {
	testCasesString := []string{"a", "b", "c", "d", "e", ""}

	for idx, in := range testCasesString {
		out := String(in)
		assertValues(t, in, out, true, idx)

		out2 := StringValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, StringValue(nil), "expected conversion from nil to return zero value")
}

func TestBoolValue(t *testing.T) {
	testCasesBool := []bool{true, false}

	for idx, in := range testCasesBool {
		out := Bool(in)
		assertValues(t, in, out, true, idx)

		out2 := BoolValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, BoolValue(nil), "expected conversion from nil to return zero value")
}

func TestIntValue(t *testing.T) {
	testCasesInt := []int{1, 2, 3, 0}

	for idx, in := range testCasesInt {
		out := Int(in)
		assertValues(t, in, out, true, idx)

		out2 := IntValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, IntValue(nil), "expected conversion from nil to return zero value")
}

func TestInt32Value(t *testing.T) {
	testCasesInt32 := []int32{1, 2, 3, 0}

	for idx, in := range testCasesInt32 {
		out := Int32(in)
		assertValues(t, in, out, true, idx)

		out2 := Int32Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Int32Value(nil), "expected conversion from nil to return zero value")
}

func TestInt64Value(t *testing.T) {
	testCasesInt64 := []int64{1, 2, 3, 0}

	for idx, in := range testCasesInt64 {
		out := Int64(in)
		assertValues(t, in, out, true, idx)

		out2 := Int64Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Int64Value(nil), "expected conversion from nil to return zero value")
}

func TestUintValue(t *testing.T) {
	testCasesUint := []uint{1, 2, 3, 0}

	for idx, in := range testCasesUint {
		out := Uint(in)
		assertValues(t, in, out, true, idx)

		out2 := UintValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, UintValue(nil), "expected conversion from nil to return zero value")
}

func TestUint32Value(t *testing.T) {
	testCasesUint32 := []uint32{1, 2, 3, 0}

	for idx, in := range testCasesUint32 {
		out := Uint32(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint32Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Uint32Value(nil), "expected conversion from nil to return zero value")
}

func TestUint64Value(t *testing.T) {
	testCasesUint64 := []uint64{1, 2, 3, 0}

	for idx, in := range testCasesUint64 {
		out := Uint64(in)
		assertValues(t, in, out, true, idx)

		out2 := Uint64Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Uint64Value(nil), "expected conversion from nil to return zero value")
}

func TestFloat32Value(t *testing.T) {
	testCasesFloat32 := []float32{1, 2, 3, 0}

	for idx, in := range testCasesFloat32 {
		out := Float32(in)
		assertValues(t, in, out, true, idx)

		out2 := Float32Value(out)
		assertValues(t, in, out2, false, idx)
	}

	assert.Zerof(t, Float32Value(nil), "expected conversion from nil to return zero value")
}

func TestFloat64Value(t *testing.T) {
	testCasesFloat64 := []float64{1, 2, 3, 0}

	for idx, in := range testCasesFloat64 {
		out := Float64(in)
		assertValues(t, in, out, true, idx)

		out2 := Float64Value(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, Float64Value(nil), "expected conversion from nil to return zero value")
}

func TestTimeValue(t *testing.T) {
	testCasesTime := []time.Time{
		time.Now().AddDate(-100, 0, 0), time.Now(),
	}

	for idx, in := range testCasesTime {
		out := Time(in)
		assertValues(t, in, out, true, idx)

		out2 := TimeValue(out)
		assertValues(t, in, out2, false, idx)
	}
	assert.Zerof(t, TimeValue(nil), "expected conversion from nil to return zero value")
}

func assertSingleValue(t *testing.T, inElem, elem reflect.Value, expectPointer bool, idx int) {
	require.Equalf(t,
		expectPointer, (elem.Kind() == reflect.Ptr),
		"unexpected expectPointer=%t value type", expectPointer,
	)

	if inElem.Kind() == reflect.Ptr && !inElem.IsNil() {
		inElem = reflect.Indirect(inElem)
	}

	if elem.Kind() == reflect.Ptr && !elem.IsNil() {
		elem = reflect.Indirect(elem)
	}

	require.Truef(t,
		(elem.Kind() == reflect.Ptr && elem.IsNil()) ||
			IsZero(elem.Interface()) == (inElem.Kind() == reflect.Ptr && inElem.IsNil()) ||
			IsZero(inElem.Interface()),
		"unexpected nil pointer at idx %d", idx,
	)

	if !((elem.Kind() == reflect.Ptr && elem.IsNil()) || IsZero(elem.Interface())) {
		require.IsTypef(t, inElem.Interface(), elem.Interface(), "Expected in/out to match types")
		assert.EqualValuesf(t, inElem.Interface(), elem.Interface(), "Unexpected value at idx %d: %v", idx, elem.Interface())
	}
}

// assertValues checks equivalent representation pointer vs values for single var, slices and maps
func assertValues(t *testing.T, in, out interface{}, expectPointer bool, idx int) {
	vin := reflect.ValueOf(in)
	vout := reflect.ValueOf(out)

	switch vin.Kind() { //nolint:exhaustive
	case reflect.Slice, reflect.Map:
		require.Equalf(t, vin.Kind(), vout.Kind(), "Unexpected output type at idx %d", idx)
		require.Equalf(t, vin.Len(), vout.Len(), "Unexpected len at idx %d", idx)

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
