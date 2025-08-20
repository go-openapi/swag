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
	"math/bits"
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
	t.Run("should not be integers", testNotIntegers(IsFloat64AJSONInteger, false))
	t.Run("should be integers", testIntegers(IsFloat64AJSONInteger, false))
}

func TestPreviousIsFloat64AJSONInteger(t *testing.T) {
	t.Run("should not be integers", testNotIntegers(previousIsFloat64JSONInteger, false))
	t.Run("should be integers", testIntegers(previousIsFloat64JSONInteger, true))
}

func TestBitWiseIsFloat64AJSONInteger(t *testing.T) {
	t.Run("should not be integers", testNotIntegers(bitwiseIsFloat64JSONInteger, false))
	t.Run("should be integers", testIntegers(bitwiseIsFloat64JSONInteger, false))
}

func TestBitWise2IsFloat64AJSONInteger(t *testing.T) {
	t.Run("should not be integers", testNotIntegers(bitwiseIsFloat64JSONInteger2, false))
	t.Run("should be integers", testIntegers(bitwiseIsFloat64JSONInteger2, false))
}

func TestStdlib2IsFloat64AJSONInteger(t *testing.T) {
	t.Run("should not be integers", testNotIntegers(stdlibIsFloat64JSONInteger, true))
	t.Run("should be integers", testIntegers(stdlibIsFloat64JSONInteger, true))
}

func testNotIntegers(fn func(float64) bool, skipKnownFailure bool) func(*testing.T) {
	_ = skipKnownFailure

	return func(t *testing.T) {
		assert.False(t, fn(math.Inf(1)))
		assert.False(t, fn(maxJSONFloat+1))
		assert.False(t, fn(minJSONFloat-1))
		assert.False(t, fn(math.SmallestNonzeroFloat64))
		assert.False(t, fn(0.5))
		assert.False(t, fn(0.25))
		assert.False(t, fn(1.00/func() float64 { return 2.00 }()))
		assert.False(t, fn(1.00/func() float64 { return 4.00 }()))
		assert.False(t, fn(epsilon))
	}
}

func testIntegers(fn func(float64) bool, skipKnownFailure bool) func(*testing.T) {
	// wrapping in a function forces non-constant evaluation to test float64 rounding behavior
	return func(t *testing.T) {
		assert.True(t, fn(0.0))
		assert.True(t, fn(1.0))
		assert.True(t, fn(maxJSONFloat))
		assert.True(t, fn(minJSONFloat))
		if !skipKnownFailure {
			assert.True(t, fn(1/0.01*67.15000001))
		}
		if !skipKnownFailure {
			assert.True(t, fn(1.00/func() float64 { return 0.01 }()*4643.4))
		}
		assert.True(t, fn(1.00/func() float64 { return 1.00 / 3.00 }()))
		assert.True(t, fn(math.SmallestNonzeroFloat64/2))
		assert.True(t, fn(math.SmallestNonzeroFloat64/3))
		assert.True(t, fn(math.SmallestNonzeroFloat64/4))
	}
}

func BenchmarkIsFloat64JSONInteger(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	b.Run("new float vs integer comparison", benchmarkIsFloat64JSONInteger(IsFloat64AJSONInteger))
	b.Run("previous float vs integer comparison", benchmarkIsFloat64JSONInteger(previousIsFloat64JSONInteger))
	b.Run("bitwise float vs integer comparison", benchmarkIsFloat64JSONInteger(bitwiseIsFloat64JSONInteger))
	b.Run("bitwise float vs integer comparison (2)", benchmarkIsFloat64JSONInteger(bitwiseIsFloat64JSONInteger2))
	b.Run("stdlib float vs integer comparison (2)", benchmarkIsFloat64JSONInteger(stdlibIsFloat64JSONInteger))
}

func BenchmarkBitwise(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	b.Run("bitwise float vs integer comparison (2)", benchmarkIsFloat64JSONInteger(bitwiseIsFloat64JSONInteger2))
}

func previousIsFloat64JSONInteger(f float64) bool {
	if math.IsNaN(f) || math.IsInf(f, 0) || f < minJSONFloat || f > maxJSONFloat {
		return false
	}
	fa := math.Abs(f)
	g := float64(uint64(f))
	ga := math.Abs(g)

	diff := math.Abs(f - g)

	// more info: https://floating-point-gui.de/errors/comparison/#look-out-for-edge-cases
	switch {
	case f == g: // best case
		return true
	case f == float64(int64(f)) || f == float64(uint64(f)): // optimistic case
		return true
	case f == 0 || g == 0 || diff < math.SmallestNonzeroFloat64: // very close to 0 values
		return diff < (epsilon * math.SmallestNonzeroFloat64)
	}
	// check the relative error
	return diff/math.Min(fa+ga, math.MaxFloat64) < epsilon
}

func stdlibIsFloat64JSONInteger(f float64) bool {
	if f < minJSONFloat || f > maxJSONFloat {
		return false
	}
	var bf big.Float
	bf.SetFloat64(f)

	return bf.IsInt()
}

