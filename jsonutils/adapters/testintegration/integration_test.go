// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testintegration

import (
	"testing"

	stdlib "github.com/go-openapi/swag/jsonutils/adapters/stdlib/json"

	"github.com/go-openapi/testify/v2/require"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	a := stdlib.BorrowAdapter()
	defer func() {
		stdlib.RedeemAdapter(a)
	}()

	const reasonableLength = 10
	constructor := func() *stdlib.MapSlice {
		m := a.NewOrderedMap(reasonableLength) // returns ifaces.OrderedMap
		stdm, ok := m.(*stdlib.MapSlice)

		require.True(t, ok)

		return stdm
	}

	t.Run("with stdlib OrderedMap implementation", runTestSuite(constructor, constructor, assertionTypeOrdered))
	t.Run("with stdlib unordered", runTestSuite[any, any](nil, nil, assertionTypeUnordered))

	t.Run("with easyjson ordered object", runTestSuite(newEasyOrderedObject, newEasyOrderedObject, assertionTypeOrdered,
		withAssertionsAfterRead(func(v any) func(*testing.T) {
			return func(t *testing.T) {
				value, ok := v.(EasyOrderedObject)
				require.True(t, ok)
				require.Len(t, value.UnmarshalEasyJSONCalls(), 1)
			}
		}),
		withAssertionsAfterWrite(func(v any) func(*testing.T) {
			return func(t *testing.T) {
				value, ok := v.(EasyOrderedObject)
				require.True(t, ok)
				require.Len(t, value.MarshalEasyJSONCalls(), 1)
			}
		}),
	))
	t.Run("with easyjson unordered", runTestSuite(newEasyObject, newEasyObject, assertionTypeUnordered))
}
