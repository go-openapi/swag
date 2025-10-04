// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestWriter(t *testing.T) {
	const (
		marker = "START"
		str    = "should not be written"
	)

	t.Run("writer should be interruptible when in error state", func(t *testing.T) {
		w := newJWriter()
		w.RawString(marker)

		w.err = ErrStdlib
		w.RawString(str)
		require.Equal(t, marker, w.buf.String())

		raw := func() ([]byte, error) { return []byte(str), nil }
		w.Raw(raw())
		require.Equal(t, marker, w.buf.String())

		w.RawByte('x')
		require.Equal(t, marker, w.buf.String())

		w.String(str)
		require.Equal(t, marker, w.buf.String())

		result, err := w.BuildBytes()
		require.Nil(t, result)
		require.ErrorIs(t, err, ErrStdlib)
	})

	t.Run("Raw should not write any output if error", func(t *testing.T) {
		w := newJWriter()
		w.RawString(marker)
		raw := func() ([]byte, error) { return []byte(str), ErrStdlib }
		w.Raw(raw())
		require.Equal(t, marker, w.buf.String())

		result, err := w.BuildBytes()
		require.Nil(t, result)
		require.ErrorIs(t, err, ErrStdlib)
	})
}
