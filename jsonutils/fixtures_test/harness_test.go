// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fixtures

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestHarness(t *testing.T) {
	h := NewHarness(t)

	h.Init()

	f1, ok := h.Get("with JSON object")
	require.True(t, ok)
	require.NotNil(t, f1)

	f2 := h.ShouldGet("with JSON object inside")
	require.NotNil(t, f2)

	count := 0
	for name, test := range h.AllTests() {
		count++
		require.NotEmpty(t, name)
		require.NotEmptyf(t, test.JSONPayload, "JSON payload empty for %s", name)
		require.NotEmptyf(t, test.YAMLPayload, "YAML payload empty for %s", name)
		require.NotEmptyf(t, test.JSONBytes(), "JSON bytes payload empty for %s", name)
		require.Equal(t, test.Error, test.ExpectError())
		if !test.ExpectError() {
			require.NotEmptyf(t, test.YAMLBytes(), "YAML bytes payload empty for %s", name)
		}
	}
	require.NotZero(t, count)

	countWithoutError := 0
	for name, test := range h.AllTests(WithoutError(true)) {
		require.Falsef(t, test.ExpectError(), "test %s did not expect an error", name)
		countWithoutError++
	}
	require.NotZero(t, countWithoutError)

	countWithError := 0
	for name, test := range h.AllTests(WithError(true)) {
		require.Truef(t, test.ExpectError(), "test %s expected an error", name)
		countWithError++
	}
	require.NotZero(t, countWithError)
	require.Equal(t, count, countWithoutError+countWithError)
}

func TestLoadFixture(t *testing.T) {
	require.NotPanics(t, func() {
		MustLoadFixture(EmbeddedFixtures, "ordered_fixtures.yaml")
	})
}

func TestJSONEqOrdered(t *testing.T) {
	JSONEqualOrdered(t,
		`{"a": 1.0, "b": "x"   , "c": 1e-2}`,
		`{"a"  : 1,"b":"x","c": 0.01}`,
	)
}

func TestYAMLEqOrdered(t *testing.T) {
	y1 := `
---
a:
  b:
    c: [1,2,3]
b: true
`

	y2 := `
b: true
a:
  b:
    c:
      - 1
      - 2
      - 3
`
	YAMLEqualOrdered(t,
		y1, y2,
	)
}
