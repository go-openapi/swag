// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
