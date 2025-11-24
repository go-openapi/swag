// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package mangling

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

const (
	// parts used to build fixtures

	textTitle  = "Text"
	blankText  = " text"
	dashText   = "-text"
	uscoreText = "_text"

	sampleTitle  = "Sample"
	sampleString = "sample"
	sampleBlank  = "sample "
	sampleDash   = "sample-"
	sampleUscore = "sample_"
)

func TestManglerToGoName_Issue158(t *testing.T) {
	m := NewNameMangler()

	t.Run("should detect trailing pluralized initialisms", func(t *testing.T) {
		require.Equal(t, "LinkLocalIPs", m.ToGoName("LinkLocalIPs"))
		require.Equal(t, "NativeBaseURLs", m.ToGoName("nativeBaseURLs"))
		require.Equal(t, "SiteURLs", m.ToGoName("siteURLs"))
	})
}

func TestManglerToGoName_Issue159(t *testing.T) {
	m := NewNameMangler()
	// overlapping initialisms
	// TTLs
	//  TLS

	// require.Equal(t, "TTLs", m.ToGoName("TTLs"))
	// require.Equal(t, "TTLsS", m.ToGoName("TTLsS"))
	require.Equal(t, "TTLss", m.ToGoName("TTLss"))
}

func TestManglerToGoName(t *testing.T) {
	m := NewNameMangler()

	t.Run("with simple input", func(t *testing.T) {
		samples := []translationSample{
			// input, expected
			{"@Type", "AtType"},
			{"Sample@where", "SampleAtWhere"},
			{"Id", "ID"},
			{"SomethingTTLSeconds", "SomethingTTLSeconds"},
			{"sample text", "SampleText"},
			{"IPv6Address", "IPv6Address"}, // changed assertion: favor IPv6 over IPV6
			{"IPv4Address", "IPv4Address"}, // changed assertion: favor IPv4 over IPV4
			{"sample-text", "SampleText"},
			{"sample_text", "SampleText"},
			{"sampleText", "SampleText"},
			{"sample 2 Text", "Sample2Text"},
			{"findThingById", "FindThingByID"},
			{"日本語sample 2 Text", "X日本語sample2Text"},
			{"日本語findThingById", "X日本語findThingByID"},
			{"findTHINGSbyID", "FindTHINGSbyID"},
			{"x-isAnOptionalHeader0", "XIsAnOptionalHeader0"},
			{"get$ref", "GetDollarRef"},
			{"éget$ref", "ÉgetDollarRef"},
			{"日get$ref", "X日getDollarRef"},
			{"", ""},
			{"?", ""},
			{"!", "Bang"},
			{"", ""},
			{"Http Server", "HTTPServer"},
		}

		t.Run("ToGoName should convert names as expected", func(t *testing.T) {
			for _, sample := range samples {
				result := m.ToGoName(sample.str)
				assert.Equal(t, sample.out, result,
					"expected ToGoName(%q) == %q but got %q", sample.str, sample.out, result)
			}
		})
	})

	t.Run("with composed with initialism sample", func(t *testing.T) {
		for _, k := range m.Initialisms() {
			samples := []translationSample{
				// input, expected
				{sampleBlank + lower(k) + blankText, sampleTitle + k + textTitle},
				{sampleDash + lower(k) + dashText, sampleTitle + k + textTitle},
				{sampleUscore + lower(k) + uscoreText, sampleTitle + k + textTitle},
				{sampleString + titleize(k) + textTitle, sampleTitle + k + textTitle},
				{sampleBlank + lower(k), sampleTitle + k},
				{sampleDash + lower(k), sampleTitle + k},
				{sampleUscore + lower(k), sampleTitle + k},
				{sampleString + titleize(k), sampleTitle + k},
				{sampleBlank + titleize(k) + blankText, sampleTitle + k + textTitle},
				{sampleDash + titleize(k) + dashText, sampleTitle + k + textTitle},
				{sampleUscore + titleize(k) + uscoreText, sampleTitle + k + textTitle},
				// leading initialism preserving case in initialism, e.g. Ipv4_Address -> IPv4Address and no IPV4Address
				{titleize(k) + uscoreText, k + textTitle},
				// leading initialism preserving case in initialism, e.g. ipv4_Address -> IPv4Address and no IPV4Address
				{lower(k) + uscoreText, k + textTitle},
			}

			for _, sample := range samples {
				result := m.ToGoName(sample.str)
				assert.Equal(t, sample.out, result,
					"with initialism %q, expected ToGoName(%q) == %q but got %q", k, sample.str, sample.out, result)
			}
		}
	})

	t.Run("with prefix rule", func(t *testing.T) {
		samples := []translationSample{
			{"123_a", "Nr123a"},
			{"!123_a", "Bang123a"},
			{"+123_a", "Plus123a"},
			{"abc", "Abc"},
			{"éabc", "Éabc"},
			{":éabc", "Éabc"},
			{"get$ref", "GetDollarRef"},
			{"get!ref", "GetBangRef"},
			{"get&ref", "GetAndRef"},
			{"get|ref", "GetPipeRef"},
		}
		t.Run("with GoNamePrefixFunc", func(t *testing.T) {
			m := NewNameMangler(
				WithGoNamePrefixFunc(func(name string) string {
					// this is the pascalize func from go-swagger codegen
					arg := []rune(name)
					if len(arg) == 0 || arg[0] > '9' {
						return ""
					}
					if arg[0] == '+' {
						return "Plus"
					}
					if arg[0] == '-' {
						return "Minus"
					}

					return "Nr"
				}),
			)

			for _, sample := range samples {
				assert.Equal(t, sample.out, m.ToGoName(sample.str))
			}
		})

		t.Run("with GoNamePrefixFuncPtr", func(t *testing.T) {
			var fn PrefixFunc = func(name string) string {
				arg := []rune(name)
				if len(arg) == 0 || arg[0] > '9' {
					return ""
				}
				if arg[0] == '+' {
					return "Plus"
				}
				if arg[0] == '-' {
					return "Minus"
				}

				return "Nr"
			}

			m := NewNameMangler(
				WithGoNamePrefixFuncPtr(&fn),
			)

			for _, sample := range samples {
				assert.Equal(t, sample.out, m.ToGoName(sample.str))
			}
		})
	})

	t.Run("with Unicode edge cases", func(t *testing.T) {
		m := NewNameMangler()

		samples := []translationSample{
			{
				// single letter rune
				"ã",
				`Ã`,
			},
			{
				// single non letter rune (ascii)
				"3",
				`X3`,
			},
			{
				// multi non letter rune (ascii)
				"23",
				`X23`,
			},
			{
				// single non letter rune (devanagari digit)
				"१",
				`X१`,
			},
			{
				// single letter, no uppercase rune (devanagari letter)
				"आ",
				`Xआ`,
			},
			// TODO: non unicode char
		}

		for _, sample := range samples {
			assert.Equal(t, sample.out, m.ToGoName(sample.str))
		}
	})

	t.Run("with replace table", func(t *testing.T) {
		m := NewNameMangler(
			WithReplaceFunc(func(r rune) (string, bool) {
				switch r {
				case '$':
					return "Dollar ", true
				case '€':
					return "Euro ", true
				default:
					return "", false
				}
			}),
		)

		samples := []translationSample{
			{"a$ b", "ADollarb"},
			{"a€ b", "AEurob"},
		}
		for _, sample := range samples {
			assert.Equal(t, sample.out, m.ToGoName(sample.str))
		}
	})
}

