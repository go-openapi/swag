package mangling

import (
	"bytes"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestLexemEdgeCases(t *testing.T) {
	t.Run("with single rune, letter but not casable", func(t *testing.T) {
		const input = "आ"
		l := newCasualNameLexem(input)
		b := bytes.Buffer{}

		t.Run("should not titleize", func(t *testing.T) {
			ok := l.WriteTitleized(&b, true)
			assert.False(t, ok)
			assert.Empty(t, b.Bytes())
		})
		t.Run("should not lower", func(t *testing.T) {
			ok := l.WriteLower(&b, true)
			assert.False(t, ok)
			assert.Empty(t, b.Bytes())
		})
	})

	t.Run("with single rune, not letter", func(t *testing.T) {
		const input = "१"
		l := newCasualNameLexem(input)
		b := bytes.Buffer{}

		t.Run("should not titleize", func(t *testing.T) {
			ok := l.WriteTitleized(&b, true)
			assert.False(t, ok)
			assert.Empty(t, b.Bytes())
		})
		t.Run("should not lower", func(t *testing.T) {
			ok := l.WriteLower(&b, true)
			assert.False(t, ok)
		})
	})

	t.Run("with empty lexem", func(t *testing.T) {
		const input = ""
		l := newCasualNameLexem(input)
		b := bytes.Buffer{}

		t.Run("should titleize but do nothing", func(t *testing.T) {
			ok := l.WriteTitleized(&b, true)
			assert.True(t, ok)
			assert.Empty(t, b.Bytes())
		})
		t.Run("should not lower but do nothing", func(t *testing.T) {
			ok := l.WriteLower(&b, true)
			assert.True(t, ok)
			assert.Empty(t, b.Bytes())
		})
	})
}
