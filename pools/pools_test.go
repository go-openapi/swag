// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package pools

import (
	"sync"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// resettable is a Resettable type that records how many times Reset was called
// and carries a reference we can observe to verify the pool does not pin it
// while idle.
type resettable struct {
	resets int
	ref    *int
	data   int
}

func (r *resettable) Reset() {
	r.resets++
	r.ref = nil
	r.data = 0
}

func TestPoolBorrowRedeem(t *testing.T) {
	p := New[resettable]()

	a := p.Borrow()
	a.data = 42
	x := 7
	a.ref = &x
	p.Redeem(a)

	// borrowing again should hand back a clean (reset) instance.
	b := p.Borrow()
	assert.Zerof(t, b.data, "expected data cleared on reuse, got %d", b.data)
	assert.Nilf(t, b.ref, "expected ref cleared on reuse, got %v", b.ref)
}

func TestPoolResetOnBorrowAndRedeem(t *testing.T) {
	p := New[resettable]()

	a := p.Borrow() // fresh: new(T), then reset on borrow
	require.Equalf(t, 1, a.resets, "expected 1 reset after fresh borrow, got %d", a.resets)
	p.Redeem(a) // reset on redeem
	require.Equalf(t, 2, a.resets, "expected 2 resets after redeem, got %d", a.resets)

	b := p.Borrow() // recycled: reset on borrow
	if b != a {
		t.Skip("pool did not return the same instance; reset-timing assertion not applicable")
	}
	require.Equalf(t, 3, b.resets, "expected 3 resets after re-borrow, got %d", b.resets)
}

func TestPoolRedeemNilIsSafe(t *testing.T) {
	p := New[resettable]()

	// Must not panic and must not poison the pool with a typed-nil.
	require.NotPanics(t, func() { p.Redeem(nil) })

	got := p.Borrow()
	require.NotNil(t, got, "pool handed back a nil pointer after Redeem(nil) poisoned it")
	got.data = 1 // would panic on a nil pointer
}

func TestRedeemableBorrowWithRedeem(t *testing.T) {
	p := NewRedeemable[resettable]()

	v, redeem := p.BorrowWithRedeem()
	require.NotNil(t, v, "expected non-nil instance")
	require.NotNil(t, redeem, "expected non-nil redeemer")
	v.data = 99
	x := 3
	v.ref = &x
	redeem()

	w, redeem2 := p.BorrowWithRedeem()
	assert.Zerof(t, w.data, "expected clean instance on reuse, got data=%d", w.data)
	assert.Nilf(t, w.ref, "expected clean instance on reuse, got ref=%v", w.ref)
	redeem2()
}

func TestRedeemableDoubleRedeemPanics(t *testing.T) {
	p := NewRedeemable[resettable]()

	_, redeem := p.BorrowWithRedeem()
	redeem() // first redeem: fine

	require.Panics(t, func() { redeem() }, "expected a panic on double redeem")
}

func TestRedeemableReborrowRearmsState(t *testing.T) {
	p := NewRedeemable[resettable]()

	// borrow/redeem several times: the state must be re-armed on each borrow so
	// redeem keeps working.
	require.NotPanics(t, func() {
		for i := 0; i < 5; i++ {
			_, redeem := p.BorrowWithRedeem()
			redeem()
		}
	})
}

func TestPoolSliceDoubleRedeemPanics(t *testing.T) {
	p := NewPoolSlice[int]()

	s, redeem := p.BorrowWithRedeem()
	s.Append(1, 2, 3)
	redeem()

	require.Panics(t, func() { redeem() }, "expected a panic on double redeem of a pooled slice")
}

func TestRedeemableZeroAllocRedeem(t *testing.T) {
	if debugBuild {
		t.Skip("the poolsdebug build allocates a per-borrow redeemer to track redemptions")
	}
	p := NewRedeemable[resettable]()

	// warm the pool
	v, redeem := p.BorrowWithRedeem()
	redeem()
	_ = v

	allocs := testing.AllocsPerRun(100, func() {
		x, r := p.BorrowWithRedeem()
		x.data++
		r()
	})
	require.Zerof(t, allocs, "expected 0 allocs on warm borrow/redeem, got %v", allocs)
}

func TestSliceBasic(t *testing.T) {
	p := NewPoolSlice[int]()

	s, redeem := p.BorrowWithRedeem()
	require.Zerof(t, s.Len(), "expected empty slice, got len %d", s.Len())
	s.Append(1, 2, 3)
	require.Equalf(t, 3, s.Len(), "expected len 3, got %d", s.Len())
	require.Equalf(t, []int{1, 2, 3}, s.Slice(), "unexpected slice contents: %v", s.Slice())
	redeem()

	// reuse: must be reset to length 0.
	s2, redeem2 := p.BorrowWithRedeem()
	require.Zerof(t, s2.Len(), "expected reset slice len 0, got %d", s2.Len())
	redeem2()
}

func TestSliceResetZeroesElementReferences(t *testing.T) {
	// White-box: verify Reset zeroes the whole backing array so pointer elements
	// are not retained.
	var s Slice[*int]
	a, b, c := 1, 2, 3
	s.Append(&a, &b, &c)
	full := s.inner[:cap(s.inner)]

	s.Reset()

	require.Zerof(t, s.Len(), "expected len 0 after reset, got %d", s.Len())
	for i, ptr := range full {
		require.Nilf(t, ptr, "element %d not cleared after reset: %v", i, ptr)
	}
}

func TestSliceWithLengthIsZeroedAndSized(t *testing.T) {
	p := NewPoolSlice[int](WithLength(4), WithMinimumCapacity(8))

	s, redeem := p.BorrowWithRedeem()
	require.Equalf(t, 4, s.Len(), "expected fixed length 4, got %d", s.Len())
	require.GreaterOrEqualf(t, s.Cap(), 8, "expected cap >= 8, got %d", s.Cap())
	for i, v := range s.Slice() {
		require.Zerof(t, v, "expected zeroed element at %d, got %d", i, v)
	}
	for i, v := range s.IndexedElems() {
		require.Zerof(t, v, "expected zeroed element at %d, got %d", i, v)
	}
	// dirty it, redeem, and confirm it comes back zeroed at the configured length.
	raw := s.Slice()
	for i := range raw {
		raw[i] = i + 1
	}
	redeem()

	s2, redeem2 := p.BorrowWithRedeem()
	require.Equalf(t, 4, s2.Len(), "expected fixed length 4 on reuse, got %d", s2.Len())
	for i, v := range s2.Slice() {
		require.Zerof(t, v, "expected zeroed element at %d on reuse, got %d", i, v)
	}
	redeem2()
}

func TestSliceConcatReusesCapacity(t *testing.T) {
	var s Slice[int]
	s.Grow(16)
	before := cap(s.inner)
	s.Concat([]int{1, 2, 3})
	s.Concat([]int{4, 5})
	require.Equalf(t, []int{1, 2, 3, 4, 5}, s.Slice(), "unexpected concat result: %v", s.Slice())
	require.Equalf(t, before, cap(s.inner), "concat reallocated backing array: before=%d after=%d", before, cap(s.inner))
}

// Reset is what preserves grown capacity through a redeem (it does not clip).
//
// We test that directly: asserting capacity survives a *pool* round-trip would
// be unsound, since sync.Pool is free to drop an idle object and hand back a
// fresh one (it routinely does under -race).
func TestSliceResetPreservesCapacity(t *testing.T) {
	var s Slice[int]
	s.Grow(1024)
	grown := s.Cap()
	require.GreaterOrEqualf(t, grown, 1024, "expected cap >= 1024 after grow, got %d", grown)
	s.Append(1, 2, 3)
	s.Reset()

	require.Zerof(t, s.Len(), "expected len 0 after reset, got %d", s.Len())
	require.Equalf(t, grown, s.Cap(), "Reset must preserve grown capacity: before=%d after=%d", grown, s.Cap())
}

// resetWithCapacity is the drop path a capped pool uses to discard an oversized
// backing array.
func TestSliceResetWithCapacityDropsBacking(t *testing.T) {
	var s Slice[int]
	s.Grow(4096)
	s.resetWithCapacity(64)

	require.Zerof(t, s.Len(), "expected len 0, got %d", s.Len())
	require.Equalf(t, 64, s.Cap(), "expected capacity dropped to 64, got %d", s.Cap())
}

func TestWithMaxCapacityShrinksOversized(t *testing.T) {
	p := NewPoolSlice[int](WithMinimumCapacity(8), WithMaxCapacity(64))

	s, redeem := p.BorrowWithRedeem()
	s.Grow(1024)
	require.GreaterOrEqualf(t, s.Cap(), 1024, "expected cap >= 1024 after grow, got %d", s.Cap())
	redeem() // cap > 64 → backing should be discarded and replaced

	s2, redeem2 := p.BorrowWithRedeem()
	require.LessOrEqualf(t, s2.Cap(), 64, "expected oversized backing to be dropped on redeem, got cap %d", s2.Cap())
	require.GreaterOrEqualf(t, s2.Cap(), 8, "expected replacement to honor minimum capacity 8, got %d", s2.Cap())
	redeem2()
}

func TestWithMaxCapacityHonorsLength(t *testing.T) {
	p := NewPoolSlice[int](WithLength(4), WithMaxCapacity(64))

	s, redeem := p.BorrowWithRedeem()
	s.Grow(1024)
	redeem() // dropped: replacement must still be a clean length-4 slice

	s2, redeem2 := p.BorrowWithRedeem()
	require.Equalf(t, 4, s2.Len(), "expected replacement length 4, got %d", s2.Len())
	for i, v := range s2.Slice() {
		require.Zerof(t, v, "expected zeroed replacement element at %d, got %d", i, v)
	}
	redeem2()
}

// TestNoLeaksOnCleanRun runs in both modes: in release AssertNoLeaks is a no-op
// returning true; under -tags poolsdebug it genuinely verifies the clean
// borrow/redeem left nothing outstanding.
func TestNoLeaksOnCleanRun(t *testing.T) {
	ResetTracking()

	p := NewRedeemable[resettable]()
	for i := 0; i < 10; i++ {
		_, redeem := p.BorrowWithRedeem()
		redeem()
	}

	ps := NewPoolSlice[int]()
	s, redeem := ps.BorrowWithRedeem()
	s.Append(1, 2, 3)
	redeem()

	require.True(t, AssertNoLeaks(t), "clean run should report no leaks")
}

func TestConcurrentBorrowRedeem(t *testing.T) {
	p := New[resettable]()
	ps := NewPoolSlice[int]()

	// the whole concurrent churn must complete without panicking (and, under
	// -race, without a data race).
	require.NotPanics(t, func() {
		var wg sync.WaitGroup
		for g := 0; g < 50; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 1000; i++ {
					v := p.Borrow()
					v.data = i
					p.Redeem(v)

					s, redeem := ps.BorrowWithRedeem()
					s.Append(i, i+1)
					redeem()
				}
			}()
		}
		wg.Wait()
	})
}
