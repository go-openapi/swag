// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

// defaultMaxNestingDepth is the default maximum number of nested JSON containers
// ('{' or '[') that the ordered-JSON marshaler and unmarshaler will process before
// returning an error.
//
// It mirrors the limit enforced by the standard library's [encoding/json] decoder
// (see encoding/json's internal maxNestingDepth), which this adapter would otherwise
// not benefit from since it drives [encoding/json.Decoder.Token] directly.
const defaultMaxNestingDepth = 10000

// Option selects options for the stdlib adapter.
type Option func(o *options)

type options struct {
	maxNestingDepth int
}

// maxDepth returns the configured maximum nesting depth, or the default when unset.
func (o options) maxDepth() int {
	if o.maxNestingDepth <= 0 {
		return defaultMaxNestingDepth
	}

	return o.maxNestingDepth
}

// WithMaxNestingDepth sets the maximum number of nested JSON containers accepted
// when marshaling or unmarshaling ordered JSON.
//
// A value <= 0 selects the default ([defaultMaxNestingDepth]).
//
// This guards against stack-overflow crashes on deeply nested (possibly adversarial)
// JSON documents or in-memory structures.
func WithMaxNestingDepth(depth int) Option {
	return func(o *options) {
		o.maxNestingDepth = depth
	}
}
