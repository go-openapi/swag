// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

// defaultMaxNestingDepth is the default maximum number of nested JSON containers
// ('{' or '[') that the ordered-JSON marshaler and unmarshaler will process before
// returning an error.
//
// It mirrors the limit enforced by the standard library's [encoding/json] decoder,
// guarding against stack-overflow crashes on deeply nested (possibly adversarial) input.
const defaultMaxNestingDepth = 10000

// Option selects options for the easyjson adapter.
type Option func(o options) options

type options struct {
	writerOptions
	lexerOptions
}

type lexerOptions struct {
	useMultipleErrors bool
	maxNestingDepth   int
}

func buildOptions(o options, opts []Option) options {
	for _, apply := range opts {
		o = apply(o)
	}

	return o
}

// maxDepth returns the configured maximum nesting depth, or the default when unset.
func (o options) maxDepth() int {
	if o.maxNestingDepth <= 0 {
		return defaultMaxNestingDepth
	}

	return o.maxNestingDepth
}

type writerOptions struct {
	nilMapAsEmpty   bool
	nilSliceAsEmpty bool
	noEscapeHTML    bool
}

func WithLexerUseMultipleErrors(enabled bool) Option {
	return func(o options) options {
		o.useMultipleErrors = enabled

		return o
	}
}

// WithMaxNestingDepth sets the maximum number of nested JSON containers accepted
// when marshaling or unmarshaling ordered JSON.
//
// A value <= 0 selects the default (10,000).
//
// This guards against stack-overflow crashes on deeply nested (possibly adversarial)
// JSON documents or in-memory structures.
func WithMaxNestingDepth(depth int) Option {
	return func(o options) options {
		o.maxNestingDepth = depth

		return o
	}
}

func WithWriterNilMapAsEmpty(enabled bool) Option {
	return func(o options) options {
		o.nilMapAsEmpty = enabled

		return o
	}
}

func WithWriterNilSliceAsEmpty(enabled bool) Option {
	return func(o options) options {
		o.nilSliceAsEmpty = enabled

		return o
	}
}

func WithWriterNoEscapeHTML(noescape bool) Option {
	return func(o options) options {
		o.noEscapeHTML = noescape

		return o
	}
}
