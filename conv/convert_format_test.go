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
	"fmt"
	"io"
	"math"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var evaluatesAsTrue = map[string]struct{}{
	"true":     {},
	"1":        {},
	"yes":      {},
	"ok":       {},
	"y":        {},
	"on":       {},
	"selected": {},
	"checked":  {},
	"t":        {},
	"enabled":  {},
}

func TestConvertBool(t *testing.T) {
	for k := range evaluatesAsTrue {
		r, err := ConvertBool(k)
		require.NoError(t, err)
		assert.True(t, r)
	}
	for _, k := range []string{"a", "", "0", "false", "unchecked", "anythingElse"} {
		r, err := ConvertBool(k)
		require.NoError(t, err)
		assert.False(t, r)
	}
}

func TestFormatBool(t *testing.T) {
	assert.Equal(t, "true", FormatBool(true))
	assert.Equal(t, "false", FormatBool(false))
}

func TestConvertFloat(t *testing.T) {
	t.Run("with float32", func(t *testing.T) {
		validFloats := []float32{1.0, -1, math.MaxFloat32, math.SmallestNonzeroFloat32, 0, 5.494430303}
		invalidFloats := []string{"a", strconv.FormatFloat(math.MaxFloat64, 'f', -1, 64), "true", float64OverflowStr()}

		for _, f := range validFloats {
			str := FormatFloat(f)
			c1, err := ConvertFloat32(str)
			require.NoError(t, err)
			assert.InDelta(t, f, c1, 1e-6)

			c2, err := ConvertFloat[float32](str)
			require.NoError(t, err)
			assert.InDelta(t, c1, c2, 1e-6)
		}

		for _, f := range invalidFloats {
			_, err := ConvertFloat32(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertFloat[float32](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with float64", func(t *testing.T) {
		validFloats := []float64{1.0, -1, float64(math.MaxFloat32), float64(math.SmallestNonzeroFloat32), math.MaxFloat64, math.SmallestNonzeroFloat64, 0, 5.494430303}
		invalidFloats := []string{"a", "true", float64OverflowStr()}

		for _, f := range validFloats {
			str := FormatFloat(f)
			c1, err := ConvertFloat64(str)
			require.NoError(t, err)
			assert.InDelta(t, f, c1, 1e-6)

			c2, err := ConvertFloat64(str)
			require.NoError(t, err)
			assert.InDelta(t, c1, c2, 1e-6)
		}

		for _, f := range invalidFloats {
			_, err := ConvertFloat64(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertFloat[float64](f)
			require.Error(t, err, testErrMsg(f))
		}
	})
}

func TestConvertInteger(t *testing.T) {
	t.Run("with int8", func(t *testing.T) {
		validInts := []int8{0, 1, -1, math.MaxInt8, math.MinInt8}
		invalidInts := []string{"1.233", "a", "false", strconv.FormatInt(int64(math.MaxInt64), 10)}

		for _, f := range validInts {
			str := FormatInteger(f)
			c1, err := ConvertInt8(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertInteger[int8](str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidInts {
			_, err := ConvertInt8(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertInteger[int8](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with int16", func(t *testing.T) {
		validInts := []int16{0, 1, -1, math.MaxInt8, math.MinInt8, math.MaxInt16, math.MinInt16}
		invalidInts := []string{"1.233", "a", "false", strconv.FormatInt(int64(math.MaxInt64), 10)}

		for _, f := range validInts {
			str := FormatInteger(f)
			c1, err := ConvertInt16(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertInteger[int16](str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidInts {
			_, err := ConvertInt16(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertInteger[int16](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with int32", func(t *testing.T) {
		validInts := []int32{0, 1, -1, math.MaxInt8, math.MinInt8, math.MaxInt16, math.MinInt16, math.MinInt32, math.MaxInt32}
		invalidInts := []string{"1.233", "a", "false", strconv.FormatInt(int64(math.MaxInt64), 10)}

		for _, f := range validInts {
			str := FormatInteger(f)
			c1, err := ConvertInt32(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertInteger[int32](str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidInts {
			_, err := ConvertInt32(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertInteger[int32](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with int64", func(t *testing.T) {
		validInts := []int64{0, 1, -1, math.MaxInt8, math.MinInt8, math.MaxInt16, math.MinInt16, math.MinInt32, math.MaxInt32, math.MaxInt64, math.MinInt64}
		invalidInts := []string{"1.233", "a", "false"}

		for _, f := range validInts {
			str := FormatInteger(f)
			c1, err := ConvertInt64(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertInt64(str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidInts {
			_, err := ConvertInt64(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertInteger[int64](f)
			require.Error(t, err, testErrMsg(f))
		}
	})
}

func TestConvertUinteger(t *testing.T) {
	t.Run("with uint8", func(t *testing.T) {
		validInts := []uint8{0, 1, math.MaxUint8}
		invalidInts := []string{"1.233", "a", "false", strconv.FormatUint(math.MaxUint64, 10), "-1"}

		for _, f := range validInts {
			str := FormatUinteger(f)
			c1, err := ConvertUint8(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertUinteger[uint8](str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidInts {
			_, err := ConvertUint8(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertUinteger[uint8](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with uint16", func(t *testing.T) {
		validUints := []uint16{0, 1, math.MaxUint8, math.MaxUint16}
		invalidUints := []string{"1.233", "a", "false", strconv.FormatUint(math.MaxUint64, 10), strconv.FormatInt(-1, 10)}

		for _, f := range validUints {
			str := FormatUinteger(f)
			c1, err := ConvertUint16(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertUinteger[uint16](str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidUints {
			_, err := ConvertUint16(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertUinteger[uint16](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with uint32", func(t *testing.T) {
		validUints := []uint32{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32}
		invalidUints := []string{"1.233", "a", "false", strconv.FormatUint(math.MaxUint64, 10), strconv.FormatInt(-1, 10)}

		for _, f := range validUints {
			str := FormatUinteger(f)
			c1, err := ConvertUint32(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertUint32(str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}

		for _, f := range invalidUints {
			_, err := ConvertUint32(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertUinteger[uint32](f)
			require.Error(t, err, testErrMsg(f))
		}
	})

	t.Run("with uint64", func(t *testing.T) {
		validUints := []uint64{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}
		invalidUints := []string{"1.233", "a", "false", strconv.FormatInt(-1, 10), uint64OverflowStr()}

		for _, f := range validUints {
			str := FormatUinteger(f)
			c1, err := ConvertUint64(str)
			require.NoError(t, err)
			assert.Equal(t, f, c1)

			c2, err := ConvertUinteger[uint64](str)
			require.NoError(t, err)
			assert.Equal(t, c1, c2)
		}
		for _, f := range invalidUints {
			_, err := ConvertUint64(f)
			require.Error(t, err, testErrMsg(f))

			_, err = ConvertUinteger[uint64](f)
			require.Error(t, err, testErrMsg(f))
		}
	})
}

func TestIsFloat64AJSONInteger(t *testing.T) {
	assert.False(t, IsFloat64AJSONInteger(math.Inf(1)))
	assert.False(t, IsFloat64AJSONInteger(maxJSONFloat+1))
	assert.False(t, IsFloat64AJSONInteger(minJSONFloat-1))
	assert.False(t, IsFloat64AJSONInteger(math.SmallestNonzeroFloat64))

	assert.True(t, IsFloat64AJSONInteger(1.0))
	assert.True(t, IsFloat64AJSONInteger(maxJSONFloat))
	assert.True(t, IsFloat64AJSONInteger(minJSONFloat))
	assert.True(t, IsFloat64AJSONInteger(1/0.01*67.15000001))
	assert.True(t, IsFloat64AJSONInteger(math.SmallestNonzeroFloat64/2))
	assert.True(t, IsFloat64AJSONInteger(math.SmallestNonzeroFloat64/3))
	assert.True(t, IsFloat64AJSONInteger(math.SmallestNonzeroFloat64/4))
}

// test utilities

func testErrMsg(f string) string {
	const (
		expectedQuote = "expected '"
		errSuffix     = "' to generate an error"
	)

	return expectedQuote + f + errSuffix
}

func uint64OverflowStr() string {
	var one, maxUint, overflow big.Int
	one.SetUint64(1)
	maxUint.SetUint64(math.MaxUint64)
	overflow.Add(&maxUint, &one)

	return overflow.String()
}

func float64OverflowStr() string {
	var one, maxFloat64, overflow big.Float
	one.SetFloat64(1.00)
	maxFloat64.SetFloat64(math.MaxFloat64)
	overflow.Add(&maxFloat64, &one)

	return overflow.String()
}

// benchmarks
func BenchmarkConvertBool(b *testing.B) {
	inputs := []string{
		"a", "t", "ok", "false", "true", "TRUE", "no", "n", "y",
	}
	var isTrue bool

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("use switch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			isTrue, _ = ConvertBool(inputs[i%len(inputs)])
		}
		fmt.Fprintln(io.Discard, isTrue)
	})

	b.Run("use map (previous version)", func(b *testing.B) {
		previousConvertBool := func(str string) (bool, error) {
			_, ok := evaluatesAsTrue[strings.ToLower(str)]
			return ok, nil
		}

		for i := 0; i < b.N; i++ {
			isTrue, _ = previousConvertBool(inputs[i%len(inputs)])
		}
		fmt.Fprintln(io.Discard, isTrue)
	})

	b.Run("use slice.Contains", func(b *testing.B) {
		sliceContainsConvertBool := func(str string) (bool, error) {
			return slices.Contains(
				[]string{"true", "1", "yes", "ok", "y", "on", "selected", "checked", "t", "enabled"},
				strings.ToLower(str),
			), nil
		}

		for i := 0; i < b.N; i++ {
			isTrue, _ = sliceContainsConvertBool(inputs[i%len(inputs)])
		}
		fmt.Fprintln(io.Discard, isTrue)
	})
}
