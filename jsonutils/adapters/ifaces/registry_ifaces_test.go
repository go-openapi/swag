package ifaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryIfaces(t *testing.T) {
	t.Run("capability should be a Stringer for debugging and error formatting purpose", func(t *testing.T) {
		for _, test := range []struct {
			in       Capability
			expected string
		}{
			{
				in:       CapabilityMarshalJSON,
				expected: "MarshalJSON",
			},
			{
				in:       CapabilityUnmarshalJSON,
				expected: "UnmarshalJSON",
			},
			{
				in:       CapabilityOrderedMarshalJSON,
				expected: "OrderedMarshalJSON",
			},
			{
				in:       CapabilityOrderedUnmarshalJSON,
				expected: "OrderedUnmarshalJSON",
			},
			{
				in:       CapabilityOrderedMap,
				expected: "OrderedMap",
			},
			{
				in:       Capability(99),
				expected: "<unknown>",
			},
		} {
			assert.Equal(t, test.expected, test.in.String())
		}
	})

	t.Run("capabilities should be a Stringer for debugging and error formatting purpose", func(t *testing.T) {
		for _, test := range []struct {
			in       Capabilities
			expected string
		}{
			{
				in:       AllCapabilities,
				expected: "MarshalJSON|UnmarshalJSON|OrderedMarshalJSON|OrderedUnmarshalJSON|OrderedMap",
			},
			{
				in:       AllUnorderedCapabilities,
				expected: "MarshalJSON|UnmarshalJSON",
			},
			{
				in:       Capabilities(CapabilityMarshalJSON | CapabilityOrderedMap),
				expected: "MarshalJSON|OrderedMap",
			},
			{
				in:       Capabilities(0),
				expected: "",
			},
		} {
			assert.Equal(t, test.expected, test.in.String())
		}
	})
}
