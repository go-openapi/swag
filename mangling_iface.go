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

import "github.com/go-openapi/swag/mangling"

// GoNamePrefixFunc sets an optional rule to prefix go names
// which do not start with a letter.
//
// Deprecated: use [mangling.GoNamePrefixFunc] instead.
var GoNamePrefixFunc = mangling.GoNamePrefixFunc

// AddInitialisms adds additional initialisms.
//
// Deprecated: use [mangling.AddInitialisms] instead.
func AddInitialisms(word ...string) { mangling.AddInitialisms(word...) }

// Camelize an uppercased word.
//
// Deprecated: use [mangling.Camelize] instead.
func Camelize(word string) string { return mangling.Camelize(word) }

// ToFileName lowercases and underscores a go type name.
//
// Deprecated: use [mangling.ToFileName] instead.
func ToFileName(name string) string { return mangling.ToFileName(name) }

// ToCommandName lowercases and underscores a go type name.
//
// Deprecated: use [mangling.ToCommandName] instead.
func ToCommandName(name string) string { return mangling.ToCommandName(name) }

// ToHumanNameLower represents a code name as a human series of words.
//
// Deprecated: use [mangling.ToHumanNameLower] instead.
func ToHumanNameLower(name string) string { return mangling.ToHumanNameLower(name) }

// ToHumanNameTitle represents a code name as a human series of words with the first letters titleized.
//
// Deprecated: use [mangling.ToHumanNameTitle] instead.
func ToHumanNameTitle(name string) string { return mangling.ToHumanNameTitle(name) }

// ToJSONName camelcases a name which can be underscored or pascal cased.
//
// Deprecated: use [mangling.ToJSONName] instead.
func ToJSONName(name string) string { return mangling.ToJSONName(name) }

// ToVarName camelcases a name which can be underscored or pascal cased.
//
// Deprecated: use [mangling.ToVarName] instead.
func ToVarName(name string) string { return mangling.ToVarName(name) }

// ToGoName translates a swagger name which can be underscored or camel cased to a name that golint likes.
//
// Deprecated: use [mangling.ToGoName] instead.
func ToGoName(name string) string { return mangling.ToGoName(name) }
