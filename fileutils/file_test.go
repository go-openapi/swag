// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils

import (
	"io"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestFileImplementsIOReader(t *testing.T) {
	var file any = &File{}
	expected := "that File implements io.Reader"
	assert.Implements(t, new(io.Reader), file, expected)
}

func TestFileImplementsIOReadCloser(t *testing.T) {
	var file any = &File{}
	expected := "that File implements io.ReadCloser"
	assert.Implements(t, new(io.ReadCloser), file, expected)
}
