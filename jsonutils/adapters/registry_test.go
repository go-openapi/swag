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

package adapters

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	"github.com/go-openapi/swag/jsonutils/adapters/ifaces/mocks"
	stdlib "github.com/go-openapi/swag/jsonutils/adapters/stdlib/json"
	"github.com/stretchr/testify/require"
)

func TestRegistryUnmarshal(t *testing.T) {
	t.Parallel()
	reg := NewRegistrar()

	t.Run("should handle new registration for all capabilities", func(t *testing.T) {
		register1(reg)
		require.Len(t, reg.marshalerRegistry, 2)
		require.Len(t, reg.unmarshalerRegistry, 2)
		require.Len(t, reg.orderedMarshalerRegistry, 2)
		require.Len(t, reg.orderedUnmarshalerRegistry, 2)

		t.Run("should serve adapter for capability Unmarshal", testUnmarshal[any, *MockAdapter1](reg))

		t.Run("should retrieve route from cache when calling Unmarshal", func(t *testing.T) {
			var value any
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityUnmarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			jazon := []byte("null")
			err := adapter.Unmarshal(jazon, &value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			mockAdapter, ok := adapter.(*MockAdapter1)
			require.True(t, ok)
			calls := mockAdapter.UnmarshalCalls()
			require.Len(t, calls, 1)
		})

		t.Run("should serve adapter for capability OrderedUnmarshalJSON", func(t *testing.T) {
			value := newMockOrdered()
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityOrderedUnmarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter *MockAdapter1
			require.IsType(t, expectedAdapter, adapter)

			jazon := []byte("null")
			err := adapter.OrderedUnmarshal(jazon, value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's OrderedUnmarshal method", func(t *testing.T) {
				mockAdapter, ok := adapter.(*MockAdapter1)
				require.True(t, ok)
				calls := mockAdapter.OrderedUnmarshalCalls()
				require.Len(t, calls, 1)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.orderedUnmarshalerCache, 1)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.orderedUnmarshalerCache, key)
				entry := reg.orderedUnmarshalerCache[key]
				require.Equal(t, "github.com/go-openapi/swag/jsonutils/adapters.MockAdapter1", entry.Who)
				require.True(t, entry.What.Has(ifaces.CapabilityOrderedUnmarshalJSON))
			})
		})

		t.Run("should retrieve route from cache when calling OrderedUnmarshal", func(t *testing.T) {
			value := newMockOrdered()
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityOrderedUnmarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			jazon := []byte("null")
			err := adapter.OrderedUnmarshal(jazon, value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			mockAdapter, ok := adapter.(*MockAdapter1)
			require.True(t, ok)
			calls := mockAdapter.OrderedUnmarshalCalls()
			require.Len(t, calls, 1)
		})
	})
}

