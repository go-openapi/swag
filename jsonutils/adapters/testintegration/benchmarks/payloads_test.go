// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	stdjson "encoding/json"
	"fmt"
	"testing"

	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
	"github.com/go-openapi/testify/v2/require"
	"github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

func TestPayloads(t *testing.T) {
	t.Run("SmallPayload should ReadJSON and WriteJSON", verifyPayload(NewSmallPayload))
	t.Run("MediumPayload should ReadJSON and WriteJSON", verifyPayload(NewMediumPayload))
	t.Run("LargePayload should ReadJSON and WriteJSON", verifyPayload(NewLargePayload))
}

func TestFixtures(t *testing.T) {
	for i, jazon := range [][]byte{
		fixtures.ShouldLoadFixture(t, EmbeddedFixtures, "fixtures/small_sample.json"),
		fixtures.ShouldLoadFixture(t, EmbeddedFixtures, "fixtures/medium_sample.json"),
		fixtures.ShouldLoadFixture(t, EmbeddedFixtures, "fixtures/large_sample.json"),
	} {
		t.Run(fmt.Sprintf("[%d] json should be valid", i), func(t *testing.T) {
			var value any
			require.NoError(t, stdjson.Unmarshal(jazon, &value))
		})
	}
}

func verifyPayload[T any](constructor func() *T) func(*testing.T) {
	return func(t *testing.T) {
		value := constructor()

		t.Run(fmt.Sprintf("value of type %T should MarshalJSON", value), func(t *testing.T) {
			jazon, err := stdjson.Marshal(value)
			require.NoError(t, err)
			require.NotEmpty(t, jazon)

			t.Run(fmt.Sprintf("value of type %T should MarshalEasyJSON", value), func(t *testing.T) {
				var val any = value
				easyMarshaler, ok := val.(easyjson.Marshaler)
				require.True(t, ok)
				jw := jwriter.Writer{}
				easyMarshaler.MarshalEasyJSON(&jw)
				data, err := jw.BuildBytes()
				require.NoError(t, err)
				require.NotEmpty(t, data)
				require.JSONEqBytes(t, jazon, data)

				t.Run(fmt.Sprintf("value of type %T should UnmarshalEasyJSON", value), func(t *testing.T) {
					target := new(T)
					var tgt any = target
					easyUnmarshaler, ok := tgt.(easyjson.Unmarshaler)
					require.True(t, ok)
					jl := jlexer.Lexer{Data: data}
					easyUnmarshaler.UnmarshalEasyJSON(&jl)
					require.NoError(t, jl.Error())

					require.Equal(t, *value, *target)
				})
			})
		})
	}
}
