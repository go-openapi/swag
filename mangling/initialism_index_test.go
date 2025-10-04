// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package mangling

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestInitialismSorted(t *testing.T) {
	configuredInitialisms := map[string]struct{}{
		"ACL":   {},
		"API":   {},
		"ASCII": {},
		"CPU":   {},
		"CSS":   {},
		"DNS":   {},
		"VM":    {},
		"XML":   {},
		"IPv4":  {},
		"IPV4":  {},
		"XMPP":  {},
		"XSRF":  {},
		"XSS":   {},
	}

	// now the order is reverse lexicographic.
	// With this ordering, when several initialisms differ in case only,
	// lowercase comes first.
	//
	// Example below: IPv4 and IPV4 favors IPv4.
	goldenSample := []string{
		"ASCII",
		"XSRF",
		"XMPP",
		"IPv4",
		"IPV4",
		"XSS",
		"XML",
		"DNS",
		"CSS",
		"CPU",
		"API",
		"ACL",
		"VM",
	}
	for i := 0; i < 50; i++ {
		idx := newIndexOfInitialisms()
		for w := range configuredInitialisms {
			idx.add(w) // add in random order
		}
		sample := idx.sorted()
		failMsg := "equal sorted initialisms should be always equal"

		if !assert.Equal(t, goldenSample, sample, failMsg) {
			t.FailNow()
		}
	}
}

func TestInitialismPlural(t *testing.T) {
	idx := newIndexOfInitialisms()
	for _, word := range DefaultInitialisms() {
		idx.add(word)
	}
	idx.add("Series")
	idx.add("Serie")

	assert.Equal(t, simplePlural, idx.pluralForm("ID"))
	assert.Equal(t, invariantPlural, idx.pluralForm("HTTP"))
	assert.Equal(t, invariantPlural, idx.pluralForm("HTTPS"))
	assert.Equal(t, invariantPlural, idx.pluralForm("DNS"))
	assert.Equal(t, invariantPlural, idx.pluralForm("Serie"))
	assert.Equal(t, invariantPlural, idx.pluralForm("Series"))
	assert.Equal(t, notPlural, idx.pluralForm("xyz"))
}

func TestInitialismSanitize(t *testing.T) {
	t.Run("should be ignored", func(t *testing.T) {
		idx := newIndexOfInitialisms()
		ignoredKeys := []string{
			"1",
			"+ABC",
		}

		for _, key := range ignoredKeys {
			idx.add(key)
			_, ok := idx.index[key]
			assert.Falsef(t, ok,
				"expected key %q not to be indexed", key,
			)
		}
	})

	t.Run("should be unique trimmed", func(t *testing.T) {
		idx := newIndexOfInitialisms()
		trimmedKeys := []string{
			" aBc ",
			" DeF",
			"DeF\t",
			"GHI\u2007",
			"\u2002GHI",
		}

		for _, key := range trimmedKeys {
			idx.add(key)
			_, ok := idx.index[key]
			assert.Falsef(t, ok,
				"expected key %q not to be indexed", key,
			)

			trimmedKey := strings.TrimSpace(key)
			require.Len(t, trimmedKey, 3) // ensure trimmed
			_, trimmedOk := idx.index[trimmedKey]
			assert.Truef(t, trimmedOk,
				"expected %q (trimmed as %q) to be indexed", key, trimmedKey,
			)
		}

		assert.Len(t, idx.index, 3)
	})

	t.Run("should be uppercased", func(t *testing.T) {
		const key = "abc"
		idx := newIndexOfInitialisms()
		idx.add(key)

		_, ok := idx.index[key]
		assert.False(t, ok)

		_, capitalizedOk := idx.index[strings.ToUpper(key)]
		assert.True(t, capitalizedOk)
	})
}
