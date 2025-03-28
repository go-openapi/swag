// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringUtilsIface(t *testing.T) {
	t.Run("deprecated functions should work", func(t *testing.T) {
		assert.True(t, ContainsStrings([]string{"a", "b"}, "a"))
		assert.True(t, ContainsStringsCI([]string{"a", "b"}, "A"))
		require.Len(t, JoinByFormat([]string{"a", "b"}, "pipes"), 1)
		require.Len(t, SplitByFormat("a|b", "pipes"), 2)
	})
}
