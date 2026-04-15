// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonname

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type testAltEmbedded struct {
	Nested string `json:"nested"`
}

type testAltDeep struct {
	Deep string `json:"deep"`
}

type testAltMiddle struct {
	testAltDeep

	Middle string `json:"middle"`
}

// testAltStruct exercises the stdlib-aligned field discovery rules:
//   - Name: explicitly tagged
//   - NotTheSame: tagged with a different json name
//   - Ignored: fully excluded via `json:"-"`
//   - DashField: stdlib quirk, literally named "-" in json
//   - Untagged: empty name in tag → keeps Go name
//   - Optional: options-only tag → keeps Go name
//   - NoTag: no tag at all → keeps Go name
//   - unexported: excluded
//   - testAltEmbedded: fields promoted to the parent
//   - testAltMiddle: embedded struct itself embedding another → transitively promoted
type testAltStruct struct {
	testAltEmbedded
	testAltMiddle

	Name       string `json:"name"`
	NotTheSame int64  `json:"plain"`
	Ignored    string `json:"-"`
	DashField  string `json:"-,"` //nolint:staticcheck  // deliberate: exercise stdlib "-," quirk
	Untagged   string `json:""`
	Optional   string `json:",omitempty"`
	NoTag      string
	unexported string //nolint:unused  // exercised to confirm it is filtered out
}

// testAltShadow verifies the depth-based conflict resolution: the outer field
// must win over one promoted from an embedded type.
type testAltShadow struct {
	testAltEmbedded

	Nested string `json:"nested"`
}

func TestGoNameProvider(t *testing.T) {
	provider := NewGoNameProvider()
	obj := testAltStruct{}
	tpe := reflect.TypeOf(obj)
	ptr := &obj

	t.Run("GetGoName resolves tagged fields", func(t *testing.T) {
		for _, tc := range []struct {
			jsonName string
			goName   string
		}{
			{"name", "Name"},
			{"plain", "NotTheSame"},
			{"-", "DashField"}, // stdlib `json:"-,"` quirk
			{"Untagged", "Untagged"},
			{"Optional", "Optional"},
			{"NoTag", "NoTag"},
			{"nested", "Nested"},
			{"middle", "Middle"},
			{"deep", "Deep"},
		} {
			nm, ok := provider.GetGoName(obj, tc.jsonName)
			assert.TrueT(t, ok, "expected json name %q to resolve", tc.jsonName)
			assert.EqualT(t, tc.goName, nm)
		}
	})

	t.Run("GetGoName rejects excluded or unknown names", func(t *testing.T) {
		for _, bad := range []string{"ignored", "Ignored", "unexported", "doesNotExist"} {
			nm, ok := provider.GetGoName(obj, bad)
			assert.FalseT(t, ok, "did not expect %q to resolve", bad)
			assert.Empty(t, nm)
		}
	})

	t.Run("GetGoNameForType mirrors GetGoName", func(t *testing.T) {
		nm, ok := provider.GetGoNameForType(tpe, "plain")
		assert.TrueT(t, ok)
		assert.EqualT(t, "NotTheSame", nm)

		_, ok = provider.GetGoNameForType(tpe, "doesNotExist")
		assert.FalseT(t, ok)
	})

	t.Run("GetGoName accepts pointer subjects", func(t *testing.T) {
		nm, ok := provider.GetGoName(ptr, "name")
		assert.TrueT(t, ok)
		assert.EqualT(t, "Name", nm)

		nm, ok = provider.GetGoName(ptr, "nested")
		assert.TrueT(t, ok)
		assert.EqualT(t, "Nested", nm)
	})

	t.Run("GetJSONName is the inverse mapping", func(t *testing.T) {
		for _, tc := range []struct {
			goName   string
			jsonName string
		}{
			{"Name", "name"},
			{"NotTheSame", "plain"},
			{"DashField", "-"},
			{"Untagged", "Untagged"},
			{"Optional", "Optional"},
			{"NoTag", "NoTag"},
			{"Nested", "nested"},
			{"Middle", "middle"},
			{"Deep", "deep"},
		} {
			nm, ok := provider.GetJSONName(obj, tc.goName)
			assert.TrueT(t, ok, "expected go name %q to resolve", tc.goName)
			assert.EqualT(t, tc.jsonName, nm)
		}

		_, ok := provider.GetJSONName(obj, "Ignored")
		assert.FalseT(t, ok)

		_, ok = provider.GetJSONNameForType(tpe, "DoesNotExist")
		assert.FalseT(t, ok)
	})

	t.Run("GetJSONNames lists every discoverable field exactly once", func(t *testing.T) {
		names := provider.GetJSONNames(ptr)
		sort.Strings(names)
		assert.Equal(t, []string{
			"-",
			"NoTag",
			"Optional",
			"Untagged",
			"deep",
			"middle",
			"name",
			"nested",
			"plain",
		}, names)
	})

	t.Run("index caches per type", func(t *testing.T) {
		// Re-query to confirm no duplicate entries are created on repeat access.
		_, _ = provider.GetGoName(obj, "name")
		_, _ = provider.GetGoName(ptr, "name")
		assert.Len(t, provider.index, 1)
	})
}

