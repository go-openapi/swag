// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonutils_test

import (
	"fmt"

	"github.com/go-openapi/swag/jsonutils"
)

func ExampleReadJSON() {
	const jazon = `{"a": 1,"b": "x"}`
	var value any

	if err := jsonutils.ReadJSON([]byte(jazon), &value); err != nil {
		panic(err)
	}

	reconstructed, err := jsonutils.WriteJSON(value)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(reconstructed))

	// Output:
	// {"a":1,"b":"x"}
}

type A struct {
	A string
	B []int
	C struct {
		X string
		Y bool
	}
}

func ExampleFromDynamicJSON_struct() {
	source := A{
		A: "x",
		B: []int{0, 1},
		C: struct {
			X string
			Y bool
		}{X: "y", Y: true},
	}

	var target any

	if err := jsonutils.FromDynamicJSON(source, &target); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", target)

	// Output:
	// map[string]interface {}{"A":"x", "B":[]interface {}{0, 1}, "C":map[string]interface {}{"X":"y", "Y":true}}
}

func ExampleFromDynamicJSON_orderedmap() {
	source := jsonutils.JSONMapSlice{
		{Key: "A", Value: "x"},
		{Key: "B", Value: []int{0, 1}},
		{Key: "C", Value: jsonutils.JSONMapSlice{
			{Key: "X", Value: "y"},
			{Key: "Y", Value: true},
		}},
	}

	var target jsonutils.JSONMapSlice

	if err := jsonutils.FromDynamicJSON(source, &target); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", target)

	// Output:
	// jsonutils.JSONMapSlice{jsonutils.JSONMapItem{Key:"A", Value:"x"}, jsonutils.JSONMapItem{Key:"B", Value:[]interface {}{0, 1}}, jsonutils.JSONMapItem{Key:"C", Value:json.MapSlice{json.MapItem{Key:"X", Value:"y"}, json.MapItem{Key:"Y", Value:true}}}}
}

func ExampleConcatJSON_objects() {
	blob1 := []byte(`{"a": 1,"b": "x"}`)
	blob2 := []byte(`{"c": 1,"d": "z"}`)
	blob3 := []byte(`{"e": "y"}`)

	// blobs are not merged: if common keys appear, duplicates will be created
	reunited := jsonutils.ConcatJSON(blob1, blob2, blob3)

	fmt.Println(string(reunited))

	// Output:
	// {"a": 1,"b": "x","c": 1,"d": "z","e": "y"}
}

func ExampleConcatJSON_arrays() {
	blob1 := []byte(`["a","b","x"]`)
	blob2 := []byte(`["c","d"]`)
	blob3 := []byte(`["e","y"]`)

	// blobs are not merged: if common keys appear, duplicates will be created
	reunited := jsonutils.ConcatJSON(blob1, blob2, blob3)

	fmt.Println(string(reunited))

	// Output:
	// ["a","b","x","c","d","e","y"]
}

func ExampleJSONMapSlice() {
	const jazon = `{"a": 1,"c": "x", "b": 2}`
	var value jsonutils.JSONMapSlice

	if err := value.UnmarshalJSON([]byte(jazon)); err != nil {
		panic(err)
	}

	reconstructed, err := value.MarshalJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(reconstructed))
	fmt.Printf("%#v\n", value)

	// Output:
	// {"a":1,"c":"x","b":2}
	// jsonutils.JSONMapSlice{jsonutils.JSONMapItem{Key:"a", Value:1}, jsonutils.JSONMapItem{Key:"c", Value:"x"}, jsonutils.JSONMapItem{Key:"b", Value:2}}
}
