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

// Package swag contains a bunch of helper functions for go-openapi and go-swagger projects.
//
// You may also use it standalone for your projects.
//
// Here is what is inside:
//
//	Sub-package conv:
//	  - convert between value and pointers for builtin types
//	  - convert from string to builtin types (wraps strconv)
//
//	Package swag:
//	- fast json concatenation
//	- search in path
//	- load from file or http
//	- name mangling
//
// This repo has a few dependencies outside of the standard library:
//
//   - YAML utilities depend on [gopkg.in/yaml.v3]
//   - JSON utilities depend on [mailru/easyjson]
package swag