func TestManglerToFileName(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)
	samples := []translationSample{
		{"SampleText", "sample_text"},
		{"FindThingByID", "find_thing_by_id"},
		{"FindThingByIDs", "find_thing_by_ids"},
		{"CAPWD.folwdBylc", "capwd_folwd_bylc"},
		{"CAPWDfolwdBylc", "cap_w_dfolwd_bylc"},
		{"CAP_WD_folwdBylc", "cap_wd_folwd_bylc"},
		{"TypeOAI_alias", "type_oai_alias"},
		{"Type_OAI_alias", "type_oai_alias"},
		{"Type_OAIAlias", "type_oai_alias"},
		{"ELB.HTTPLoadBalancer", "elb_http_load_balancer"},
		{"elbHTTPLoadBalancer", "elb_http_load_balancer"},
		{"ELBHTTPLoadBalancer", "elb_http_load_balancer"},
		{"get$Ref", "get_dollar_ref"},
	}
	for _, k := range m.Initialisms() {
		samples = append(samples,
			translationSample{sampleTitle + k + textTitle, sampleUscore + lower(k) + uscoreText},
		)
	}

	for _, sample := range samples {
		result := m.ToFileName(sample.str)
		assert.Equal(t, sample.out, m.ToFileName(sample.str),
			"ToFileName(%q) == %q but got %q", sample.str, sample.out, result)
	}
}

