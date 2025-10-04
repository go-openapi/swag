// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestJSONNameIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		assert.NotNil(t, NewNameProvider())
	})
}
