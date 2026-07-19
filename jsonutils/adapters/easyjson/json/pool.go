// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	"github.com/go-openapi/swag/pools"
	"github.com/mailru/easyjson/buffer"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

var (
	emptyLexer     = jlexer.Lexer{}
	poolOfAdapters = pools.New[Adapter]()
	poolOfWriters  = pools.NewRedeemable[jwriter.Writer]()
	poolOfLexers   = pools.NewRedeemable[jlexer.Lexer]()
)

// BorrowAdapter borrows an [Adapter] from the pool, recycling already allocated instances.
func BorrowAdapter() *Adapter {
	return poolOfAdapters.Borrow()
}

func BorrowAdapterIface() ifaces.Adapter {
	return poolOfAdapters.Borrow()
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
func BorrowWriter() (*jwriter.Writer, func()) {
	w, redeem := poolOfWriters.BorrowWithRedeem()
	w.Error = nil
	w.NoEscapeHTML = false
	// w.Flags = jwriter.NilMapAsEmpty | jwriter.NilSliceAsEmpty
	w.Flags = 0
	w.Buffer = buffer.Buffer{}

	return w, redeem
}

// BorrowLexer borrows a [jlexer.Lexer] from the pool, recycling already allocated instances.
func BorrowLexer(data []byte) (*jlexer.Lexer, func()) {
	l, redeem := poolOfLexers.BorrowWithRedeem()
	*l = emptyLexer
	l.Data = data

	return l, redeem
}
