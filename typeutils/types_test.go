// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package typeutils

import (
	"testing"
	"time"
	"unsafe"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type SimpleZeroes struct {
	ID   string
	Name string
}

type ZeroesWithTime struct {
	Time time.Time
}

type dummyZeroable struct {
	zero bool
}

func (d dummyZeroable) IsZero() bool {
	return d.zero
}

func TestIsZero(t *testing.T) {
	var strs [5]string
	var strss []string
	var a int
	var b int8
	var c int16
	var d int32
	var e int64
	var f uint
	var g uint8
	var h uint16
	var i uint32
	var j uint64
	var k map[string]string
	var l any
	var m *SimpleZeroes
	var n string
	var o SimpleZeroes
	var p ZeroesWithTime
	var q time.Time
	var z bool
	data := []struct {
		Data     any
		Expected bool
	}{
		{a, true},
		{b, true},
		{c, true},
		{d, true},
		{e, true},
		{f, true},
		{g, true},
		{h, true},
		{i, true},
		{j, true},
		{k, true},
		{l, true},
		{m, true},
		{n, true},
		{o, true},
		{p, true},
		{q, true},
		{strss, true},
		{strs, true},
		{"", true},
		{nil, true},
		{1, false},
		{0, true},
		{int8(1), false},
		{int8(0), true},
		{int16(1), false},
		{int16(0), true},
		{int32(1), false},
		{int32(0), true},
		{int64(1), false},
		{int64(0), true},
		{uint(1), false},
		{uint(0), true},
		{uint8(1), false},
		{uint8(0), true},
		{uint16(1), false},
		{uint16(0), true},
		{uint32(1), false},
		{uint32(0), true},
		{uint64(1), false},
		{uint64(0), true},
		{0.0, true},
		{0.1, false},
		{float32(0.0), true},
		{float32(0.1), false},
		{float64(0.0), true},
		{float64(0.1), false},
		{[...]string{}, true},
		{[...]string{"hello"}, false},
		{[]string(nil), true},
		{[]string{"a"}, false},
		{&dummyZeroable{true}, true},
		{&dummyZeroable{false}, false},
		{(*dummyZeroable)(nil), true},
		{z, true},
	}

	for _, it := range data {
		assert.Equalf(t, it.Expected, IsZero(it.Data), "expected %#v, but got %#v", it.Expected, it.Data)
	}
}

func TestIsNil(t *testing.T) {
	var (
		c     chan<- int
		f     func() bool
		s     []int
		m     map[string]any
		struc struct{}
	)

	for _, value := range []any{
		nil,
		[]string(nil),
		zeroable(nil),
		s,
		m,
		unsafe.Pointer(nil),
		c,
		f,
	} {
		require.True(t, IsNil(value))
	}

	for _, value := range []any{
		[]string{},
		map[string]string{},
		struc,
		0,
		0.00,
		"",
		false,
	} {
		require.False(t, IsNil(value))
	}

}
