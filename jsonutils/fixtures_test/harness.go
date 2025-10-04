// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fixtures

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	yaml "go.yaml.in/yaml/v3"
)

// embedded test files

//go:embed *.yaml
var EmbeddedFixtures embed.FS

// TestCases is a collection of test cases [Fixture]
type TestCases []Fixture

// Fixture holds a JSON payload and its equivalent YAML payload.
type Fixture struct {
	Name          string `yaml:"name"`
	Comment       string `yaml:"comment"`
	JSONPayload   string `yaml:"json_payload"`
	YAMLPayload   string `yaml:"yaml_payload"`
	Error         bool   `yaml:"error"`
	ErrorContains string `yaml:"error_contains"`
}

// JSONBytes returns the JSON payload string as bytes.
func (f Fixture) JSONBytes() []byte {
	return []byte(f.JSONPayload)
}

// YAMLBytes returns the YAML payload string as bytes.
func (f Fixture) YAMLBytes() []byte {
	return []byte(f.YAMLPayload)
}

// ExpectError indicates if this test case expect an error.
func (f Fixture) ExpectError() bool {
	return f.Error
}

// Harness is a test helper to retrieve or scan over a collection
// of predefined test cases loaded from the embedded file system.
type Harness struct {
	t     testing.TB
	index map[string]*Fixture
	tests TestCases
}

// NewHarness yields a new test [Harness] for a given test.
//
// The [Harness] requires a call to [Harness.Init] to be ready.
func NewHarness(t testing.TB) *Harness {
	return &Harness{
		t: t,
	}
}

type testFile struct {
	TestCases TestCases `yaml:"testcases"`
}

// Init loads the set of fixtures from the embedded YAML configuration file:
// "ordered_fixtures.yaml".
func (h *Harness) Init() {
	const sensibleAlloc = 20
	fixtures := ShouldLoadFixture(h.t, EmbeddedFixtures, "ordered_fixtures.yaml")

	testCases := testFile{
		TestCases: make(TestCases, 0, sensibleAlloc),
	}
	require.NoError(h.t, yaml.Unmarshal(fixtures, &testCases))

	h.tests = testCases.TestCases
	h.index = make(map[string]*Fixture, len(h.tests))
	for i := range h.tests {
		name := h.tests[i].Name
		h.index[name] = &h.tests[i]
	}
}

func (h *Harness) Get(name string) (Fixture, bool) {
	fixture, ok := h.index[name]

	return *fixture, ok
}

func (h *Harness) ShouldGet(name string) Fixture {
	fixture, ok := h.Get(name)
	require.True(h.t, ok)

	return fixture
}

type Filter func(*filters)

type filters struct {
	withoutError      bool
	withError         bool
	withExcludeRegexp *regexp.Regexp
	withIncludeRegexp *regexp.Regexp
}

func WithError(only bool) Filter {
	return func(f *filters) {
		f.withError = only
	}
}

func WithoutError(only bool) Filter {
	return func(f *filters) {
		f.withoutError = only
	}
}

func WithExcludePattern(rex *regexp.Regexp) Filter {
	return func(f *filters) {
		f.withExcludeRegexp = rex
	}
}

func WithIncludePattern(rex *regexp.Regexp) Filter {
	return func(f *filters) {
		f.withIncludeRegexp = rex
	}
}

func (h *Harness) AllTests(filter ...Filter) iter.Seq2[string, Fixture] {
	var f filters
	for _, apply := range filter {
		apply(&f)
	}

	return func(yield func(string, Fixture) bool) {
		for _, test := range h.tests {
			if f.withoutError && test.Error || f.withError && !test.Error {
				continue
			}
			if f.withExcludeRegexp != nil && f.withExcludeRegexp.MatchString(test.Name) {
				continue
			}
			if f.withIncludeRegexp != nil && !f.withIncludeRegexp.MatchString(test.Name) {
				continue
			}

			key := test.Name
			value := test

			if !yield(key, value) {
				return
			}
		}
	}
}

func ShouldLoadFixture(t testing.TB, fsys fs.FS, pth string) []byte {
	data, err := loadFixture(fsys, pth)
	require.NoError(t, err)

	return data
}

func MustLoadFixture(fsys fs.FS, pth string) []byte {
	data, err := loadFixture(fsys, pth)
	if err != nil {
		panic(err)
	}

	return data
}

func loadFixture(fsys fs.FS, pth string) ([]byte, error) {
	pth = filepath.ToSlash(pth) // "/" even on windows

	return fs.ReadFile(fsys, pth)
}

