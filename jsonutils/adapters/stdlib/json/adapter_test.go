// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"regexp"
	"testing"

	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
	"github.com/go-openapi/testify/v2/require"
)

func TestAdapter(t *testing.T) {
	t.Parallel()

	const reasonableCapacity = 10
	a := BorrowAdapter()
	defer func() {
		RedeemAdapter(a)
	}()

	harness := fixtures.NewHarness(t)
	harness.Init()

	for name, test := range harness.AllTests(
		// in these test conditions we do not return nil when token is null, but an empty slice.
		fixtures.WithExcludePattern(regexp.MustCompile(`^with null value$`)),
	) {
		t.Run(name, func(t *testing.T) {
			t.Run("should Unmarshal JSON", func(t *testing.T) {
				value := a.NewOrderedMap(reasonableCapacity)

				if test.ExpectError() {
					require.Error(t, a.Unmarshal(test.JSONBytes(), value))

					return
				}

				require.NoError(t, a.Unmarshal(test.JSONBytes(), value))

				t.Run("should Marshal JSON with equivalent JSON", func(t *testing.T) {
					jazon, err := a.Marshal(value)
					require.NoError(t, err)

					require.JSONEq(t, test.JSONPayload, string(jazon))
				})
			})

			t.Run("should OrderedUnmarshal JSON", func(t *testing.T) {
				value := a.NewOrderedMap(reasonableCapacity)

				if test.ExpectError() {
					require.Error(t, a.OrderedUnmarshal(test.JSONBytes(), value))

					return
				}

				require.NoError(t, a.OrderedUnmarshal(test.JSONBytes(), value))

				t.Run("should OrderedMarshal JSON with identical JSON", func(t *testing.T) {
					jazon, err := a.OrderedMarshal(value)
					require.NoError(t, err)

					fixtures.JSONEqualOrdered(t, test.JSONPayload, string(jazon))
				})
			})
		})
	}
}
