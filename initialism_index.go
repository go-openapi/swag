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

package swag

import (
	"sort"
	"strings"
	"sync"
)

var (
	// commonInitialisms are common acronyms that are kept as whole uppercased words.
	commonInitialisms *indexOfInitialisms

	// initialisms is a slice of sorted initialisms
	initialisms []string

	// a copy of initialisms pre-baked as []rune
	initialismsRunes      [][]rune
	initialismsUpperCased [][]rune
	initialismsPluralForm []pluralForm // pre-baked indexed support for pluralization

	isInitialism func(string) bool

	maxAllocMatches int
)

func init() {
	// List of initialisms taken from https://github.com/golang/lint/blob/3390df4df2787994aea98de825b964ac7944b817/lint.go#L732-L769
	//
	// Now superseded by: https://github.com/mgechev/revive/blob/master/lint/name.go#L93
	//
	// Notice that initialisms are not necessarily uppercased.
	// In particular, we may find plural forms with mixed case like "IDs" or legit values like "IPv4" or "IPv6".
	//
	// At this moment, we don't support pluralization of terms that ends with an 's' (or 'S').
	// We don't want to support pluralization of terms which would otherwise conflict with another one,
	// like "HTTPs" vs "HTTPS". All these should be considered invariant. Hence: "Https" matches "HTTPS" and
	// "HTTPSS" is "HTTPS" followed by "S".
	configuredInitialisms := map[string]bool{
		// initialism: true|false = accept a pluralized form 'Xs' - false means invariant plural
		"ACL":   true,
		"API":   true,
		"ASCII": true,
		"CPU":   true,
		"CSS":   false,
		"DNS":   false,
		"EOF":   true,
		"GUID":  true,
		"HTML":  true,
		"HTTPS": false,
		"HTTP":  false,
		"ID":    true,
		"IP":    true,
		"IPv4":  true, // prefer the mixed case outcome IPv4 over the capitalized IPV4
		"IPv6":  true, // prefer the mixed case outcome
		"JSON":  true,
		"LHS":   true,
		"OAI":   true, // not in the linter's list, but added for the openapi context
		"QPS":   false,
		"RAM":   true,
		"RHS":   false,
		"RPC":   true,
		"SLA":   true,
		"SMTP":  true,
		"SQL":   true,
		"SSH":   true,
		"TCP":   true,
		"TLS":   false,
		"TTL":   true,
		"UDP":   true,
		"UI":    true,
		"UID":   true,
		"UUID":  true,
		"URI":   true,
		"URL":   true,
		"UTF8":  true,
		"VM":    true,
		"XML":   true,
		"XMPP":  true,
		"XSRF":  true,
		"XSS":   false,
	}

	// a thread-safe index of initialisms
	commonInitialisms = newIndexOfInitialisms().load(configuredInitialisms)
	initialisms = commonInitialisms.sorted()
	initialismsRunes = asRunes(initialisms)
	initialismsUpperCased = asUpperCased(initialisms)
	maxAllocMatches = maxAllocHeuristic(initialismsRunes)
	initialismsPluralForm = asPluralForms(initialisms, commonInitialisms)

	// a test function
	isInitialism = commonInitialisms.isInitialism
}

func asRunes(in []string) [][]rune {
	out := make([][]rune, len(in))
	for i, initialism := range in {
		out[i] = []rune(initialism)
	}

	return out
}

func asUpperCased(in []string) [][]rune {
	out := make([][]rune, len(in))

	for i, initialism := range in {
		out[i] = []rune(upper(trim(initialism)))
	}

	return out
}

// asPluralForms bakes an index of pluralization support.
func asPluralForms(in []string, idx *indexOfInitialisms) []pluralForm {
	out := make([]pluralForm, len(in))
	for i, initialism := range in {
		out[i] = idx.pluralForm(initialism)
	}

	return out
}

func maxAllocHeuristic(in [][]rune) int {
	heuristic := make(map[rune]int)
	for _, initialism := range in {
		heuristic[initialism[0]]++
	}

	var maxAlloc int
	for _, val := range heuristic {
		if val > maxAlloc {
			maxAlloc = val
		}
	}

	return maxAlloc
}

// AddInitialisms add additional initialisms.
//
// This method adds extra words as "initialisms" (i.e. words that won't be camel cased or titled cased),
// to the existing list of common initialisms (such as ID, HTTP...).
//
// The list of initialisms is maintained at the package level, so this method can't be used concurrently.
//
// It is typically used when initializing a command line utility, such as go-swagger.
func AddInitialisms(words ...string) {
	for _, word := range words {
		// commonInitialisms[upper(word)] = true
		uword := upper(word)
		commonInitialisms.add(uword, !strings.HasSuffix(uword, "S"))
	}
	// sort again
	initialisms = commonInitialisms.sorted()
	initialismsRunes = asRunes(initialisms)
	initialismsUpperCased = asUpperCased(initialisms)
	initialismsPluralForm = asPluralForms(initialisms, commonInitialisms)
}

// indexOfInitialisms is a thread-safe implementation of the sorted index of initialisms.
// Since go1.9, this may be implemented with sync.Map.
type indexOfInitialisms struct {
	sortMutex *sync.Mutex
	index     *sync.Map
}

func newIndexOfInitialisms() *indexOfInitialisms {
	return &indexOfInitialisms{
		sortMutex: new(sync.Mutex),
		index:     new(sync.Map),
	}
}

func (m *indexOfInitialisms) load(initial map[string]bool) *indexOfInitialisms {
	m.sortMutex.Lock()
	defer m.sortMutex.Unlock()
	for k, v := range initial {
		m.index.Store(k, v)
	}
	return m
}

func (m *indexOfInitialisms) isInitialism(key string) bool {
	_, ok := m.index.Load(key)
	return ok
}

func (m *indexOfInitialisms) add(key string, hasPlural bool) *indexOfInitialisms {
	m.index.Store(key, hasPlural)
	return m
}

func (m *indexOfInitialisms) sorted() (result []string) {
	m.sortMutex.Lock()
	defer m.sortMutex.Unlock()
	m.index.Range(func(key, _ interface{}) bool {
		k := key.(string)
		result = append(result, k)
		return true
	})
	sort.Sort(sort.Reverse(byInitialism(result)))
	return
}

// pluralForm denotes the kind of pluralization to be used for initialisms.
//
// At this moment, initialisms are either invariant or follow a simple plural form with an
// extra (lower case) "s".
type pluralForm uint8

const (
	notPlural pluralForm = iota
	invariantPlural
	simplePlural
)

// pluralForm indicates how we want to pluralize a given initialism.
//
// Besides configured invariant forms (like HTTP and HTTPS),
// an initialism is normally pluralized by adding a single 's', like in IDs.
//
// Initialisms ending with an 'S' or an 's' are configured as invariant (we don't
// support plural forms like CSSes or DNSes, however the mechanism could be extended to
// do just that).
func (m *indexOfInitialisms) pluralForm(key string) pluralForm {
	v, ok := m.index.Load(key)
	if !ok {
		return notPlural
	}

	acceptsPlural := v.(bool)
	if !acceptsPlural {
		return invariantPlural
	}

	return simplePlural
}

type byInitialism []string

func (s byInitialism) Len() int {
	return len(s)
}
func (s byInitialism) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less specifies the order in which initialisms are prioritized:
// 1. match longest first
// 2. when equal length, match in reverse lexicographical order, lower case match comes first
func (s byInitialism) Less(i, j int) bool {
	if len(s[i]) != len(s[j]) {
		return len(s[i]) < len(s[j])
	}

	return s[i] < s[j]
}
