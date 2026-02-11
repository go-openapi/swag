// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package stringutils

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestContainsStringsCI(t *testing.T) {
	list := []string{"hello", "world", "and", "such"}

	assert.TrueT(t, ContainsStringsCI(list, "hELLo"))
	assert.TrueT(t, ContainsStringsCI(list, "world"))
	assert.TrueT(t, ContainsStringsCI(list, "AND"))
	assert.FalseT(t, ContainsStringsCI(list, "nuts"))
}

func TestContainsStrings(t *testing.T) {
	list := []string{"hello", "world", "and", "such"}

	assert.TrueT(t, ContainsStrings(list, "hello"))
	assert.FalseT(t, ContainsStrings(list, "hELLo"))
	assert.TrueT(t, ContainsStrings(list, "world"))
	assert.FalseT(t, ContainsStrings(list, "World"))
	assert.TrueT(t, ContainsStrings(list, "and"))
	assert.FalseT(t, ContainsStrings(list, "AND"))
	assert.FalseT(t, ContainsStrings(list, "nuts"))
}