func TestRegistryMarshal(t *testing.T) {
	t.Parallel()
	reg := NewRegistrar()

	t.Run("should handle new registration for all capabilities", func(t *testing.T) {
		register1(reg)
		require.Len(t, reg.marshalerRegistry, 2)
		require.Len(t, reg.unmarshalerRegistry, 2)
		require.Len(t, reg.orderedMarshalerRegistry, 2)
		require.Len(t, reg.orderedUnmarshalerRegistry, 2)

		t.Run("should serve adapter for capability MarshalJSON", testMarshal[any, *MockAdapter1](reg))

		t.Run("should retrieve route from cache when calling Marshal", func(t *testing.T) {
			var value any
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityMarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()

			jazon, err := adapter.Marshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			mockAdapter, ok := adapter.(*MockAdapter1)
			require.True(t, ok)
			calls := mockAdapter.MarshalCalls()
			require.Len(t, calls, 1)
		})

		t.Run("should serve adapter for capability OrderedMarshalJSON", func(t *testing.T) {
			value := newMockOrdered()
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityOrderedMarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter *MockAdapter1
			require.IsType(t, expectedAdapter, adapter)

			jazon, err := adapter.OrderedMarshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's OrderedMarshal method", func(t *testing.T) {
				mockAdapter, ok := adapter.(*MockAdapter1)
				require.True(t, ok)
				calls := mockAdapter.OrderedMarshalCalls()
				require.Len(t, calls, 1)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.orderedMarshalerCache, 1)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.orderedMarshalerCache, key)
				entry := reg.orderedMarshalerCache[key]
				require.Equal(t, "github.com/go-openapi/swag/jsonutils/adapters.MockAdapter1", entry.Who)
				require.True(t, entry.What.Has(ifaces.CapabilityOrderedMarshalJSON))
			})
		})

		t.Run("should retrieve route from cache when calling OrderedMarshal", func(t *testing.T) {
			value := newMockOrdered()
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityOrderedMarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			jazon, err := adapter.OrderedMarshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			mockAdapter, ok := adapter.(*MockAdapter1)
			require.True(t, ok)
			calls := mockAdapter.OrderedMarshalCalls()
			require.Len(t, calls, 1)
		})

		t.Run("should panic on unsupported capability", func(t *testing.T) {
			require.Panics(t, func() {
				var value any
				_, _ = reg.AdapterFor(ifaces.Capability(99), value)
			})
		})
	})

	t.Run("should register new adapter with limited capabilities", func(t *testing.T) {
		register2(reg)
		require.Len(t, reg.marshalerRegistry, 3)
		require.Len(t, reg.unmarshalerRegistry, 3)
		require.Len(t, reg.orderedMarshalerRegistry, 3)
		require.Len(t, reg.orderedUnmarshalerRegistry, 3)

		t.Run("should serve new adapter for capability MarshalJSON when type is supported", func(t *testing.T) {
			var value supportedType
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityMarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter *MockAdapter2
			require.IsType(t, expectedAdapter, adapter)

			jazon, err := adapter.Marshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's MarshalJSON method", func(t *testing.T) {
				mockAdapter, ok := adapter.(*MockAdapter2)
				require.True(t, ok)
				calls := mockAdapter.MarshalCalls()
				require.Len(t, calls, 1)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.marshalerCache, 2)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.marshalerCache, key)
				entry := reg.marshalerCache[key]
				require.Equal(t, "github.com/go-openapi/swag/jsonutils/adapters.MockAdapter2", entry.Who)
				require.True(t, entry.What.Has(ifaces.CapabilityMarshalJSON))
			})
		})

		t.Run("should serve previous adapter for capability MarshalJSON when type is NOT supported", func(t *testing.T) {
			var value struct{}
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityMarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter *MockAdapter1
			require.IsType(t, expectedAdapter, adapter)

			jazon, err := adapter.Marshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's MarshalJSON method", func(t *testing.T) {
				mockAdapter, ok := adapter.(*MockAdapter1)
				require.True(t, ok)
				calls := mockAdapter.MarshalCalls()
				require.Len(t, calls, 1)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.marshalerCache, 3)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.marshalerCache, key)
				entry := reg.marshalerCache[key]
				require.Equal(t, "github.com/go-openapi/swag/jsonutils/adapters.MockAdapter1", entry.Who)
				require.True(t, entry.What.Has(ifaces.CapabilityMarshalJSON))
			})
		})

		t.Run("should serve previous adapter for capability Unmarshal", func(t *testing.T) {
			var value supportedType
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityUnmarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter *MockAdapter1
			require.IsType(t, expectedAdapter, adapter)

			jazon := []byte("null")
			err := adapter.Unmarshal(jazon, &value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's Unmarshal method", func(t *testing.T) {
				mockAdapter, ok := adapter.(*MockAdapter1)
				require.True(t, ok)
				calls := mockAdapter.UnmarshalCalls()
				require.Len(t, calls, 1)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.unmarshalerCache, 1)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.unmarshalerCache, key)
				entry := reg.unmarshalerCache[key]
				require.Equal(t, "github.com/go-openapi/swag/jsonutils/adapters.MockAdapter1", entry.Who)
				require.True(t, entry.What.Has(ifaces.CapabilityUnmarshalJSON))
			})
		})
	})
}

func TestRegistryOrderedMap(t *testing.T) {
	t.Parallel()
	reg := NewRegistrar()

	t.Run("should handle new registration for all capabilities", func(t *testing.T) {
		register1(reg)
		require.Len(t, reg.marshalerRegistry, 2)
		require.Len(t, reg.unmarshalerRegistry, 2)
		require.Len(t, reg.orderedMarshalerRegistry, 2)
		require.Len(t, reg.orderedUnmarshalerRegistry, 2)

		t.Run("should serve adapter for capability OrderedMap", func(t *testing.T) {
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityOrderedMap, nil)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter *MockAdapter1
			require.IsType(t, expectedAdapter, adapter)

			orderedMap := adapter.NewOrderedMap(1)
			expectedMap := newMockOrdered()
			require.NotNil(t, orderedMap)
			require.IsType(t, expectedMap, orderedMap)

			t.Run("should have called the adapter's NewOrderedMap method", func(t *testing.T) {
				mockAdapter, ok := adapter.(*MockAdapter1)
				require.True(t, ok)
				calls := mockAdapter.NewOrderedMapCalls()
				require.Len(t, calls, 1)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.orderedMapCache, 1)
				key := reflect.TypeOf(nil)
				require.Contains(t, reg.orderedMapCache, key)
				entry := reg.orderedMapCache[key]
				require.Equal(t, "github.com/go-openapi/swag/jsonutils/adapters.MockAdapter1", entry.Who)
				require.True(t, entry.What.Has(ifaces.CapabilityOrderedMap))
			})
		})
	})
}

