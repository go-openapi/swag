// Package shared provides ready-made, process-wide pools for objects that are
// commonly recycled across potentially many packages: byte slices,
// [bytes.Buffer] and [bytes.Reader].
//
// Sharing a single warm pool per object type across the whole program gives a
// better hit rate than many small, private pools.
//
// The shared pools are capacity-bounded: an occasional large object is dropped
// on return rather than circulating forever and bloating process memory (see
// [pools.WithMaxCapacity]).
//
// Callers that need different sizing should build their own pool with the
// parent [pools] package instead of using these globals.
package shared

import (
	"bytes"

	"github.com/go-openapi/swag/pools"
)

// maxSharedCapacity bounds the size of objects recycled by the shared pools.
//
// Objects grown past it are dropped on return so a single large request does
// not inflate the shared pool's steady-state footprint.
// 64 KiB comfortably covers typical scratch buffers while keeping the bound
// modest.
const maxSharedCapacity = 1 << 16

// Bytes is the process-wide pool of []byte scratch buffers.
//
// Borrow with [pools.PoolSlice.BorrowWithRedeem] and grow through the returned
// wrapper's methods so growth is tracked and recycled:
//
//	s, redeem := shared.Bytes.BorrowWithRedeem()
//	defer redeem()
//	s.Append(data...)
//	use(s.Slice())
//
// Buffers grown beyond 64 KiB are dropped on redeem instead of being recycled.
var Bytes = pools.NewPoolSlice[byte](pools.WithMaxCapacity(maxSharedCapacity))

var (
	buffers       = pools.New[bytes.Buffer]()
	redeemBuffers = pools.NewRedeemable[bytes.Buffer]()
)

// BorrowBuffer borrows a reset [bytes.Buffer] from the shared pool.
func BorrowBuffer() *bytes.Buffer {
	return buffers.Borrow()
}

// BorrowBufferWithRedeem borrows a reset [bytes.Buffer] from the shared pool,
// with its redeem closure.
func BorrowBufferWithRedeem() (*bytes.Buffer, func()) {
	return redeemBuffers.BorrowWithRedeem()
}

// RedeemBuffer returns a [bytes.Buffer] to the shared pool.
//
// A nil buffer is ignored, and a buffer that has grown beyond 64 KiB is dropped
// (left for the GC) rather than bloating the shared pool.
//
// The buffer is reset before being recycled; do not use it after calling
// RedeemBuffer.
func RedeemBuffer(b *bytes.Buffer) {
	if b == nil || b.Cap() > maxSharedCapacity {
		return
	}
	buffers.Redeem(b)
}

var readers = pools.New[bytes.Reader]()

// BorrowReader borrows a [bytes.Reader] from the shared pool, positioned to
// read from b.
//
// [bytes.Reader] is not auto-resettable (its Reset takes an argument), so it is
// reinitialized here explicitly.
//
// Pair every call to [BorrowReader] with a [RedeemReader].
func BorrowReader(b []byte) *bytes.Reader {
	r := readers.Borrow()
	r.Reset(b)

	return r
}

// RedeemReader returns a [bytes.Reader] to the shared pool, first clearing its
// reference to the underlying data so the idle pooled reader does not keep that
// data alive.
//
// A nil reader is ignored; do not use the reader after calling RedeemReader.
func RedeemReader(r *bytes.Reader) {
	if r == nil {
		return
	}
	r.Reset(nil)
	readers.Redeem(r)
}

// maxStringSliceCapacity bounds the size of objects recycled by the shared pools.
//
// Objects grown past it are dropped on return so a single large request does
// not inflate the shared pool's steady-state footprint.
// 1024 slots comfortably covers typical string slices while keeping the bound
// modest.
const maxStringSliceCapacity = 1024

// Strings is the process-wide pool of []string scratch buffers.
//
// Borrow with [String.BorrowWithRedeem] and grow through the returned
// wrapper's methods so growth is tracked and recycled:
//
//	s, redeem := shared.Strings.BorrowWithRedeem()
//	defer redeem()
//	s.Append(data...)
//	use(s.Slice())
//
// Buffers grown beyond 1024 slices are dropped on redeem instead of being recycled.
var Strings = pools.NewPoolSlice[string](pools.WithMaxCapacity(maxStringSliceCapacity))
