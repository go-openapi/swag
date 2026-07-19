// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	stdjson "encoding/json"
	"errors"

	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	"github.com/go-openapi/swag/typeutils"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
)

// ErrMaxNestingDepth is returned when a JSON document or in-memory structure nests
// deeper than the configured maximum (see [WithMaxNestingDepth]). It guards against
// stack-overflow crashes on adversarially deep input.
var ErrMaxNestingDepth = errors.New("maximum JSON nesting depth exceeded")

var _ ifaces.Adapter = &Adapter{}

type Adapter struct {
	options
}

// NewAdapter yields a JSON adapter for [easyjson].
func NewAdapter(opts ...Option) *Adapter {
	var o options

	return &Adapter{
		options: buildOptions(o, opts),
	}
}

func (a *Adapter) Marshal(value any) ([]byte, error) {
	marshaler, ok := value.(easyjson.Marshaler)
	if ok {
		w, redeem := BorrowWriter()
		defer redeem()

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
		l, redeem := BorrowLexer(data)
		defer redeem()

		l.UseMultipleErrors = a.useMultipleErrors
		unmarshaler.UnmarshalEasyJSON(l)

		return l.Error()
	}

	return stdjson.Unmarshal(data, value)
}

func (a *Adapter) OrderedMarshal(value ifaces.Ordered) ([]byte, error) {
	w, redeem := BorrowWriter()
	defer redeem()

	a.orderedMarshal(w, value, a.maxDepth())

	return w.BuildBytes() // this actually copies data, so its okay to redeem the writer
}

func (a *Adapter) OrderedUnmarshal(data []byte, value ifaces.SetOrdered) error {
	var m MapSlice
	if err := m.orderedUnmarshalJSON(data, a.maxDepth()); err != nil {
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

// orderedMarshal writes value to w, decreasing budget for every nested container to
// guard against stack overflow on deeply nested structures.
func (a *Adapter) orderedMarshal(w *jwriter.Writer, value ifaces.Ordered, budget int) {
	if typeutils.IsNil(value) {
		w.RawString("null")

		return
	}

	if budget <= 0 {
		w.Error = ErrMaxNestingDepth

		return
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
		case ifaces.Ordered:
			// ordered values (including this package's MapSlice) recurse through the
			// depth-guarded path rather than their own unbounded MarshalEasyJSON.
			a.orderedMarshal(w, val, budget-1)
		case easyjson.Marshaler:
			val.MarshalEasyJSON(w)
		default:
			w.Raw(stdjson.Marshal(v))
		}
	}

	w.RawByte('}')
}
