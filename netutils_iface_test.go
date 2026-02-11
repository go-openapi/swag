// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestNetUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		host, port, err := SplitHostPort("localhost:1000")
		require.NoError(t, err)
		assert.EqualT(t, "localhost", host)
		assert.EqualT(t, 1000, port)
	})
}