func TestManglerToCommandName(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)
	samples := []translationSample{
		{"SampleText", "sample-text"},
		{"FindThingByID", "find-thing-by-id"},
		{"elbHTTPLoadBalancer", "elb-http-load-balancer"},
		{"get$ref", "get-dollar-ref"},
		{"get!ref", "get-bang-ref"},
	}

	for _, k := range m.Initialisms() {
		samples = append(samples,
			translationSample{sampleTitle + k + textTitle, sampleDash + lower(k) + dashText},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, m.ToCommandName(sample.str))
	}
}

func TestManglerToHumanName(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)
	samples := []translationSample{
		{"Id", "Id"},
		{"IDs", "IDs"},
		{"SampleText", "sample text"},
		{"FindThingByID", "find thing by ID"},
		{"elbHTTPLoadBalancer", "elb HTTP load balancer"},
	}

	for _, k := range m.Initialisms() {
		samples = append(samples,
			translationSample{sampleTitle + k + textTitle, sampleBlank + k + blankText},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, m.ToHumanNameLower(sample.str))
	}
}

func TestManglerToJSONName(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)
	samples := []translationSample{
		{"SampleText", "sampleText"},
		{"FindThingByID", "findThingById"},
		{"elbHTTPLoadBalancer", "elbHttpLoadBalancer"},
		{"get$ref", "getDollarRef"},
		{"get!ref", "getBangRef"},
	}

	for _, k := range m.Initialisms() {
		samples = append(samples,
			translationSample{sampleTitle + k + textTitle, sampleString + titleize(k) + textTitle},
		)
	}

	for _, sample := range samples {
		assert.Equal(t, sample.out, m.ToJSONName(sample.str))
	}
}

func TestManglerCamelize(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)

	t.Run("with empty input", func(t *testing.T) {
		assert.Empty(t, m.Camelize(""))
	})

	t.Run("with single byte", func(t *testing.T) {
		assert.Equal(t, "A", m.Camelize("a"))
	})

	t.Run("with single multi-byte rune", func(t *testing.T) {
		assert.Equal(t, "Ã", m.Camelize("ã"))
	})

	samples := []translationSample{
		{"SampleText", "Sampletext"},
		{"FindThingByID", "Findthingbyid"},
		{"CAPWD.folwdBylc", "Capwd.folwdbylc"},
		{"CAPWDfolwdBylc", "Capwdfolwdbylc"},
		{"CAP_WD_folwdBylc", "Cap_wd_folwdbylc"},
		{"TypeOAI_alias", "Typeoai_alias"},
		{"Type_OAI_alias", "Type_oai_alias"},
		{"Type_OAIAlias", "Type_oaialias"},
		{"ELB.HTTPLoadBalancer", "Elb.httploadbalancer"},
		{"elbHTTPLoadBalancer", "Elbhttploadbalancer"},
		{"ELBHTTPLoadBalancer", "Elbhttploadbalancer"},
		{"12ab", "12ab"},
		{"get$Ref", "Get$ref"},
		{"get!Ref", "Get!ref"},
	}

	for _, sample := range samples {
		res := m.Camelize(sample.str)
		assert.Equalf(t, sample.out, res, "expected Camelize(%q)=%q, got %q", sample.str, sample.out, res)
	}
}

