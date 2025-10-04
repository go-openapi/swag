// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestConvIface(t *testing.T) {
	const epsilon = 1e-6

	t.Run("deprecated Convert functions should work", func(t *testing.T) {
		// only check happy path - more comprehensive testing is carried out inside the called packag
		assert.True(t, IsFloat64AJSONInteger(1.00))

		b, err := ConvertBool("true")
		require.NoError(t, err)
		assert.True(t, b)

		f32, err := ConvertFloat32("1.05")
		require.NoError(t, err)
		assert.InDelta(t, float32(1.05), f32, epsilon)

		f64, err := ConvertFloat64("1.05")
		require.NoError(t, err)
		assert.InDelta(t, float32(1.05), f64, epsilon)

		i8, err := ConvertInt8("2")
		require.NoError(t, err)
		assert.Equal(t, int8(2), i8)

		i16, err := ConvertInt16("2")
		require.NoError(t, err)
		assert.Equal(t, int16(2), i16)

		i32, err := ConvertInt32("2")
		require.NoError(t, err)
		assert.Equal(t, int32(2), i32)

		i64, err := ConvertInt64("2")
		require.NoError(t, err)
		assert.Equal(t, int64(2), i64)

		u8, err := ConvertUint8("2")
		require.NoError(t, err)
		assert.Equal(t, uint8(2), u8)

		u16, err := ConvertUint16("2")
		require.NoError(t, err)
		assert.Equal(t, uint16(2), u16)

		u32, err := ConvertUint32("2")
		require.NoError(t, err)
		assert.Equal(t, uint32(2), u32)

		u64, err := ConvertUint64("2")
		require.NoError(t, err)
		assert.Equal(t, uint64(2), u64)
	})

	t.Run("deprecated Format functions should work", func(t *testing.T) {
		assert.Equal(t, "true", FormatBool(true))
		assert.Equal(t, "1.05", FormatFloat32(1.05))
		assert.Equal(t, "1.05", FormatFloat64(1.05))
		assert.Equal(t, "1", FormatInt8(1))
		assert.Equal(t, "1", FormatInt16(1))
		assert.Equal(t, "1", FormatInt32(1))
		assert.Equal(t, "1", FormatInt64(1))
		assert.Equal(t, "1", FormatUint8(1))
		assert.Equal(t, "1", FormatUint16(1))
		assert.Equal(t, "1", FormatUint32(1))
		assert.Equal(t, "1", FormatUint64(1))
	})

	t.Run("deprecated pointer functions should work", func(t *testing.T) {
		assert.Equal(t, "a", StringValue(String("a")))
		assert.Equal(t, []string{"a"}, StringValueSlice(StringSlice([]string{"a"})))
		assert.Equal(t, map[string]string{"1": "a"}, StringValueMap(StringMap(map[string]string{"1": "a"})))

		assert.True(t, BoolValue(Bool(true)))
		assert.Equal(t, []bool{true}, BoolValueSlice(BoolSlice([]bool{true})))
		assert.Equal(t, map[string]bool{"1": true}, BoolValueMap(BoolMap(map[string]bool{"1": true})))

		assert.Equal(t, 1, IntValue(Int(1)))
		assert.Equal(t, []int{1}, IntValueSlice(IntSlice([]int{1})))
		assert.Equal(t, map[string]int{"1": 1}, IntValueMap(IntMap(map[string]int{"1": 1})))

		assert.Equal(t, int32(1), Int32Value(Int32(1)))
		assert.Equal(t, []int32{1}, Int32ValueSlice(Int32Slice([]int32{1})))
		assert.Equal(t, map[string]int32{"1": 1}, Int32ValueMap(Int32Map(map[string]int32{"1": 1})))

		assert.Equal(t, int64(1), Int64Value(Int64(1)))
		assert.Equal(t, []int64{1}, Int64ValueSlice(Int64Slice([]int64{1})))
		assert.Equal(t, map[string]int64{"1": 1}, Int64ValueMap(Int64Map(map[string]int64{"1": 1})))

		assert.Equal(t, uint16(1), Uint16Value(Uint16(1)))
		assert.Equal(t, []uint16{1}, Uint16ValueSlice(Uint16Slice([]uint16{1})))
		assert.Equal(t, map[string]uint16{"1": 1}, Uint16ValueMap(Uint16Map(map[string]uint16{"1": 1})))

		assert.Equal(t, uint32(1), Uint32Value(Uint32(1)))
		assert.Equal(t, []uint32{1}, Uint32ValueSlice(Uint32Slice([]uint32{1})))
		assert.Equal(t, map[string]uint32{"1": 1}, Uint32ValueMap(Uint32Map(map[string]uint32{"1": 1})))

		assert.Equal(t, uint64(1), Uint64Value(Uint64(1)))
		assert.Equal(t, []uint64{1}, Uint64ValueSlice(Uint64Slice([]uint64{1})))
		assert.Equal(t, map[string]uint64{"1": 1}, Uint64ValueMap(Uint64Map(map[string]uint64{"1": 1})))

		assert.Equal(t, uint(1), UintValue(Uint(1)))
		assert.Equal(t, []uint{1}, UintValueSlice(UintSlice([]uint{1})))
		assert.Equal(t, map[string]uint{"1": 1}, UintValueMap(UintMap(map[string]uint{"1": 1})))

		assert.InDelta(t, float32(1.00), Float32Value(Float32(1.00)), epsilon)
		assert.Equal(t, []float32{1.00}, Float32ValueSlice(Float32Slice([]float32{1.00})))
		assert.Equal(t, map[string]float32{"1": 1.00}, Float32ValueMap(Float32Map(map[string]float32{"1": 1.00})))

		assert.InDelta(t, float64(1.00), Float64Value(Float64(1)), epsilon)
		assert.Equal(t, []float64{1.00}, Float64ValueSlice(Float64Slice([]float64{1.00})))
		assert.Equal(t, map[string]float64{"1": 1.00}, Float64ValueMap(Float64Map(map[string]float64{"1": 1.00})))

		assert.Equal(t, time.Unix(0, 0), TimeValue(Time(time.Unix(0, 0))))
		assert.Equal(t, []time.Time{time.Unix(0, 0)}, TimeValueSlice(TimeSlice([]time.Time{time.Unix(0, 0)})))
		assert.Equal(t, map[string]time.Time{"1": time.Unix(0, 0)}, TimeValueMap(TimeMap(map[string]time.Time{"1": time.Unix(0, 0)})))
	})
}