func TestEmptyRegistry(t *testing.T) {
	t.Parallel()

	reg := NewRegistrar()
	reg.marshalerRegistry = reg.marshalerRegistry[:0]

	t.Run("should not find an adapter for capability MarshalJSON", func(t *testing.T) {
		var value any
		adapter, redeem := reg.AdapterFor(ifaces.CapabilityMarshalJSON, value)
		require.Nil(t, adapter)
		require.NotNil(t, redeem)
		defer redeem()
	})
}

func TestGlobalRegistry(t *testing.T) {
	t.Parallel()

	t.Run("with default global registry", func(t *testing.T) {
		t.Run("should resolve to the stdlib adapter for MarshalJSON", func(t *testing.T) {
			var value any
			adp, redeem := MarshalAdapterFor(value)
			require.NotNil(t, adp)
			require.NotNil(t, redeem)
			defer redeem()

			_, isStdLib := adp.(*stdlib.Adapter)
			require.True(t, isStdLib)
		})

		t.Run("should resolve to the stdlib adapter for UnmarshalJSON", func(t *testing.T) {
			var value any
			adp, redeem := UnmarshalAdapterFor(value)
			require.NotNil(t, adp)
			require.NotNil(t, redeem)
			defer redeem()

			_, isStdLib := adp.(*stdlib.Adapter)
			require.True(t, isStdLib)
		})

		t.Run("should resolve to the stdlib adapter for OrderedMarshalJSON", func(t *testing.T) {
			value := newMockOrdered()
			adp, redeem := OrderedMarshalAdapterFor(value)
			require.NotNil(t, adp)
			require.NotNil(t, redeem)
			defer redeem()

			_, isStdLib := adp.(*stdlib.Adapter)
			require.True(t, isStdLib)
		})

		t.Run("should resolve to the stdlib adapter for OrderedUnmarshalJSON", func(t *testing.T) {
			value := newMockOrdered()
			adp, redeem := OrderedUnmarshalAdapterFor(value)
			require.NotNil(t, adp)
			require.NotNil(t, redeem)
			defer redeem()

			_, isStdLib := adp.(*stdlib.Adapter)
			require.True(t, isStdLib)
		})

		t.Run("should resolve to the stdlib adapter for OrderedMap", func(t *testing.T) {
			var expectedMap *stdlib.MapSlice
			orderedMap := NewOrderedMap(1)
			require.NotNil(t, orderedMap)
			require.IsType(t, expectedMap, orderedMap)

			_, isStdLib := orderedMap.(*stdlib.MapSlice)
			require.True(t, isStdLib)
		})
	})
}

func testUnmarshal[ValueType any, AdapterType ifaces.UnmarshalAdapter](reg *Registrar) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("should serve adapter for capability Unmarshal", func(t *testing.T) {
			var value ValueType
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityUnmarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter AdapterType
			require.IsType(t, expectedAdapter, adapter)

			jazon := []byte("null")
			err := adapter.Unmarshal(jazon, &value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's Unmarshal method", func(t *testing.T) {
				_, ok := adapter.(AdapterType)
				require.True(t, ok)
				auditable, ok := adapter.(interface{ UnmarshalCallsLen() int })
				require.True(t, ok)

				calls := auditable.UnmarshalCallsLen()
				require.Equal(t, 1, calls)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.unmarshalerCache, 1)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.unmarshalerCache, key)
				entry := reg.unmarshalerCache[key]
				mockAdapter, ok := adapter.(AdapterType)
				require.True(t, ok)
				require.Equal(t,
					fmt.Sprintf("github.com/go-openapi/swag/jsonutils/%s", reflect.Indirect(reflect.ValueOf(mockAdapter)).Type()),
					entry.Who,
				)
				require.True(t, entry.What.Has(ifaces.CapabilityUnmarshalJSON))
			})
		})
	}
}

