// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package yamlutils_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/swag/yamlutils"
)

func ExampleYAMLToJSON() {
	const doc = `
---
object:
  key: x
  b: true
  n: 1
`

	yml, err := yamlutils.BytesToYAMLDoc([]byte(doc))
	if err != nil {
		panic(err)
	}

	d, err := yamlutils.YAMLToJSON(yml)
	if err != nil {
		panic(err)
	}

	jazon, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jazon))
	// Output:
	// {
	//   "object": {
	//     "key": "x",
	//     "b": true,
	//     "n": 1
	//   }
	// }
}

func ExampleYAMLMapSlice() {
	const doc = `
---
object:
  key: x
  b: true
  n: 1
`

	ydoc, err := yamlutils.BytesToYAMLDoc([]byte(doc))
	if err != nil {
		panic(err)
	}

	jazon, err := yamlutils.YAMLToJSON(ydoc)
	if err != nil {
		panic(err)
	}

	var data yamlutils.YAMLMapSlice
	err = json.Unmarshal(jazon, &data)
	if err != nil {
		panic(err)
	}

	// reconstruct the initial YAML document, preserving the order of keys
	// (but not YAML specifics such as anchors, comments, ...).
	reconstructed, err := data.MarshalYAML()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(reconstructed.([]byte)))
	// Output:
	// object:
	//     key: x
	//     b: true
	//     n: 1
}
