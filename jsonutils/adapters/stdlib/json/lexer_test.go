// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	stdjson "encoding/json"
	"io"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestToken(t *testing.T) {
	t.Run("token should be stringable for debugging and error formatting", func(t *testing.T) {
		assert.EqualT(t, "invalid token", invalidToken.String())
		assert.EqualT(t, "EOF", eofToken.String())
	})

	t.Run("token should be able to map all tokens from encoding/json.Token", func(t *testing.T) {
		stdtok := stdjson.Token(stdjson.Number("123"))
		tok := token{
			Token: stdtok,
		}
		assert.EqualT(t, tokenNumber, tok.Kind())
	})

	t.Run("token should detect JSON delimiters", func(t *testing.T) {
		tok := eofToken
		assert.Zero(t, tok.Delim())
		tok = token{
			Token: stdjson.Delim(','),
		}
		assert.EqualT(t, byte(','), tok.Delim())
	})
}

func TestBytesReader(t *testing.T) {
	t.Run("should read into small buffer", func(t *testing.T) {
		r := &bytesReader{
			buf: []byte("1234567890"),
		}

		const bufferSize = 3

		buf := make([]byte, bufferSize)

		n, err := r.Read(buf)
		require.NoError(t, err)
		require.EqualT(t, bufferSize, n)
		require.Equal(t, []byte("123"), buf[:n])
		assert.EqualT(t, bufferSize, r.offset)

		n, err = r.Read(buf)
		require.NoError(t, err)
		require.EqualT(t, bufferSize, n)
		require.Equal(t, []byte("456"), buf[:n])
		assert.EqualT(t, 2*bufferSize, r.offset)

		n, err = r.Read(buf)
		require.NoError(t, err)
		require.EqualT(t, bufferSize, n)
		require.Equal(t, []byte("789"), buf[:n])
		assert.EqualT(t, 3*bufferSize, r.offset)

		n, err = r.Read(buf)
		require.NoError(t, err)
		require.EqualT(t, 1, n)
		require.Equal(t, []byte("0"), buf[:n])
		assert.EqualT(t, len(r.buf), r.offset)

		n, err = r.Read(buf)
		require.EqualT(t, 0, n)
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("should read into large buffer", func(t *testing.T) {
		r := &bytesReader{
			buf: []byte("1234567890"),
		}

		const bufferSize = 12

		buf := make([]byte, bufferSize)

		n, err := r.Read(buf)
		require.NoError(t, err)
		require.EqualT(t, len(r.buf), n)
		require.Equal(t, r.buf, buf[:n])
		assert.EqualT(t, len(r.buf), r.offset)
	})
}

func TestLexer(t *testing.T) {
	t.Run("lexer should be interruptible by setting error state", func(t *testing.T) {
		l := newLexer([]byte("123"))
		l.SetErr(ErrStdlib)
		require.FalseT(t, l.Ok())
		require.Error(t, l.Error())

		require.Equal(t, invalidToken, l.NextToken())
		require.FalseT(t, l.IsDelim(','))
		require.FalseT(t, l.IsNull())
		require.Zero(t, l.Number())
		require.NotPanics(t, func() {
			l.Null()
		})
	})

	t.Run("lexer should detect delimiter (comma and colon are elided)", func(t *testing.T) {
		l := newLexer([]byte{})
		l.next = token{Token: stdjson.Delim('{')}

		l.Delim('{')
		require.TrueT(t, l.Ok())

		l.next = token{Token: "123"}
		l.Delim('{')
		require.FalseT(t, l.Ok())
	})

	t.Run("lexer should detect null", func(t *testing.T) {
		l := newLexer([]byte{})
		l.next = token{Token: nil}

		l.Null()
		require.TrueT(t, l.Ok())

		l.next = token{Token: "123"}
		l.Null()
		require.FalseT(t, l.Ok())
	})

	t.Run("lexer should detect bool", func(t *testing.T) {
		l := newLexer([]byte{})
		l.next = token{Token: false}

		b := l.Bool()
		require.TrueT(t, l.Ok())
		require.FalseT(t, b)

		l.next = token{Token: true}
		b = l.Bool()
		require.TrueT(t, l.Ok())
		require.TrueT(t, b)

		l.next = token{Token: "x"}
		b = l.Bool()
		require.FalseT(t, l.Ok())
		require.FalseT(t, b)
	})

	t.Run("lexer should detect JSON number as string", func(t *testing.T) {
		const epsilon = 1e-9

		l := newLexer([]byte{})
		l.next = token{Token: stdjson.Number("123")}

		n := l.Number()
		require.TrueT(t, l.Ok())
		require.Equal(t, int64(123), n)

		l.next = token{Token: stdjson.Number("123.4")}
		n = l.Number()
		require.TrueT(t, l.Ok())
		require.InDelta(t, float64(123.4), n, epsilon)

		l.next = token{Token: 123.4}
		n = l.Number()
		require.TrueT(t, l.Ok())
		require.InDelta(t, float64(123.4), n, epsilon)

		l.next = token{Token: "123.4"}
		n = l.Number()
		require.FalseT(t, l.Ok())
		require.Zero(t, n)
	})
}
