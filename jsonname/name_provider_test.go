// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package jsonname

import (
	"reflect"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

type testNameStruct struct {
	Name       string `json:"name"`
	NotTheSame int64  `json:"plain"`
	Ignored    string `json:"-"`
}

func TestNameProvider(t *testing.T) {
	provider := NewNameProvider()

	var obj = testNameStruct{}

	nm, ok := provider.GetGoName(obj, "name")
	assert.TrueT(t, ok)
	assert.EqualT(t, "Name", nm)

	nm, ok = provider.GetGoName(obj, "plain")
	assert.TrueT(t, ok)
	assert.EqualT(t, "NotTheSame", nm)

	nm, ok = provider.GetGoName(obj, "doesNotExist")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetGoName(obj, "ignored")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	tpe := reflect.TypeOf(obj)
	nm, ok = provider.GetGoNameForType(tpe, "name")
	assert.TrueT(t, ok)
	assert.EqualT(t, "Name", nm)

	nm, ok = provider.GetGoNameForType(tpe, "plain")
	assert.TrueT(t, ok)
	assert.EqualT(t, "NotTheSame", nm)

	nm, ok = provider.GetGoNameForType(tpe, "doesNotExist")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetGoNameForType(tpe, "ignored")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	ptr := &obj
	nm, ok = provider.GetGoName(ptr, "name")
	assert.TrueT(t, ok)
	assert.EqualT(t, "Name", nm)

	nm, ok = provider.GetGoName(ptr, "plain")
	assert.TrueT(t, ok)
	assert.EqualT(t, "NotTheSame", nm)

	nm, ok = provider.GetGoName(ptr, "doesNotExist")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetGoName(ptr, "ignored")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetJSONName(obj, "Name")
	assert.TrueT(t, ok)
	assert.EqualT(t, "name", nm)

	nm, ok = provider.GetJSONName(obj, "NotTheSame")
	assert.TrueT(t, ok)
	assert.EqualT(t, "plain", nm)

	nm, ok = provider.GetJSONName(obj, "DoesNotExist")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetJSONName(obj, "Ignored")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetJSONNameForType(tpe, "Name")
	assert.TrueT(t, ok)
	assert.EqualT(t, "name", nm)

	nm, ok = provider.GetJSONNameForType(tpe, "NotTheSame")
	assert.TrueT(t, ok)
	assert.EqualT(t, "plain", nm)

	nm, ok = provider.GetJSONNameForType(tpe, "doesNotExist")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetJSONNameForType(tpe, "Ignored")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetJSONName(ptr, "Name")
	assert.TrueT(t, ok)
	assert.EqualT(t, "name", nm)

	nm, ok = provider.GetJSONName(ptr, "NotTheSame")
	assert.TrueT(t, ok)
	assert.EqualT(t, "plain", nm)

	nm, ok = provider.GetJSONName(ptr, "doesNotExist")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nm, ok = provider.GetJSONName(ptr, "Ignored")
	assert.FalseT(t, ok)
	assert.Empty(t, nm)

	nms := provider.GetJSONNames(ptr)
	assert.Len(t, nms, 2)

	assert.Len(t, provider.index, 1)
}
