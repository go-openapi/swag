package shared

import (
	"io"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestBytesBorrowRedeemReuse(t *testing.T) {
	s, redeem := Bytes.BorrowWithRedeem()
	require.Zerof(t, s.Len(), "expected empty buffer, got len %d", s.Len())
	s.Append([]byte("hello")...)
	require.Equalf(t, "hello", string(s.Slice()), "unexpected contents: %q", string(s.Slice()))
	redeem()

	s2, redeem2 := Bytes.BorrowWithRedeem()
	require.Zerof(t, s2.Len(), "expected reset buffer on reuse, got len %d", s2.Len())
	redeem2()
}

func TestBufferBorrowRedeemReuse(t *testing.T) {
	b := BorrowBuffer()
	require.Zerof(t, b.Len(), "expected empty buffer, got len %d", b.Len())
	b.WriteString("payload")
	RedeemBuffer(b)

	b2 := BorrowBuffer()
	require.Zerof(t, b2.Len(), "expected reset buffer on reuse, got len %d", b2.Len())
	RedeemBuffer(b2)
}

func TestStringsBorrowRedeemReuse(t *testing.T) {
	s, redeem := Strings.BorrowWithRedeem()
	require.Zerof(t, s.Len(), "expected empty buffer, got len %d", s.Len())
	s.Append([]string{"hello"}...)
	require.Equalf(t, []string{"hello"}, s.Slice(), "unexpected contents: %v", s.Slice())

	redeem()

	s2, redeem2 := Strings.BorrowWithRedeem()
	require.Zerof(t, s2.Len(), "expected reset buffer on reuse, got len %d", s2.Len())
	redeem2()
}

func TestBufferRedeemNilAndOversizedAreSafe(t *testing.T) {
	require.NotPanics(t, func() { RedeemBuffer(nil) }, "RedeemBuffer(nil) must not panic")

	b := BorrowBuffer()
	b.Grow(maxSharedCapacity * 2)                                                                    // oversized
	require.NotPanics(t, func() { RedeemBuffer(b) }, "redeeming an oversized buffer must not panic") // dropped, not recycled

	b2 := BorrowBuffer()
	require.LessOrEqualf(t, b2.Cap(), maxSharedCapacity, "oversized buffer should not have been recycled, got cap %d", b2.Cap())
	RedeemBuffer(b2)
}

func TestReaderBorrowRedeemReuse(t *testing.T) {
	r := BorrowReader([]byte("first"))
	got, err := io.ReadAll(r)
	require.NoError(t, err, "read error")
	require.Equalf(t, "first", string(got), "unexpected read: %q", got)
	RedeemReader(r)

	r2 := BorrowReader([]byte("second"))
	got2, err := io.ReadAll(r2)
	require.NoError(t, err, "read error")
	require.Equalf(t, "second", string(got2), "reader not reinitialized on reuse: %q", got2)
	RedeemReader(r2)
}

func TestReaderRedeemClearsData(t *testing.T) {
	r := BorrowReader([]byte("data"))
	RedeemReader(r)

	// After RedeemReader the reader must not reference the old data (Len reports
	// 0).
	require.Zerof(t, r.Len(), "expected reader to release its data on Redeem, Len = %d", r.Len())
}

func TestReaderRedeemNilIsSafe(t *testing.T) {
	require.NotPanics(t, func() { RedeemReader(nil) }, "RedeemReader(nil) must not panic")
}
