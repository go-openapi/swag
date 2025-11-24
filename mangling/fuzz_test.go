// SPDX-FileCopyrightText: Copyright (c) 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package mangling

import (
	"iter"
	"slices"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func FuzzToGoName(f *testing.F) {
	// initial seed
	cumulated := make([]string, 0, 100)
	for generator := range generators() {
		f.Add(generator)

		cumulated = append(cumulated, generator)
		f.Add(strings.Join(cumulated, ""))
	}
	mangler := NewNameMangler()

	f.Fuzz(func(t *testing.T, input string) {
		require.NotPanics(t, func() {
			_ = mangler.ToGoName(input)
		})
	})
}

func generators() iter.Seq[string] {
	return slices.Values([]string{
		"                    ",
		"!",
		"!123_a",
		"",
		"+123_a",
		"-|x>",
		"123_a",
		":éabc",
		"?",
		"@Type",
		"AbC",
		"Abc",
		"AtType",
		"Bang",
		"Bang123a",
		"FindTHINGSbyID",
		"FindThingByID",
		"GetAndRef",
		"GetBangRef",
		"GetDollarRef",
		"GetDollarRef",
		"GetPipeRef",
		"HTTPServer",
		"Http Server",
		"ID",
		"IPv4Address",
		"IPv4Address",
		"IPv6Address",
		"IPv6Address",
		"Id",
		"LinkLocalIPs",
		"Nr123a",
		"Plus123a",
		"Sample2Text",
		"Sample@where",
		"SampleAtWhere",
		"SampleText",
		"SampleText",
		"SampleText",
		"SampleText",
		"SomethingTTLSeconds",
		"SomethingTTLSeconds",
		"XIsAnOptionalHeader0",
		"X日getDollarRef",
		"X日本語findThingByID",
		"X日本語sample2Text",
		"a b c",
		"abc",
		"findTHINGSbyID",
		"findThingById",
		"get!ref",
		"get$ref",
		"get$ref",
		"get&ref",
		"get|ref",
		"nativeBaseURLs",
		"sample 2 Text",
		"sample text",
		"sample-text",
		"sampleText",
		"sample_text",
		"siteURLs",
		"x-isAnOptionalHeader0",
		"Éabc",
		"Éabc",
		"ÉgetDollarRef",
		"éabc",
		"éget$ref",
		"日get$ref",
		"日本語findThingById",
		"日本語sample 2 Text",
	})
}
