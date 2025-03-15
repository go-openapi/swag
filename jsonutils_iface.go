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

import "github.com/go-openapi/swag/jsonutils"

// DefaultJSONNameProvider is the default cache for types
var DefaultJSONNameProvider = jsonutils.DefaultJSONNameProvider

// WriteJSON writes json data.
//
// See [jsonutils.WriteJSON]
func WriteJSON(data interface{}) ([]byte, error) { return jsonutils.WriteJSON(data) }

// ReadJSON reads json data.
//
// See [jsonutils.ReadJSON]
func ReadJSON(data []byte, value interface{}) error { return jsonutils.ReadJSON(data, value) }

// DynamicJSONToStruct converts an untyped json structure into a struct
func DynamicJSONToStruct(data interface{}, target interface{}) error {
	return jsonutils.DynamicJSONToStruct(data, target)
}

// ConcatJSON concatenates multiple json objects efficiently
func ConcatJSON(blobs ...[]byte) []byte { return jsonutils.ConcatJSON(blobs...) }

// ToDynamicJSON turns an object into a properly JSON typed structure
func ToDynamicJSON(data interface{}) interface{} { return jsonutils.ToDynamicJSON(data) }

// FromDynamicJSON turns an object into a properly JSON typed structure
func FromDynamicJSON(data, target interface{}) error { return jsonutils.FromDynamicJSON(data, target) }

// NameProvider represents an object capable of translating from go property names
// to json property names.
type NameProvider = jsonutils.NameProvider

// NewNameProvider creates a new name provider
func NewNameProvider() *NameProvider { return jsonutils.NewNameProvider() }
