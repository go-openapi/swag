// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package mangling

import (
	"fmt"
	"io"
	"testing"
)

func BenchmarkToXXXName(b *testing.B) {
	samples := []string{
		"sample text",
		"sample-text",
		"sample_text",
		"sampleText",
		"sample 2 Text",
		"findThingById",
		"日本語sample 2 Text",
		"日本語findThingById",
		"findTHINGSbyID",
	}
	m := NewNameMangler()

	b.Run("ToGoName", benchmarkFunc(m.ToGoName, samples))
	b.Run("ToVarName", benchmarkFunc(m.ToVarName, samples))
	b.Run("ToFileName", benchmarkFunc(m.ToFileName, samples))
	b.Run("ToCommandName", benchmarkFunc(m.ToCommandName, samples))
	b.Run("ToHumanNameLower", benchmarkFunc(m.ToHumanNameLower, samples))
	b.Run("ToHumanNameTitle", benchmarkFunc(m.ToHumanNameTitle, samples))
}

func benchmarkFunc(fn func(string) string, samples []string) func(*testing.B) {
	return func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var res string
		for i := 0; i < b.N; i++ {
			res = fn(samples[i%len(samples)])
		}

		fmt.Fprintln(io.Discard, res)
	}
}
