// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package stringutils

import (
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

const collectionFormatComma = "csv"

func TestSplitByFormat(t *testing.T) {
	expected := []string{"one", "two", "three"}
	for _, fmt := range []string{collectionFormatComma, collectionFormatPipe, collectionFormatTab, collectionFormatSpace, collectionFormatMulti} {

		var actual []string
		switch fmt {
		case collectionFormatMulti:
			assert.Nil(t, SplitByFormat("", fmt))
			assert.Nil(t, SplitByFormat("blah", fmt))
		case collectionFormatSpace:
			actual = SplitByFormat(strings.Join(expected, " "), fmt)
			assert.Equal(t, expected, actual)
		case collectionFormatPipe:
			actual = SplitByFormat(strings.Join(expected, "|"), fmt)
			assert.Equal(t, expected, actual)
		case collectionFormatTab:
			actual = SplitByFormat(strings.Join(expected, "\t"), fmt)
			assert.Equal(t, expected, actual)
		default:
			actual = SplitByFormat(strings.Join(expected, ","), fmt)
			assert.Equal(t, expected, actual)
		}
	}
}

func TestJoinByFormat(t *testing.T) {
	for _, fmt := range []string{collectionFormatComma, collectionFormatPipe, collectionFormatTab, collectionFormatSpace, collectionFormatMulti} {

		lval := []string{"one", "two", "three"}
		var expected []string
		switch fmt {
		case collectionFormatMulti:
			expected = lval
		case collectionFormatSpace:
			expected = []string{strings.Join(lval, " ")}
		case collectionFormatPipe:
			expected = []string{strings.Join(lval, "|")}
		case collectionFormatTab:
			expected = []string{strings.Join(lval, "\t")}
		default:
			expected = []string{strings.Join(lval, ",")}
		}
		assert.Nil(t, JoinByFormat(nil, fmt))
		assert.Equal(t, expected, JoinByFormat(lval, fmt))
	}
}
