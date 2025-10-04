// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"sync"

	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	"github.com/mailru/easyjson/buffer"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type adaptersPool struct {
	sync.Pool
}

func (p *adaptersPool) Borrow() *Adapter {
	return p.Get().(*Adapter)
}

func (p *adaptersPool) BorrowIface() ifaces.Adapter {
	return p.Get().(*Adapter)
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

func BorrowAdapterIface() ifaces.Adapter {
	return poolOfAdapters.BorrowIface()
}

// RedeemAdapter redeems an [Adapter] to the pool, so it may be recycled.
func RedeemAdapter(a *Adapter) {
	poolOfAdapters.Redeem(a)
}

func RedeemAdapterIface(a ifaces.Adapter) {
	concrete, ok := a.(*Adapter)
	if ok {
		poolOfAdapters.Redeem(concrete)
	}
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
