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

package jsonutils

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

// NOTE: the nolint:testifylint directives no longer apply on our fork.

func TestJSONConcatenation(t *testing.T) {
	t.Run("should concat nothing", func(t *testing.T) {
		assert.Nil(t, ConcatJSON())
	})

	// we require an exact assertion (with ordering), not just JSON equivalence. Hence: testifylint disabled.

	t.Run("should concat with nothing more", func(t *testing.T) {
		assert.Equal(t, []byte(`{"id":1}`),
			ConcatJSON([]byte(`{"id":1}`)),
		)
		assert.Equal(t, []byte(`{}`),
			ConcatJSON([]byte(`{}`), []byte(`{}`)),
		)
		assert.Equal(t, []byte(`[]`),
			ConcatJSON([]byte(`[]`), []byte(`[]`)),
		)
	})

	t.Run("should concat objects and arrays", func(t *testing.T) {
		assert.Equal(t, []byte(`{"id":1,"name":"Rachel"}`),
			ConcatJSON([]byte(`{"id":1}`), []byte(`{"name":"Rachel"}`)),
		)
		assert.Equal(t, []byte(`[{"id":1},{"name":"Rachel"}]`),
			ConcatJSON([]byte(`[{"id":1}]`), []byte(`[{"name":"Rachel"}]`)),
		)
		assert.Equal(t, []byte(`{"name":"Rachel"}`),
			ConcatJSON([]byte(`{}`), []byte(`{"name":"Rachel"}`)),
		)
		assert.Equal(t, []byte(`[{"name":"Rachel"}]`),
			ConcatJSON([]byte(`[]`), []byte(`[{"name":"Rachel"}]`)),
		)
		assert.Equal(t, []byte(`{"id":1}`),
			ConcatJSON([]byte(`{"id":1}`), []byte(`{}`)),
		)
		assert.Equal(t, []byte(`[{"id":1}]`),
			ConcatJSON([]byte(`[{"id":1}]`), []byte(`[]`)),
		)
		assert.Equal(t, []byte(`{"id":1,"name":"Rachel","age":32}`),
			ConcatJSON([]byte(`{"id":1}`), []byte(`{"name":"Rachel"}`), []byte(`{"age":32}`)),
		)
		assert.Equal(t, []byte(`[{"id":1},{"name":"Rachel"},{"age":32}]`),
			ConcatJSON([]byte(`[{"id":1}]`), []byte(`[{"name":"Rachel"}]`), []byte(`[{"age":32}]`)),
		)
		assert.Equal(t, []byte(`{"name":"Rachel","age":32}`),
			ConcatJSON([]byte(`{}`), []byte(`{"name":"Rachel"}`), []byte(`{"age":32}`)),
		)
		assert.Equal(t, []byte(`[{"name":"Rachel"},{"age":32}]`),
			ConcatJSON([]byte(`[]`), []byte(`[{"name":"Rachel"}]`), []byte(`[{"age":32}]`)),
		)
		assert.Equal(t, []byte(`{"id":1,"age":32}`),
			ConcatJSON([]byte(`{"id":1}`), []byte(`{}`), []byte(`{"age":32}`)),
		)
		assert.Equal(t, []byte(`[{"id":1},{"age":32}]`),
			ConcatJSON([]byte(`[{"id":1}]`), []byte(`[]`), []byte(`[{"age":32}]`)),
		)
		assert.Equal(t, []byte(`{"id":1,"name":"Rachel"}`),
			ConcatJSON([]byte(`{"id":1}`), []byte(`{"name":"Rachel"}`), []byte(`{}`)),
		)
		assert.Equal(t, []byte(`[{"id":1},{"name":"Rachel"}]`),
			ConcatJSON([]byte(`[{"id":1}]`), []byte(`[{"name":"Rachel"}]`), []byte(`[]`)),
		)
	})

	t.Run("should concat empty objects and arrays", func(t *testing.T) {
		assert.Equal(t, []byte(`{}`),
			ConcatJSON([]byte(`{}`), []byte(`{}`), []byte(`{}`)),
		)
		assert.Equal(t, []byte(`[]`),
			ConcatJSON([]byte(`[]`), []byte(`[]`), []byte(`[]`)),
		)
	})

	t.Run("should concat objects with null", func(t *testing.T) {
		assert.Equal(t, []byte(nil),
			ConcatJSON([]byte(nil)),
		)
		assert.Equal(t, []byte(nil),
			ConcatJSON([]byte(`null`)),
		)
		assert.Equal(t, []byte(nil),
			ConcatJSON([]byte(nil), []byte(`null`)),
		)
		assert.Equal(t, []byte(`{"id":null}`),
			ConcatJSON([]byte(`{"id":null}`), []byte(`null`)),
		)
		assert.Equal(t, []byte(`{"id":null,"name":"Rachel"}`),
			ConcatJSON([]byte(`{"id":null}`), []byte(`null`), []byte(`{"name":"Rachel"}`)),
		)
	})

	t.Run("should concat arrays with null", func(t *testing.T) {
		assert.Equal(t, []byte(`[{"id":1}]`),
			ConcatJSON([]byte(`[{"id":1}]`), []byte(nil)),
		)
	})

	t.Run("should NOT concat non-containers", func(t *testing.T) {
		assert.Equal(t, []byte(nil),
			ConcatJSON([]byte(`"a"`), []byte(`1`)),
		)
	})
}