func TestManglerToHumanNameTitle(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)
	samples := []translationSample{
		{"SampleText", "Sample Text"},
		{"FindThingByID", "Find Thing By ID"},
		{"CAPWD.folwdBylc", "CAPWD Folwd Bylc"},
		{"CAPWDfolwdBylc", "CAP W Dfolwd Bylc"},
		{"CAP_WD_folwdBylc", "CAP WD Folwd Bylc"},
		{"TypeOAI_alias", "Type OAI Alias"},
		{"Type_OAI_alias", "Type OAI Alias"},
		{"Type_OAIAlias", "Type OAI Alias"},
		{"ELB.HTTPLoadBalancer", "ELB HTTP Load Balancer"},
		{"elbHTTPLoadBalancer", "elb HTTP Load Balancer"},
		{"ELBHTTPLoadBalancer", "ELB HTTP Load Balancer"},
		{"get$ref", "Get Dollar Ref"},
		{"get!ref", "Get Bang Ref"},
	}

	for _, sample := range samples {
		res := m.ToHumanNameTitle(sample.str)
		assert.Equalf(t, sample.out, res, "expected ToHumanNameTitle(%q)=%q, got %q", sample.str, sample.out, res)
	}
}

func TestManglerToVarName(t *testing.T) {
	m := NewNameMangler(
		WithAdditionalInitialisms("elb", "cap", "capwd", "wd"),
	)
	samples := []translationSample{
		{"SampleText", "sampleText"},
		{"FindThingByID", "findThingByID"},
		{"CAPWD.folwdBylc", "capwdFolwdBylc"},
		{"CAPWDfolwdBylc", "capWDfolwdBylc"},   // first part not detected as initialism (contentious point)
		{"CAP_WD_folwdBylc", "capWDFolwdBylc"}, // first part is initialism
		{"TypeOAI_alias", "typeOAIAlias"},
		{"Type_OAI_alias", "typeOAIAlias"},
		{"Type_OAIAlias", "typeOAIAlias"},
		{"ELB.HTTPLoadBalancer", "elbHTTPLoadBalancer"}, // first part is initialism
		{"elbHTTPLoadBalancer", "elbHTTPLoadBalancer"},
		{"ELBHTTPLoadBalancer", "elbHTTPLoadBalancer"},
		{"Id", "id"},
		{"HTTP", "http"}, // single initialism
		{"A", "a"},       // single byte
		{"a", "a"},       // single byte, unchanged
		{"get$ref", "getDollarRef"},
		{"get!ref", "getBangRef"},
		{"日get$ref", "x日getDollarRef"}, // prefix rule (no uppercase letter - Japanese rune)
	}

	for _, sample := range samples {
		res := m.ToVarName(sample.str)
		assert.Equalf(t, sample.out, res, "expected ToVarName(%q)=%q, got %q", sample.str, sample.out, res)
	}

	t.Run("with Unicode edge cases", func(t *testing.T) {
		m := NewNameMangler()

		samples := []translationSample{
			{
				// single letter rune
				`Ã`,
				"ã",
			},
			{
				// single non letter rune (ascii)
				"3",
				`x3`,
			},
			{
				// multi non letter rune (ascii)
				"23",
				`x23`,
			},
			{
				// single non letter rune (devanagari digit)
				"१",
				`x१`,
			},
			{
				// single letter, no uppercase rune (devanagari letter)
				"आ",
				`xआ`,
			},
			// TODO: non unicode char
		}

		for _, sample := range samples {
			assert.Equal(t, sample.out, m.ToVarName(sample.str))
		}
	})
}

func TestManglerInitialisms(t *testing.T) {
	t.Run("with AddInitialisms", func(t *testing.T) {
		m := NewNameMangler()
		m.AddInitialisms("ELB", "OLTP")

		assert.Equal(t, "ELBEndpoint", m.ToGoName("elb_endpoint"))
		assert.Equal(t, "HTTPEndpoint", m.ToGoName("http_endpoint"))
		assert.Equal(t, "OLTPEndpoint", m.ToGoName("oltp endpoint"))
	})

	t.Run("with Initialisms", func(t *testing.T) {
		m := NewNameMangler(
			WithInitialisms("ELB", "OLTP"),
		)

		assert.Equal(t, "ELBEndpoint", m.ToGoName("elb_endpoint"))
		assert.Equal(t, "HttpEndpoint", m.ToGoName("http_endpoint")) // no recognized as initialisms (default override)
		assert.Equal(t, "OLTPEndpoint", m.ToGoName("oltp endpoint"))
	})
}
