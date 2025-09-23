package benchmarks

import (
	"embed"
	"testing"

	"github.com/go-openapi/swag/jsonutils"
	"github.com/go-openapi/swag/jsonutils/adapters"
	easyjson "github.com/go-openapi/swag/jsonutils/adapters/easyjson/json"
	fixtures "github.com/go-openapi/swag/jsonutils/fixtures_test"
)

//go:embed fixtures/*.json
var EmbeddedFixtures embed.FS

var (
	smallJSON  []byte
	mediumJSON []byte
	largeJSON  []byte
)

type benchmarkContext struct {
	small  *SmallPayload
	medium *MediumPayload
	large  *LargePayload
}

func (c *benchmarkContext) Small() *SmallPayload {
	return c.small
}
func (c *benchmarkContext) Medium() *MediumPayload {
	return c.medium
}
func (c *benchmarkContext) Large() *LargePayload {
	return c.large
}

func BenchmarkJSON(b *testing.B) {
	ctx := initBenchmarks(b)

	b.ResetTimer()

	// stdlib is registered by default: it ignores MarshalEasyJSON and UnmarshalEasyJSON
	b.Run("with standard library", allBenchs(ctx, "standard"))

	adapters.Registry.Reset()
	easyjson.Register(adapters.Registry)

	// now easyjson is registered: it prioritizes MarshalEasyJSON and UnmarshalEasyJSON
	b.Run("with easyjson library", allBenchs(ctx, "easyjson"))
}

func initBenchmarks(b *testing.B) *benchmarkContext {
	smallJSON = fixtures.ShouldLoadFixture(b, EmbeddedFixtures, "fixtures/small_sample.json")
	mediumJSON = fixtures.ShouldLoadFixture(b, EmbeddedFixtures, "fixtures/medium_sample.json")
	largeJSON = fixtures.ShouldLoadFixture(b, EmbeddedFixtures, "fixtures/large_sample.json")

	return &benchmarkContext{
		small:  NewSmallPayload(),
		medium: NewMediumPayload(),
		large:  NewLargePayload(),
	}
}

func allBenchs(ctx *benchmarkContext, library string) func(*testing.B) {
	return func(b *testing.B) {
		b.Run(library+" ReadJSON - small", benchRead[SmallPayload](smallJSON))
		b.Run(library+" WriteJSON - small", benchWrite(ctx.Small))
		b.Run(library+" FromDynamicJSON - small", benchFromDynamic(ctx.Small))
		b.Run(library+" ReadJSON - medium", benchRead[MediumPayload](mediumJSON))
		b.Run(library+" WriteJSON - medium", benchWrite(ctx.Medium))
		b.Run(library+" FromDynamicJSON - medium", benchFromDynamic(ctx.Medium))
		b.Run(library+" ReadJSON - large", benchRead[LargePayload](largeJSON))
		b.Run(library+" WriteJSON - large", benchWrite(ctx.Large))
		b.Run(library+" FromDynamicJSON - large", benchFromDynamic(ctx.Large))
	}
}

func benchRead[T any](jazon []byte) func(*testing.B) {
	return func(b *testing.B) {
		b.SetBytes(int64(len(jazon)))

		for b.Loop() {
			var data T
			if err := jsonutils.ReadJSON(jazon, &data); err != nil {
				b.Logf("unexpected error: %v", err)
				b.FailNow()
			}
		}
	}
}

func benchWrite[T any](constructor func() *T) func(*testing.B) {
	data := constructor()
	return func(b *testing.B) {
		for b.Loop() {
			jazon, err := jsonutils.WriteJSON(data)
			if err != nil {
				b.Logf("unexpected error: %v", err)
				b.FailNow()
			}
			b.SetBytes(int64(len(jazon)))
		}
	}
}

func benchFromDynamic[T any](constructor func() *T) func(*testing.B) {
	source := constructor()
	return func(b *testing.B) {
		for b.Loop() {
			var target T

			if err := jsonutils.FromDynamicJSON(source, &target); err != nil {
				b.Logf("unexpected error: %v", err)
				b.FailNow()
			}
		}
	}
}
