// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//go:build poolsdebug

package pools

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

// These tests exercise the instrumented pool and only compile/run under -tags
// poolsdebug.

// capturePanic runs f and returns the value it panicked with (nil if it did not
// panic). It lets tests assert on a panic's message with require.Contains rather
// than hand-rolling recover().
func capturePanic(f func()) (recovered any) {
	defer func() { recovered = recover() }()
	f()
	return recovered
}

func TestDebugDetectsLeak(t *testing.T) {
	ResetTracking()

	p := NewRedeemable[resettable]()
	_, _ = p.BorrowWithRedeem() // borrowed, never redeemed → a leak

	var fake fakeTB
	require.False(t, AssertNoLeaks(&fake), "expected AssertNoLeaks to report the leak")
	require.NotEmpty(t, fake.errors, "expected an Errorf about leaked objects")
	require.NotEmpty(t, fake.logs, "expected a log naming the borrow site")
	require.Contains(t, fake.logs[0], "borrowed at", "expected a log naming the borrow site")
	// the recorded borrow site should point at THIS test file, validating the
	// caller() skip depth.
	require.Contains(t, fake.logs[0], "pools_debug_test.go", "expected borrow site in pools_debug_test.go")
}

func TestDebugCleanRunHasNoLeak(t *testing.T) {
	ResetTracking()

	p := NewRedeemable[resettable]()
	_, redeem := p.BorrowWithRedeem()
	redeem()

	var fake fakeTB
	require.Truef(t, AssertNoLeaks(&fake), "clean run should have no leaks, got errors=%v logs=%v", fake.errors, fake.logs)
}

func TestDebugDoubleRedeemRichPanic(t *testing.T) {
	ResetTracking()

	p := NewRedeemable[resettable]()
	_, redeem := p.BorrowWithRedeem()
	redeem()

	rec := capturePanic(func() { redeem() })
	require.NotNil(t, rec, "expected a panic on double redeem")
	require.Contains(t, rec, "double redeem", "expected a double-redeem panic, got %q", rec)
}

func TestDebugForeignRedeemPanics(t *testing.T) {
	ResetTracking()

	p := New[resettable]()
	foreign := &resettable{} // never borrowed from p

	rec := capturePanic(func() { p.Redeem(foreign) })
	require.NotNil(t, rec, "expected a panic on redeem of a foreign object")
	require.Contains(t, rec, "never handed out", "expected a never-handed-out panic, got %q", rec)
}

func TestDebugABADetected(t *testing.T) {
	ResetTracking()

	p := NewRedeemable[resettable]()

	innerA, redeemA := p.BorrowWithRedeem()
	redeemA() // A is valid, returned to the pool

	innerB, redeemB := p.BorrowWithRedeem()
	if innerB != innerA {
		t.Skip("pool returned a different slot; ABA scenario not reproduced this run")
	}
	defer redeemB() // keep B checked out so the stale redeemA hits the re-borrowed slot

	rec := capturePanic(redeemA) // stale: the slot was re-borrowed by B since this borrow
	require.NotNil(t, rec, "expected a panic on a stale (ABA) redeem")
	msg, _ := rec.(string)
	require.Truef(t,
		strings.Contains(msg, "stale borrow") || strings.Contains(msg, "ABA"),
		"expected a stale-borrow/ABA panic, got %q", msg,
	)
}
