// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

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
