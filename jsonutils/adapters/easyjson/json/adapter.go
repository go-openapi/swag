// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	stdjson "encoding/json"

	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	"github.com/go-openapi/swag/typeutils"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

var _ ifaces.Adapter = &Adapter{}

type Adapter struct {
	options
}

// NewAdapter yields a JSON adapter for [easyjson].
func NewAdapter(opts ...Option) *Adapter {
	var o options
	for _, apply := range opts {
		apply(&o)
	}
	return &Adapter{
		options: o,
	}
}

func (a *Adapter) Marshal(value any) ([]byte, error) {
	marshaler, ok := value.(easyjson.Marshaler)
	if ok {
		w := BorrowWriter()
		defer func() {
			RedeemWriter(w)
		}()
		if a.nilMapAsEmpty {
			w.Flags |= jwriter.NilMapAsEmpty
		}
		if a.nilSliceAsEmpty {
			w.Flags |= jwriter.NilSliceAsEmpty
		}
		w.NoEscapeHTML = a.noEscapeHTML

		marshaler.MarshalEasyJSON(w)

		return w.BuildBytes() // this actually copies data, so its okay to redeem the writer
	}

	// fallback to standard library
	return stdjson.Marshal(value)
}

func (a *Adapter) Unmarshal(data []byte, value any) error {
	unmarshaler, ok := value.(easyjson.Unmarshaler)
	if ok {
		l := BorrowLexer(data)
		defer func() {
			RedeemLexer(l)
		}()
		l.UseMultipleErrors = a.useMultipleErrors

		unmarshaler.UnmarshalEasyJSON(l)
		return l.Error()
	}

	return stdjson.Unmarshal(data, value)
}

func (a *Adapter) OrderedMarshal(value ifaces.Ordered) ([]byte, error) {
	w := BorrowWriter()
	defer func() {
		RedeemWriter(w)
	}()

	if typeutils.IsNil(value) {
		w.RawString("null")

		return w.BuildBytes()
	}

	w.RawByte('{')
	first := true
	for k, v := range value.OrderedItems() {
		if first {
			first = false
		} else {
			w.RawByte(',')
		}

		w.String(k)
		w.RawByte(':')

		switch val := v.(type) {
		case easyjson.Marshaler:
			val.MarshalEasyJSON(w)
		case ifaces.Ordered:
			w.Raw(a.OrderedMarshal(val))
		default:
			w.Raw(stdjson.Marshal(v))
		}
	}

	w.RawByte('}')

	return w.BuildBytes() // this actually copies data, so its okay to redeem the writer
}

func (a *Adapter) OrderedUnmarshal(data []byte, value ifaces.SetOrdered) error {
	var m MapSlice
	if err := m.OrderedUnmarshalJSON(data); err != nil {
		return err
	}

	if typeutils.IsNil(m) {
		// force input value to nil
		value.SetOrderedItems(nil)

		return nil
	}

	value.SetOrderedItems(m.OrderedItems())

	return nil
}

func (a *Adapter) NewOrderedMap(capacity int) ifaces.OrderedMap {
	m := make(MapSlice, 0, capacity)

	return &m
}

func (a *Adapter) Redeem() {
	if a == nil {
		return
	}
	RedeemAdapter(a)
}

func (a *Adapter) Reset() {
	a.options = options{}
}

func newJWriter() *jwriter.Writer {
	return &jwriter.Writer{
		Flags: jwriter.NilMapAsEmpty | jwriter.NilSliceAsEmpty,
	}
}

func newJLexer() *jlexer.Lexer {
	return &jlexer.Lexer{}
}
