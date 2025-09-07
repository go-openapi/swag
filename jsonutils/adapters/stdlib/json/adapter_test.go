// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"regexp"
	"testing"

	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
	"github.com/stretchr/testify/require"
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