// JSONEqualOrdered is a replacement for [require.JSONEq] that checks further that
// two JSONs are exactly equal, with only the following tolerated differences:
//
//   - non-significant white space
//   - numerical encoding (e.g. 0.01 <=> 1e-2)
//   - unicode encoding (e.g. explicitly escaped unicode sequences <=> unicode rune)
func JSONEqualOrdered(t testing.TB, expected, actual string) {
	t.Helper()

	bufExpected := bytes.NewBufferString(expected)
	decExpected := json.NewDecoder(bufExpected)
	expectedTokens := make([]json.Token, 0)

	bufActual := bytes.NewBufferString(actual)
	decActual := json.NewDecoder(bufActual)

	for {
		tok, err := decExpected.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			require.NoError(t, err)
			return
		}

		expectedTokens = append(expectedTokens, tok)
	}

	count := 0
	for {
		tok, err := decActual.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			require.NoError(t, err)
			return
		}
		require.Less(t, count, len(expectedTokens))
		require.Equalf(t, expectedTokens[count], tok, "json token differs: %d", count)
		count++
	}
}

var (
	rexStripIndent        = regexp.MustCompile(`(?m)^\s+|^-{3}$`)
	rexStripComment       = regexp.MustCompile(`(?m)\s*#.+$`)
	rexStripInlineComment = regexp.MustCompile(`(?m)(.+?)\s+#.+$`)
	rexStripEmpty         = regexp.MustCompile(`(?m)^\s*$`)
)

// YAMLEqualOrdered is a replacement for [require.YAMLEq] that checks further that
// two YAML are exactly equal, with only the following tolerated differences:
//
//   - non-significant white space
//   - comments
//
// Otherwise, the representation of arrays, null values and objects must match exactly, e.g
// this checks tells us that:
//
//	a: [1,2,3]
//
// differs from:
//
//	a:
//	  - 1
//	  - 2
//	  - 3
//
// even though we know that they have equivalent semantics.
//
// NOTE: at this moment, this check does not support anchors.
func YAMLEqualOrdered(t testing.TB, expected, actual string) {
	t.Helper()

	RequireYAMLEq(t, expected, actual) // necessary but not sufficient condition

	// strip all indentation and comments (anchors not supported)
	strippedExpected := rexStripIndent.ReplaceAllString(expected, "")
	strippedExpected = rexStripComment.ReplaceAllString(strippedExpected, "")
	strippedExpected = rexStripInlineComment.ReplaceAllString(strippedExpected, "$1")
	strippedExpected = rexStripEmpty.ReplaceAllString(strippedExpected, "")

	strippedActual := rexStripIndent.ReplaceAllString(expected, "")
	strippedActual = rexStripComment.ReplaceAllString(strippedActual, "")
	strippedActual = rexStripInlineComment.ReplaceAllString(strippedActual, "$1")
	strippedActual = rexStripEmpty.ReplaceAllString(strippedActual, "")

	require.Equal(t, strippedExpected, strippedActual)
}

// RequireYAMLEq is the same as [require.YAMLEq] but without the dependency to go.pkg.in/yaml.v3.
//
// NOTE: this could be reverted once https://github.com/stretchr/testify/pull/1772 is merged.
func RequireYAMLEq(t testing.TB, expected string, actual string, msgAndArgs ...any) {
	t.Helper()

	if AssertYAMLEq(t, expected, actual, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// AssertYAMLEq is the same as [assert.YAMLEq] but without the dependency to go.pkg.in/yaml.v3.
//
// NOTE: this could be reverted once https://github.com/stretchr/testify/pull/1772 is merged.
func AssertYAMLEq(t testing.TB, expected string, actual string, msgAndArgs ...any) bool {
	t.Helper()
	var expectedYAMLAsInterface, actualYAMLAsInterface any

	if err := yaml.Unmarshal([]byte(expected), &expectedYAMLAsInterface); err != nil {
		return assert.Fail(t, fmt.Sprintf("Expected value ('%s') is not valid yaml.\nYAML parsing error: '%s'", expected, err.Error()), msgAndArgs...)
	}

	// Shortcut if same bytes
	if actual == expected {
		return true
	}

	if err := yaml.Unmarshal([]byte(actual), &actualYAMLAsInterface); err != nil {
		return assert.Fail(t, fmt.Sprintf("Input ('%s') needs to be valid yaml.\nYAML error: '%s'", actual, err.Error()), msgAndArgs...)
	}

	return assert.Equal(t, expectedYAMLAsInterface, actualYAMLAsInterface, msgAndArgs...)
}
