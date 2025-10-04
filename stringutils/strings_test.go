// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package stringutils

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestContainsStringsCI(t *testing.T) {
	list := []string{"hello", "world", "and", "such"}

	assert.True(t, ContainsStringsCI(list, "hELLo"))
	assert.True(t, ContainsStringsCI(list, "world"))
	assert.True(t, ContainsStringsCI(list, "AND"))
	assert.False(t, ContainsStringsCI(list, "nuts"))
}

func TestContainsStrings(t *testing.T) {
	list := []string{"hello", "world", "and", "such"}

	assert.True(t, ContainsStrings(list, "hello"))
	assert.False(t, ContainsStrings(list, "hELLo"))
	assert.True(t, ContainsStrings(list, "world"))
	assert.False(t, ContainsStrings(list, "World"))
	assert.True(t, ContainsStrings(list, "and"))
	assert.False(t, ContainsStrings(list, "AND"))
	assert.False(t, ContainsStrings(list, "nuts"))
}
