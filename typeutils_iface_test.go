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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeUtilsIface(t *testing.T) {
	t.Run("deprecated type utility functions should work", func(t *testing.T) {
		// only check happy path - more comprehensive testing is carried out inside the called packag
		require.True(t, IsZero(0))
	})
}
