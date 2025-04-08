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
	"encoding/json"
	"testing"

	"github.com/mailru/easyjson/jlexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONMapSlice(t *testing.T) {
	t.Run("should unmarshal and marshal MapSlice", func(t *testing.T) {
		t.Run("with object", func(t *testing.T) {
			const sd = `{"1":"the int key value","name":"a string value","y":"some value"}`
			var data JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(sd), &data))

			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			assert.JSONEq(t, sd, string(jazon))
		})

		t.Run("with nested object", func(t *testing.T) {
			const sd = `
		{
		  "1":"the int key value",
		  "name":"a string value",
		  "y":{
		    "a":"some value",
		    "b":[
		     {"x":1,"y":2},
		     {"z":4,"w":5}
		    ]
		  }
		}
		`
			var data JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(sd), &data))

			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			assert.JSONEq(t, sd, string(jazon))
		})

		t.Run("with nested array", func(t *testing.T) {
			const sd = `
	[
	  {
			"1":"the int key value",
	    "name":"a string value"
		},
		{
	    "y":{
	      "a":"some value",
	      "b": [
	         {"x":1,"y":2},
	         {"z":4,"w":5}
	      ],
			  "c": false,
			  "d": null
	    },
	    "z": true
	  },
		{
			"v": [true, "string", 10.35]
		}
	]
	`
			var data []JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(sd), &data))

			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			assert.JSONEq(t, sd, string(jazon))
		})

		t.Run("with empty array", func(t *testing.T) {
			const sd = `{"a":[]}`
			var data JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(sd), &data))

			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			assert.JSONEq(t, sd, string(jazon))
		})

		t.Run("with empty object", func(t *testing.T) {
			const sd = `{}`
			var data JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(sd), &data))

			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			assert.JSONEq(t, sd, string(jazon))
		})

		t.Run("with null value", func(t *testing.T) {
			const sd = `null`
			var data JSONMapSlice
			require.NoError(t, json.Unmarshal([]byte(sd), &data))
			assert.Nil(t, data)

			jazon, err := json.Marshal(data)
			require.NoError(t, err)

			assert.JSONEq(t, sd, string(jazon))
		})
	})

	t.Run("should keep the order of keys", func(t *testing.T) {
		const sd = `{"a":1,"b":2,"c":3,"d":4}`
		var data JSONMapSlice
		require.NoError(t, json.Unmarshal([]byte(sd), &data))
		jazon, err := json.Marshal(data)
		require.NoError(t, err)

		require.Equal(t, sd, string(jazon)) // specifically check the same order, not JSONEq()

		t.Run("should Read/Write JSON using easyJSON", func(t *testing.T) {
			var obj interface{}
			require.NoError(t, FromDynamicJSON(data, &obj))

			asMap, ok := obj.(map[string]interface{})
			require.True(t, ok)
			assert.Len(t, asMap, 4) // 3 fields in struct

			var target JSONMapSlice
			require.NoError(t, FromDynamicJSON(obj, &target))
			// the order of keys may be altered, since the intermediary representation is a map[string]interface{}
		})
	})

	t.Run("UnmarshalEasyJSON with error cases", func(t *testing.T) {
		// test directly this endpoint, as the json standard library
		// performs a preventive early check for well-formed JSON.
		t.Run("on invalid token (1)", func(t *testing.T) {
			const sd = `{"a":|,"b":2,"c":3,"d":4}`
			var data JSONMapSlice
			require.Error(t, json.Unmarshal([]byte(sd), &data))
		})
		t.Run("on invalid token (2)", func(t *testing.T) {
			const sd = `{"a":{ai+b,"b":2,"c":3,"d":4}`
			var data JSONMapSlice
			require.Error(t, json.Unmarshal([]byte(sd), &data))
		})
		t.Run("on invalid token (3)", func(t *testing.T) {
			const sd = `{"a":[ai+b,"b":2,"c":3,"d":4}`
			data := make(JSONMapSlice, 0)
			l := jlexer.Lexer{Data: []byte(sd)}
			data.UnmarshalEasyJSON(&l)
			require.Error(t, l.Error())
		})
		t.Run("on invalid delimiter (1)", func(t *testing.T) {
			const sd = `{"a":1`
			data := make(JSONMapSlice, 0)
			l := jlexer.Lexer{Data: []byte(sd)}
			data.UnmarshalEasyJSON(&l)
			require.Error(t, l.Error())
		})
		t.Run("on invalid delimiter (2)", func(t *testing.T) {
			const sd = `{"a":[1}`
			data := make(JSONMapSlice, 0)
			l := jlexer.Lexer{Data: []byte(sd)}
			data.UnmarshalEasyJSON(&l)
			require.Error(t, l.Error())
		})
		t.Run("on invalid delimiter (3)", func(t *testing.T) {
			const sd = `{"a":[1,]}`
			data := make(JSONMapSlice, 0)
			l := jlexer.Lexer{Data: []byte(sd)}
			data.UnmarshalEasyJSON(&l)
			require.Error(t, l.Error())
		})
		t.Run("on invalid delimiter (4)", func(t *testing.T) {
			const sd = `{"a":[1],}`
			data := make(JSONMapSlice, 0)
			l := jlexer.Lexer{Data: []byte(sd)}
			data.UnmarshalEasyJSON(&l)
			require.Error(t, l.Error())
		})
		t.Run("on invalid delimiter (4)", func(t *testing.T) {
			const sd = `{"a":{"b":1}`
			data := make(JSONMapSlice, 0)
			l := jlexer.Lexer{Data: []byte(sd)}
			data.UnmarshalEasyJSON(&l)
			require.Error(t, l.Error())
		})
	})
}
