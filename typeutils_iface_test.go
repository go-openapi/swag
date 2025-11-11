// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestTypeUtilsIface(t *testing.T) {
	t.Run("deprecated type utility functions should work", func(t *testing.T) {
		// only check happy path - more comprehensive testing is carried out inside the called package
		require.True(t, IsZero(0))
	})
}