// TestGoNameProvider_ShadowingMatchesStdlib pins our field selection to the
// behavior of encoding/json for shadowed promoted fields.
func TestGoNameProvider_ShadowingMatchesStdlib(t *testing.T) {
	provider := NewGoNameProvider()
	payload := `{"nested":"outer"}`

	var s testAltShadow
	require.NoError(t, json.Unmarshal([]byte(payload), &s))
	assert.Equal(t, "outer", s.Nested)
	assert.Empty(t, s.testAltEmbedded.Nested)

	goName, ok := provider.GetGoName(s, "nested")
	require.True(t, ok)
	// The outer field wins, exactly like encoding/json would pick s.Nested.
	assert.Equal(t, "Nested", goName)

	names := provider.GetJSONNames(s)
	assert.Len(t, names, 1)
}

// TestGoNameProvider_ImplementsInterface is a compile-time-ish guard that both
// providers agree on the core lookup shape expected by consumers.
func TestGoNameProvider_ImplementsInterface(t *testing.T) {
	var p providerIface = NewGoNameProvider()
	_, ok := p.GetGoName(testAltStruct{}, "name")
	assert.True(t, ok)
}

// Fixtures for the embedded-type promotion scenarios.

type testAltInner struct {
	Foo string `json:"foo"`
	Bar string
}

type testAltPromoted struct {
	testAltInner

	Baz string `json:"baz"`
}

type testAltTaggedEmbed struct {
	testAltInner `json:"inner"`

	Baz string `json:"baz"`
}

type testAltPtrEmbed struct {
	*testAltInner

	Baz string `json:"baz"`
}

type testAltUnexportedEmbed struct {
	testAltInner // exported type, will still promote

	inner testAltInner //nolint:unused  // regular unexported field, must be ignored
}

// TestGoNameProvider_EmbeddedPromotion validates how the provider resolves
// fields coming from an exported embedded type, mirroring encoding/json.
func TestGoNameProvider_EmbeddedPromotion(t *testing.T) {
	t.Run("untagged embedded struct promotes its fields", func(t *testing.T) {
		provider := NewGoNameProvider()
		obj := testAltPromoted{}

		for _, tc := range []struct {
			jsonName string
			goName   string
		}{
			{"foo", "Foo"}, // promoted, tagged on Inner
			{"Bar", "Bar"}, // promoted, untagged on Inner -> Go name kept
			{"baz", "Baz"}, // declared on Outer
		} {
			nm, ok := provider.GetGoName(obj, tc.jsonName)
			assert.TrueT(t, ok, "expected %q to resolve", tc.jsonName)
			assert.EqualT(t, tc.goName, nm)
		}

		// "Inner" must NOT appear as its own json name: its fields were promoted.
		_, ok := provider.GetJSONName(obj, "testAltInner")
		assert.False(t, ok)

		names := provider.GetJSONNames(obj)
		sort.Strings(names)
		assert.Equal(t, []string{"Bar", "baz", "foo"}, names)
	})

	t.Run("tagged embedded struct is treated as a regular named field", func(t *testing.T) {
		provider := NewGoNameProvider()
		obj := testAltTaggedEmbed{}

		nm, ok := provider.GetGoName(obj, "inner")
		assert.TrueT(t, ok)
		assert.EqualT(t, "testAltInner", nm)

		// With the tag in place, Inner's fields are NOT promoted.
		_, ok = provider.GetGoName(obj, "foo")
		assert.False(t, ok)
		_, ok = provider.GetGoName(obj, "Bar")
		assert.False(t, ok)

		names := provider.GetJSONNames(obj)
		sort.Strings(names)
		assert.Equal(t, []string{"baz", "inner"}, names)
	})

	t.Run("pointer-to-struct embedded is promoted like its elem", func(t *testing.T) {
		provider := NewGoNameProvider()
		obj := testAltPtrEmbed{}

		nm, ok := provider.GetGoName(obj, "foo")
		assert.TrueT(t, ok)
		assert.EqualT(t, "Foo", nm)

		names := provider.GetJSONNames(obj)
		sort.Strings(names)
		assert.Equal(t, []string{"Bar", "baz", "foo"}, names)
	})

	t.Run("regular unexported field alongside promotion does not leak", func(t *testing.T) {
		provider := NewGoNameProvider()
		obj := testAltUnexportedEmbed{}

		// Promotion still works for the exported embedded type.
		nm, ok := provider.GetGoName(obj, "foo")
		assert.TrueT(t, ok)
		assert.EqualT(t, "Foo", nm)

		// The regular unexported "inner" field must be invisible.
		_, ok = provider.GetGoName(obj, "inner")
		assert.False(t, ok)
	})

	t.Run("agrees with encoding/json on roundtrip", func(t *testing.T) {
		provider := NewGoNameProvider()
		payload := `{"foo":"f","Bar":"b","baz":"z"}`

		var stdVal testAltPromoted
		require.NoError(t, json.Unmarshal([]byte(payload), &stdVal))
		assert.Equal(t, "f", stdVal.Foo)
		assert.Equal(t, "b", stdVal.Bar)
		assert.Equal(t, "z", stdVal.Baz)

		// For every json key encoding/json accepted, the provider must resolve it too.
		for _, key := range []string{"foo", "Bar", "baz"} {
			_, ok := provider.GetGoName(stdVal, key)
			assert.TrueT(t, ok, "provider should resolve %q like encoding/json", key)
		}
	})
}
