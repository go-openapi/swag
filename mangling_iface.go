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

var GoNamePrefixFunc = &mangling.GoNamePrefixFunc

// AddInitialisms add additional initialisms
//
// See [mangling.AddInitialisms].
func AddInitialisms(word ...string) { mangling.AddInitialisms(word...) }

// Camelize an uppercased word
//
// See [mangling.Camelize].
func Camelize(word string) string { return mangling.Camelize(word) }

// ToFileName lowercases and underscores a go type name
//
// See [mangling.ToFileName].
func ToFileName(name string) string { return mangling.ToFileName(name) }

// ToCommandName lowercases and underscores a go type name
//
// See [mangling.ToCommandName].
func ToCommandName(name string) string { return mangling.ToCommandName(name) }

// ToHumanNameLower represents a code name as a human series of words
//
// See [mangling.ToHumanNameLower].
func ToHumanNameLower(name string) string { return mangling.ToHumanNameLower(name) }

// ToHumanNameTitle represents a code name as a human series of words with the first letters titleized
//
// See [mangling.ToHumanNameTitle].
func ToHumanNameTitle(name string) string { return mangling.ToHumanNameTitle(name) }

// ToJSONName camelcases a name which can be underscored or pascal cased
//
// See [mangling.ToJSONName].
func ToJSONName(name string) string { return mangling.ToJSONName(name) }

// ToVarName camelcases a name which can be underscored or pascal cased
//
// See [mangling.ToVarName].
func ToVarName(name string) string { return mangling.ToVarName(name) }

// ToGoName translates a swagger name which can be underscored or camel cased to a name that golint likes
//
// See [mangling.ToGoName].
func ToGoName(name string) string { return mangling.ToGoName(name) }
