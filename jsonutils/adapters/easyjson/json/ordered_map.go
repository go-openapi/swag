// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"iter"
	"math"
	"strconv"

	"github.com/go-openapi/swag/conv"
	"github.com/go-openapi/swag/jsonutils"
	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

var _ ifaces.OrderedMap = &MapSlice{}

// MapSlice represents a JSON object, with the order of keys maintained.
//
// It implements [ifaces.Ordered] and [ifaces.SetOrdered].
type MapSlice []MapItem

func (s MapSlice) OrderedItems() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for _, item := range s {
			if !yield(item.Key, item.Value) {
				return
			}
		}
	}
}

func (s *MapSlice) SetOrderedItems(items iter.Seq2[string, any]) {
	if items == nil {
		*s = nil

		return
	}

	m := *s
	if len(m) > 0 {
		// update mode
		idx := make(map[string]int, len(m))

		for i, item := range m {
			idx[item.Key] = i
		}

		for k, v := range items {
			idx, ok := idx[k]
			if ok {
				m[idx].Value = v

				continue
			}

			m = append(m, MapItem{Key: k, Value: v})
		}

		*s = m

		return
	}

	for k, v := range items {
		m = append(m, MapItem{Key: k, Value: v})
	}

	*s = m
}

// MarshalJSON renders a [MapSlice] as JSON bytes, preserving the order of keys.
func (s MapSlice) MarshalJSON() ([]byte, error) {
	return s.OrderedMarshalJSON()
}

func (s MapSlice) OrderedMarshalJSON() ([]byte, error) {
	w, redeem := BorrowWriter()
	defer redeem()

	s.MarshalEasyJSON(w)

	return w.BuildBytes() // this actually copies data, so its okay to redeem the writer
}

// MarshalEasyJSON renders a [MapSlice] as JSON bytes, using easyJSON
func (s MapSlice) MarshalEasyJSON(w *jwriter.Writer) {
	s.marshalEasyJSON(w, defaultMaxNestingDepth)
}

// UnmarshalJSON builds a [MapSlice] from JSON bytes, preserving the order of keys.
//
// Inner objects are unmarshaled as [MapSlice] slices and not map[string]any.
func (s *MapSlice) UnmarshalJSON(data []byte) error {
	return s.OrderedUnmarshalJSON(data)
}

func (s *MapSlice) OrderedUnmarshalJSON(data []byte) error {
	return s.orderedUnmarshalJSON(data, defaultMaxNestingDepth)
}

// UnmarshalEasyJSON builds a [MapSlice] from JSON bytes, using easyJSON
func (s *MapSlice) UnmarshalEasyJSON(in *jlexer.Lexer) {
	s.unmarshalEasyJSON(in, defaultMaxNestingDepth)
}

func (s MapSlice) marshalEasyJSON(w *jwriter.Writer, budget int) {
	if s == nil {
		w.RawString("null")

		return
	}

	if budget <= 0 {
		w.Error = ErrMaxNestingDepth

		return
	}

	w.RawByte('{')

	if len(s) == 0 {
		w.RawByte('}')

		return
	}

	s[0].marshalEasyJSON(w, budget)

	for i := 1; i < len(s); i++ {
		w.RawByte(',')
		s[i].marshalEasyJSON(w, budget)
	}

	w.RawByte('}')
}

func (s *MapSlice) orderedUnmarshalJSON(data []byte, budget int) error {
	l, redeem := BorrowLexer(data)
	defer redeem()

	s.unmarshalEasyJSON(l, budget)

	return l.Error()
}

// unmarshalEasyJSON parses an object, decreasing budget for every nested container
// so that adversarially deep documents return an error instead of overflowing the stack.
func (s *MapSlice) unmarshalEasyJSON(in *jlexer.Lexer, budget int) {
	if in.IsNull() {
		in.Skip()

		return
	}

	if budget <= 0 {
		in.AddError(ErrMaxNestingDepth)

		return
	}

	result := make(MapSlice, 0)
	in.Delim('{')
	for in.Ok() && !in.IsDelim('}') {
		var mi MapItem
		mi.unmarshalEasyJSON(in, budget)
		result = append(result, mi)
	}
	in.Delim('}')

	*s = result
}

