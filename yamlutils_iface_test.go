// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package swag

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestYAMLUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		t.Run("with YAML bytes to document and back as JSON", func(t *testing.T) {
			const ydoc = "x:\n  a: one\n  b: two\n"
			doc, err := BytesToYAMLDoc([]byte(ydoc))
			require.NoError(t, err)

			buf, err := YAMLToJSON(doc)
			require.NoError(t, err)

			require.JSONEq(t, `{"x":{"a":"one","b":"two"}}`, string(buf))
		})
	})
}