func bitwiseIsFloat64JSONInteger(f float64) bool {
	if math.IsNaN(f) || math.IsInf(f, 0) || f < minJSONFloat || f > maxJSONFloat {
		return false
	}

	mant, exp := math.Frexp(f) // get normalized mantissa
	if exp == 0 && mant == 0 {
		return true
	}
	if exp <= 0 {
		return false
	}

	zeros := bits.TrailingZeros64(uint64(mant))

	return bits.UintSize-zeros <= exp
}

func bitwiseIsFloat64JSONInteger2(f float64) bool {
	if f == 0 {
		return true
	}

	if f < minJSONFloat || f > maxJSONFloat || f != f || f < -math.MaxFloat64 || f > math.MaxFloat64 {
		return false
	}

	// inlined
	var (
		mant uint64
		exp  int
	)
	{
		const smallestNormal = 2.2250738585072014e-308 // 2**-1022

		if math.Abs(f) < smallestNormal {
			f *= (1 << shift) // x 2^52
			exp = -shift
		}

		x := math.Float64bits(f)
		exp += int((x>>shift)&mask) - bias + 1 //nolint:gosec // x>>12 & 0x7FF - 1022 : extract exp, recentered from bias

		x &^= mask << shift       // x= x &^ 0x7FF << 12 (clear 11 exp bits then shift 12)
		x |= (-1 + bias) << shift // x = x | 1022 << 12 ==> or with 1022 as exp location
		mant = uint64(math.Float64frombits(x))
	}
	/*
		{
			x := math.Float64bits(f)
			exp = int(x>>shift) & mask

			if exp < bias {
			} else if exp < bias+shift { // 1023 + 12
				exp -= bias
			}
		}
	*/
	/*
		e := uint(bits>>shift) & mask
		if e < bias {
			// Round abs(x) < 1 including denormals.
			bits &= signMask // +-0
			if e == bias-1 {
				bits |= uvone // +-1
			}
		} else if e < bias+shift {
			// Round any abs(x) >= 1 containing a fractional component [0,1).
			//
			// Numbers with larger exponents are returned unchanged since they
			// must be either an integer, infinity, or NaN.
			const half = 1 << (shift - 1)
			e -= bias
			bits += half >> e
			bits &^= fracMask >> e
		}
	*/

	// It returns frac and exp satisfying f == frac × 2**exp,
	// with the absolute value of frac in the interval [½, 1).
	if exp <= 0 {
		return false
	}

	zeros := bits.TrailingZeros64(mant)

	return bits.UintSize-zeros <= exp
}

const (
	mask  = 0x7FF
	shift = 64 - 11 - 1
	// uvinf    = 0x7FF0000000000000
	// uvneginf = 0xFFF0000000000000
	bias     = 1023
	fracMask = 1<<shift - 1
)

/*
func isNaN(x uint64) bool { // f != f
	return uint32(x>>shift)&mask == mask // && x != uvinf && x != uvneginf
}

func isInf(x uint64) bool { // f < - math.MaxFloat || f > math.MaxFloat
	return x == uvinf || x == uvneginf
}
*/

/*
func frexp(f float64) (frac uint64, exp int) {
	const smallestNormal = 2.2250738585072014e-308 // 2**-1022
	g := f

	if math.Abs(f) < smallestNormal {
		g *= (1 << 52)
		exp = -52
	}

	x := math.Float64bits(g)
	exp += int((x>>shift)&mask) - bias + 1
	x &^= mask << shift
	x |= (-1 + bias) << shift
	frac = uint64(math.Float64frombits(x))

	return
}
*/

func benchmarkIsFloat64JSONInteger(fn func(float64) bool) func(*testing.B) {
	assertCode := func() {
		panic("unexpected result during benchmark")
	}

	return func(b *testing.B) {
		testFunc := func() {
			if fn(math.Inf(1)) {
				assertCode()
			}
			if fn(maxJSONFloat + 1) {
				assertCode()
			}
			if fn(minJSONFloat - 1) {
				assertCode()
			}
			if fn(math.SmallestNonzeroFloat64) {
				assertCode()
			}
			if fn(0.5) {
				assertCode()
			}

			if !fn(1.0) {
				assertCode()
			}
			if !fn(maxJSONFloat) {
				assertCode()
			}
			if !fn(minJSONFloat) {
				assertCode()
			}
			if !fn(1 / 0.01 * 67.15000001) {
				assertCode()
			}
			/* can't compare both versions on this test case
			if !fn(1 / func() float64 { return 0.01 }() * 4643.4) {
				assertCode()
			}
			*/
			if !fn(math.SmallestNonzeroFloat64 / 2) {
				assertCode()
			}
			if !fn(math.SmallestNonzeroFloat64 / 3) {
				assertCode()
			}
			if !fn(math.SmallestNonzeroFloat64 / 4) {
				assertCode()
			}
		}

		for n := 0; n < b.N; n++ {
			testFunc()
		}
	}
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