func testMarshal[ValueType any, AdapterType ifaces.MarshalAdapter](reg *Registrar) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("should serve adapter for capability MarshalJSON", func(t *testing.T) {
			var value ValueType
			adapter, redeem := reg.AdapterFor(ifaces.CapabilityMarshalJSON, value)
			require.NotNil(t, adapter)
			require.NotNil(t, redeem)
			defer redeem()
			var expectedAdapter AdapterType
			require.IsType(t, expectedAdapter, adapter)

			jazon, err := adapter.Marshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run("should have called the adapter's MarshalJSON method", func(t *testing.T) {
				_, ok := adapter.(AdapterType)
				require.True(t, ok)
				auditable, ok := adapter.(interface{ MarshalCallsLen() int })
				require.True(t, ok)
				calls := auditable.MarshalCallsLen()
				require.Equal(t, 1, calls)
			})

			t.Run("should have cached the route for this type", func(t *testing.T) {
				require.Len(t, reg.marshalerCache, 1)
				key := reflect.TypeOf(value)
				require.Contains(t, reg.marshalerCache, key)
				entry := reg.marshalerCache[key]
				mockAdapter, ok := adapter.(AdapterType)
				require.True(t, ok)
				require.Equal(t,
					fmt.Sprintf("github.com/go-openapi/swag/jsonutils/%s", reflect.Indirect(reflect.ValueOf(mockAdapter)).Type()),
					entry.Who,
				)
				require.True(t, entry.What.Has(ifaces.CapabilityMarshalJSON))
			})
		})
	}
}

type MockAdapter1 struct {
	*mocks.MockAdapter
}

func (a *MockAdapter1) UnmarshalCallsLen() int {
	return len(a.UnmarshalCalls())
}

func (a *MockAdapter1) MarshalCallsLen() int {
	return len(a.MarshalCalls())
}

type MockAdapter2 struct {
	*mocks.MockAdapter
}

func (a *MockAdapter2) UnmarshalCallsLen() int {
	return len(a.UnmarshalCalls())
}

func (a *MockAdapter2) MarshalCallsLen() int {
	return len(a.MarshalCalls())
}

func newMockAdapter1() *MockAdapter1 {
	return &MockAdapter1{
		MockAdapter: newMockAdapter(),
	}
}

func newMockAdapter2() *MockAdapter2 {
	return &MockAdapter2{
		MockAdapter: newMockAdapter(),
	}
}

func newMockAdapter() *mocks.MockAdapter {
	return &mocks.MockAdapter{
		MarshalFunc: func(_ any) ([]byte, error) {
			return []byte("null"), nil
		},
		NewOrderedMapFunc: func(_ int) ifaces.OrderedMap {
			return newMockOrdered()
		},
		OrderedMarshalFunc: func(_ ifaces.Ordered) ([]byte, error) {
			return []byte("null"), nil
		},
		OrderedUnmarshalFunc: func(_ []byte, _ ifaces.SetOrdered) error {
			return nil
		},
		UnmarshalFunc: func(_ []byte, _ any) error {
			return nil
		},
	}
}

func support1(_ ifaces.Capability, _ any) bool {
	return true
}

type supportedType struct {
}

func support2(capability ifaces.Capability, value any) bool {
	switch capability { //nolint:exhaustive
	case ifaces.CapabilityMarshalJSON:
		_, ok := value.(supportedType)
		return ok
	default:
		return false
	}
}

func register1(dispatcher ifaces.Registrar) {
	t := reflect.TypeOf(MockAdapter1{})

	dispatcher.RegisterFor(
		ifaces.RegistryEntry{
			Who:  fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()),
			What: ifaces.AllCapabilities,
			Constructor: func() ifaces.Adapter {
				return newMockAdapter1()
			},
			Redeemer: func(_ ifaces.Adapter) {},
			Support:  support1,
		})
}

func register2(dispatcher ifaces.Registrar) {
	t := reflect.TypeOf(MockAdapter2{})

	dispatcher.RegisterFor(
		ifaces.RegistryEntry{
			Who:  fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()),
			What: ifaces.AllCapabilities,
			Constructor: func() ifaces.Adapter {
				return newMockAdapter2()
			},
			Redeemer: func(_ ifaces.Adapter) {},
			Support:  support2,
		})
}

var _ ifaces.OrderedMap = &MockOrdered{}

type MockOrdered struct {
	mocks.MockOrdered
	mocks.MockSetOrdered
}

func (m *MockOrdered) OrderedMarshalJSON() ([]byte, error) {
	return nil, nil
}

func (m *MockOrdered) OrderedUnmarshalJSON([]byte) error {
	return nil
}

func newMockOrdered() *MockOrdered {
	return &MockOrdered{
		MockOrdered:    mocks.MockOrdered{},
		MockSetOrdered: mocks.MockSetOrdered{},
	}
}
