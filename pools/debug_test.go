//go:build poolsdebug

// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package pools

import "fmt"

// fakeTB captures TB calls so leak assertions can be inspected without failing
// the host test.
type fakeTB struct {
	helpers int
	errors  []string
	logs    []string
}

func (f *fakeTB) Helper() { f.helpers++ }
func (f *fakeTB) Errorf(format string, args ...any) {

	f.errors = append(f.errors, fmt.Sprintf(format, args...))
}
func (f *fakeTB) Logf(format string, args ...any) {

	f.logs = append(f.logs, fmt.Sprintf(format, args...))
}
