package testintegration

import (
	"fmt"
	"maps"
	"os"
	"testing"

	"github.com/go-openapi/swag/jsonutils"
	"github.com/go-openapi/swag/jsonutils/adapters"
	easyjson "github.com/go-openapi/swag/jsonutils/adapters/easyjson/json"
	"github.com/go-openapi/swag/jsonutils/adapters/ifaces"
	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"github.com/stretchr/testify/require"
)

var (
	_ EJMarshaler   = &EasyOrderedObject{}
	_ EJUnmarshaler = &EasyOrderedObject{}
	_ EJMarshaler   = &EasyObject{}
	_ EJUnmarshaler = &EasyObject{}
	_ EJMarshaler   = &EasyTarget{}
	_ EJUnmarshaler = &EasyTarget{}
)

func TestMain(m *testing.M) {
	easyjson.Register(adapters.Registry)

	os.Exit(m.Run())
}

type EasyObject struct {
	*MockEJMarshaler
	*MockEJUnmarshaler

	inner easyjson.MapSlice
}

func newEasyObject() *EasyObject {
	a := &EasyObject{} // we may leave the inner member to nil

	a.MockEJMarshaler = &MockEJMarshaler{
		MarshalEasyJSONFunc: func(w *jwriter.Writer) {
			a.inner.MarshalEasyJSON(w)
		},
	}

	a.MockEJUnmarshaler = &MockEJUnmarshaler{
		UnmarshalEasyJSONFunc: func(l *jlexer.Lexer) {
			a.inner.UnmarshalEasyJSON(l)

			// reshuffle mappings: all inner MapSlice's become maps
			if len(a.inner) == 0 {
				return
			}

			changed := make(map[string]any)
			for k, v := range a.inner.OrderedItems() {
				ordered, ok := v.(ifaces.Ordered)
				if !ok {
					continue
				}
				m := maps.Collect(ordered.OrderedItems())
				changed[k] = m
			}
			a.inner.SetOrderedItems(maps.All(changed))
		},
	}

	return a
}

type EasyTarget struct {
	*MockEJMarshaler
	*MockEJUnmarshaler
}

type EasyOrderedObject struct {
	*MockEJMarshaler
	*MockEJUnmarshaler

	inner easyjson.MapSlice
}

func newEasyOrderedObject() *EasyOrderedObject {
	a := &EasyOrderedObject{
		// we may leave the inner member to nil
	}

	a.MockEJMarshaler = &MockEJMarshaler{
		MarshalEasyJSONFunc: func(w *jwriter.Writer) {
			a.inner.MarshalEasyJSON(w)
		},
	}

	a.MockEJUnmarshaler = &MockEJUnmarshaler{
		UnmarshalEasyJSONFunc: func(l *jlexer.Lexer) {
			a.inner.UnmarshalEasyJSON(l)
		},
	}

	return a
}

func (o EasyOrderedObject) MarshalEasyJSON(w *jwriter.Writer) {
	o.MockEJMarshaler.MarshalEasyJSON(w)
}

func (o *EasyOrderedObject) UnmarshalEasyJSON(l *jlexer.Lexer) {
	o.MockEJUnmarshaler.UnmarshalEasyJSON(l)
}

type EasyOrderedTarget struct {
	easyjson.MapSlice
	*MockEJMarshaler
	*MockEJUnmarshaler
}

func (o *EasyOrderedTarget) MarshalEasyJSON(w *jwriter.Writer) {
	o.MockEJMarshaler.MarshalEasyJSON(w)
}

func (o *EasyOrderedTarget) UnmrshalEasyJSON(l *jlexer.Lexer) {
	o.MockEJUnmarshaler.UnmarshalEasyJSON(l)
}

type assertionType uint8

const (
	assertionTypeUnordered assertionType = iota
	assertionTypeOrdered
)

type option func(*options)

type assertion func(v any) func(*testing.T)
type options struct {
	readAssertions  []assertion
	writeAssertions []assertion
}

func withAssertionsAfterRead(fn ...assertion) option {
	return func(o *options) {
		o.readAssertions = append(o.readAssertions, fn...)
	}
}

func withAssertionsAfterWrite(fn ...assertion) option {
	return func(o *options) {
		o.writeAssertions = append(o.writeAssertions, fn...)
	}
}

func runTestSuite[V any, T any](valueConstructor func() *V, targetConstructor func() *T, assertAs assertionType, extras ...option) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		harness := fixtures.NewHarness(t)
		harness.Init()

		for name, test := range harness.AllTests() {
			t.Run(name, testJSONTransforms(test, valueConstructor, targetConstructor, assertAs, extras...))
		}
	}
}

type Error string

func (e Error) Error() string { return string(e) }

const errTestConfig Error = "error in test config"

func testJSONTransforms[V any, T any](test fixtures.Fixture, valueConstructor func() *V, targetConstructor func() *T, assertAs assertionType, extras ...option) func(*testing.T) {
	var expectation string
	switch assertAs {
	case assertionTypeOrdered:
		expectation = "identical" // same JSON, with key order kept
	case assertionTypeUnordered:
		expectation = "equivalent" // same JSON, the order of keys notwithstanding
	default:
		panic(fmt.Errorf("invalid assertionType: %d: %w", assertAs, errTestConfig))
	}
	var o options
	for _, apply := range extras {
		apply(&o)
	}

	return func(t *testing.T) {
		t.Helper()

		t.Run(fmt.Sprintf("ReadJSON then WriteJSON should produce %s JSON", expectation), func(t *testing.T) {
			var value V
			if valueConstructor != nil {
				value = *valueConstructor()
			}

			if test.ExpectError() {
				require.Error(t, jsonutils.ReadJSON(test.JSONBytes(), &value))
				for _, fn := range o.readAssertions {
					if fn == nil {
						continue
					}
					fn(value)(t)
				}

				return
			}

			require.NoError(t, jsonutils.ReadJSON(test.JSONBytes(), &value))
			for _, fn := range o.readAssertions {
				if fn == nil {
					continue
				}
				fn(value)(t)
			}

			jazon, err := jsonutils.WriteJSON(value)
			require.NoError(t, err)
			for _, fn := range o.writeAssertions {
				if fn == nil {
					continue
				}
				fn(value)(t)
			}

			switch assertAs {
			case assertionTypeOrdered:
				fixtures.JSONEqualOrdered(t, test.JSONPayload, string(jazon))
			case assertionTypeUnordered:
				require.JSONEq(t, test.JSONPayload, string(jazon))
			}

			t.Run(fmt.Sprintf("FromDynamicJSON then WriteJSON should produce %s JSON", expectation), func(t *testing.T) {
				var target T
				if targetConstructor != nil {
					target = *targetConstructor()
				}

				require.NoError(t, jsonutils.FromDynamicJSON(value, &target))
				jazon, err := jsonutils.WriteJSON(target)
				require.NoError(t, err)
				for _, fn := range o.writeAssertions {
					if fn == nil {
						continue
					}
					fn(target)(t)
				}

				switch assertAs {
				case assertionTypeOrdered:
					fixtures.JSONEqualOrdered(t, test.JSONPayload, string(jazon))
				case assertionTypeUnordered:
					require.JSONEq(t, test.JSONPayload, string(jazon))
				}
			})
		})
	}
}