// MapItem represents the value of a key in a JSON object held by [MapSlice].
//
// Notice that MapItem should not be marshaled to or unmarshaled from JSON directly,
// use this type as part of a [MapSlice] when dealing with JSON bytes.
type MapItem struct {
	Key   string
	Value any
}

// MarshalEasyJSON renders a [MapItem] as JSON bytes, using easyJSON
func (s MapItem) MarshalEasyJSON(w *jwriter.Writer) {
	s.marshalEasyJSON(w, defaultMaxNestingDepth)
}

// UnmarshalEasyJSON builds a [MapItem] from JSON bytes, using easyJSON
func (s *MapItem) UnmarshalEasyJSON(in *jlexer.Lexer) {
	s.unmarshalEasyJSON(in, defaultMaxNestingDepth)
}

func (s MapItem) marshalEasyJSON(w *jwriter.Writer, budget int) {
	w.String(s.Key)
	w.RawByte(':')

	// Recurse internally for nested ordered maps so the depth guard survives, instead
	// of dispatching to MarshalEasyJSON which would reset it and re-enable overflow.
	if nested, ok := s.Value.(MapSlice); ok {
		nested.marshalEasyJSON(w, budget-1)

		return
	}

	if val, ok := s.Value.(easyjson.Marshaler); ok {
		val.MarshalEasyJSON(w)

		return
	}

	w.Raw(jsonutils.WriteJSON(s.Value))
}

func (s *MapItem) unmarshalEasyJSON(in *jlexer.Lexer, budget int) {
	key := in.UnsafeString()
	in.WantColon()
	value := s.asInterface(in, budget-1)
	in.WantComma()

	s.Key = key
	s.Value = value
}

// asInterface is very much like [jlexer.Lexer.Interface], but unmarshals an object
// into a [MapSlice], not a map[string]any.
//
// We have to force parsing errors somehow, since [jlexer.Lexer] doesn't let us
// set a parsing error directly.
func (s *MapItem) asInterface(in *jlexer.Lexer, budget int) any {
	tokenKind := in.CurrentToken()

	if !in.Ok() {
		return nil
	}

	switch tokenKind {
	case jlexer.TokenString:
		return in.String()

	case jlexer.TokenNumber:
		// determine if we may use an integer type
		n := in.JsonNumber().String()
		f, _ := strconv.ParseFloat(n, 64)
		if conv.IsFloat64AJSONInteger(f) {
			return int64(math.Trunc(f))
		}
		return f

	case jlexer.TokenBool:
		return in.Bool()

	case jlexer.TokenNull:
		in.Null()
		return nil

	case jlexer.TokenDelim:
		if in.IsDelim('{') {
			ret := make(MapSlice, 0)
			ret.unmarshalEasyJSON(in, budget)

			if in.Ok() {
				return ret
			}

			// lexer is in an error state: will exhaust
			return nil
		}

		if in.IsDelim('[') {
			if budget <= 0 {
				in.AddError(ErrMaxNestingDepth)

				return nil
			}

			in.Delim('[') // consume

			ret := []any{}
			for in.Ok() && !in.IsDelim(']') {
				ret = append(ret, s.asInterface(in, budget-1))
				in.WantComma()
			}
			in.Delim(']')

			if in.Ok() {
				return ret
			}

			// lexer is in an error state: will exhaust
			return nil
		}

		if in.Ok() {
			in.Delim('{') // force error
		}

		return nil

	case jlexer.TokenUndef:
		fallthrough
	default:
		if in.Ok() {
			in.Delim('{') // force error
		}

		return nil
	}
}
