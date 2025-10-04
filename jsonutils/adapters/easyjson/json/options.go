// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

// Option selects options for the easyjson adapter.
type Option func(o *options)

type options struct {
	writerOptions
	lexerOptions
}

type lexerOptions struct {
	useMultipleErrors bool
}

type writerOptions struct {
	nilMapAsEmpty   bool
	nilSliceAsEmpty bool
	noEscapeHTML    bool
}

func WithLexerUseMultipleErrors(enabled bool) Option {
	return func(o *options) {
		o.useMultipleErrors = enabled
	}
}

func WithWriterNilMapAsEmpty(enabled bool) Option {
	return func(o *options) {
		o.nilMapAsEmpty = enabled
	}
}

func WithWriterNilSliceAsEmpty(enabled bool) Option {
	return func(o *options) {
		o.nilSliceAsEmpty = enabled
	}
}

func WithWriterNoEscapeHTML(noescape bool) Option {
	return func(o *options) {
		o.noEscapeHTML = noescape
	}
}
