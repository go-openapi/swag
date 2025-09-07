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

import (
	"sync"

	"github.com/mailru/easyjson/buffer"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type adaptersPool struct {
	sync.Pool
}

func (p *adaptersPool) Borrow() *Adapter {
	ptr := p.Get()

	return ptr.(*Adapter)
}

func (p *adaptersPool) Redeem(a *Adapter) {
	p.Put(a)
}

type writersPool struct {
	sync.Pool
}

func (p *writersPool) Borrow() *jwriter.Writer {
	ptr := p.Get()

	w := ptr.(*jwriter.Writer)
	w.Error = nil
	w.NoEscapeHTML = false
	w.Flags = 0
	w.Buffer = buffer.Buffer{}

	return w
}

func (p *writersPool) Redeem(w *jwriter.Writer) {
	p.Put(w)
}

type lexersPool struct {
	sync.Pool
}

var emptyLexer = jlexer.Lexer{}

func (p *lexersPool) Borrow(data []byte) *jlexer.Lexer {
	ptr := p.Get()

	l := ptr.(*jlexer.Lexer)
	*l = emptyLexer
	l.Data = data

	return l
}

func (p *lexersPool) Redeem(l *jlexer.Lexer) {
	p.Put(l)
}

var (
	poolOfAdapters = &adaptersPool{
		Pool: sync.Pool{
			New: func() any {
				return NewAdapter()
			},
		},
	}

	poolOfWriters = &writersPool{
		Pool: sync.Pool{
			New: func() any {
				return newJWriter()
			},
		},
	}

	poolOfLexers = &lexersPool{
		Pool: sync.Pool{
			New: func() any {
				return newJLexer()
			},
		},
	}
)

// BorrowAdapter borrows an [Adapter] from the pool, recycling already allocated instances.
func BorrowAdapter() *Adapter {
	return poolOfAdapters.Borrow()
}

// RedeemAdapter redeems an [Adapter] to the pool, so it may be recycled.
func RedeemAdapter(a *Adapter) {
	poolOfAdapters.Redeem(a)
}

// BorrowWriter borrows a [jwriter.Writer] from the pool, recycling already allocated instances.
func BorrowWriter() *jwriter.Writer {
	return poolOfWriters.Borrow()
}

// RedeemWriter redeems a [jwriter.Writer] to the pool, so it may be recycled.
func RedeemWriter(w *jwriter.Writer) {
	poolOfWriters.Redeem(w)
}

// BorrowLexer borrows a [jlexer.Lexer] from the pool, recycling already allocated instances.
func BorrowLexer(data []byte) *jlexer.Lexer {
	return poolOfLexers.Borrow(data)
}

// RedeemLexer redeems a [jlexer.Lexer] to the pool, so it may be recycled.
func RedeemLexer(l *jlexer.Lexer) {
	poolOfLexers.Redeem(l)
}
