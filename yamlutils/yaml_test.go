// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package yamlutils

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	_ "github.com/go-openapi/testify/enable/yaml/v2" // enable YAMLEq in testify
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	yaml "go.yaml.in/yaml/v3"
)

// embedded test files

//go:embed fixtures/*.yaml
var embeddedFixtures embed.FS

var fixtureSpecTags, fixture2224, fixtureWithQuotedYKey, fixtureWithYKey []byte

func TestMain(m *testing.M) {
	fixtureSpecTags = mustLoadFixture("fixture_spec_tags.yaml")
	fixture2224 = mustLoadFixture("fixture_2224.yaml")
	fixtureWithQuotedYKey = mustLoadFixture("fixture_with_quoted.yaml")
	fixtureWithYKey = mustLoadFixture("fixture_with_ykey.yaml")

	os.Exit(m.Run())
}

func TestYAMLToJSON(t *testing.T) {
	const sd = `---
1: the int key value
name: a string value
'y': some value
`

	t.Run("with initial YAML doc", func(t *testing.T) {
		var data yaml.Node
		_ = yaml.Unmarshal([]byte(sd), &data)

		t.Run("should convert YAML doc to JSON", func(t *testing.T) {
			d, err := YAMLToJSON(data)
			require.NoError(t, err)
			require.NotNil(t, d)
			expected := []byte(`{"1":"the int key value","name":"a string value","y":"some value"}`)
			assert.JSONEqBytes(t, expected, d)
		})

		t.Run("should NOT convert appended YAML doc to JSON", func(t *testing.T) {
			ns := []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: "true",
					Tag:   "!!bool",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: "the bool value",
					Tag:   "!!str",
				},
			}
			data.Content[0].Content = append(data.Content[0].Content, ns...)

			d, err := YAMLToJSON(data)
			require.Error(t, err)
			require.Nil(t, d)
			require.ErrorContains(t, err, "is not supported as map key")
		})
	})

	t.Run("with initial YAML doc", func(t *testing.T) {
		var data yaml.Node
		_ = yaml.Unmarshal([]byte(sd), &data)

		tag := []*yaml.Node{
			{
				Kind:  yaml.ScalarNode,
				Value: "tag",
				Tag:   "!!str",
			},
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Value: "name",
						Tag:   "!!str",
					},
					{
						Kind:  yaml.ScalarNode,
						Value: "tag name",
						Tag:   "!!str",
					},
				},
			},
		}
		data.Content[0].Content = append(data.Content[0].Content, tag...)

		t.Run("should convert appended YAML doc to JSON", func(t *testing.T) {
			d, err := YAMLToJSON(data)
			require.NoError(t, err)
			expected := []byte(`{"1":"the int key value","name":"a string value","y":"some value","tag":{"name":"tag name"}}`)
			assert.JSONEqBytes(t, expected, d)
		})

		t.Run("should NOT convert appended YAML doc to JSON: key cannot be a bool", func(t *testing.T) {
			tag[1].Content = []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: "true",
					Tag:   "!!bool",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: "the bool tag name",
					Tag:   "!!str",
				},
			}

			d, err := YAMLToJSON(data)
			require.Error(t, err)
			require.Nil(t, d)
			require.ErrorContains(t, err, "is not supported as map key")
		})
	})

	t.Run("with initial YAML doc", func(t *testing.T) {
		var data yaml.Node
		_ = yaml.Unmarshal([]byte(sd), &data)

		t.Run("should convert any array to JSON", func(t *testing.T) {
			var lst []any
			lst = append(lst, "hello")

			d, err := YAMLToJSON(lst)
			require.NoError(t, err)
			require.NotNil(t, d)
			assert.JSONEqBytes(t, []byte(`["hello"]`), d)

			t.Run("should convert object appended to array to JSON", func(t *testing.T) {
				lst = append(lst, data)

				d, err = YAMLToJSON(lst)
				require.NoError(t, err)
				require.NotEmpty(t, d)
				assert.JSONEqBytes(t, []byte(`["hello",{"1":"the int key value","name":"a string value","y":"some value"}]`), d)
			})
		})
	})

	t.Run("from YAML bytes", func(t *testing.T) {
		t.Run("root document is an object. Should not convert", func(t *testing.T) {
			_, err := BytesToYAMLDoc([]byte("- name: hello\n"))
			require.Error(t, err)
		})

		t.Run("document is invalid YAML. Should not convert", func(t *testing.T) {
			_, err := BytesToYAMLDoc([]byte("name:\tgreetings: hello\n"))
			require.Error(t, err)
		})

		t.Run("root document is an object. Should convert", func(t *testing.T) {
			dd, err := BytesToYAMLDoc([]byte("description: 'object created'\n"))
			require.NoError(t, err)

			t.Run("should convert YAML object to JSON", func(t *testing.T) {
				d, err := YAMLToJSON(dd)
				require.NoError(t, err)
				assert.JSONEqBytes(t, []byte(`{"description":"object created"}`), d)
			})
		})
	})
}

func TestWithYKey(t *testing.T) {
	doc, err := BytesToYAMLDoc(fixtureWithYKey)
	require.NoError(t, err)

	_, err = YAMLToJSON(doc)
	require.NoError(t, err)

	doc, err = BytesToYAMLDoc(fixtureWithQuotedYKey)
	require.NoError(t, err)
	jsond, err := YAMLToJSON(doc)
	require.NoError(t, err)

	var yt struct {
		Definitions struct {
			Viewbox struct {
				Properties struct {
					Y struct {
						Type string `json:"type"`
					} `json:"y"`
				} `json:"properties"`
			} `json:"viewbox"`
		} `json:"definitions"`
	}
	require.NoError(t, json.Unmarshal(jsond, &yt))
	assert.Equal(t, "integer", yt.Definitions.Viewbox.Properties.Y.Type)
}

func TestMapKeyTypes(t *testing.T) {
	dm := map[any]any{
		"12345":             "string",
		12345:               "int",
		int8(1):             "int8",
		int16(12345):        "int16",
		int32(12345678):     "int32",
		int64(12345678910):  "int64",
		uint(12345):         "uint",
		uint8(1):            "uint8",
		uint16(12345):       "uint16",
		uint32(12345678):    "uint32",
		uint64(12345678910): "uint64",
	}
	_, err := YAMLToJSON(dm)
	require.NoError(t, err)
}

func TestYAMLTags(t *testing.T) {
	t.Run("should marshal as a YAML doc", func(t *testing.T) {
		doc, err := BytesToYAMLDoc(fixtureSpecTags)
		require.NoError(t, err)

		t.Run("doc should marshal as the original doc", func(t *testing.T) {
			text, err := yaml.Marshal(doc)
			require.NoError(t, err)
			assert.YAMLEq(t, string(fixtureSpecTags), string(text))
		})

		t.Run("doc should marshal to JSON", func(t *testing.T) {
			jazon, err := YAMLToJSON(doc)
			require.NoError(t, err)

			t.Run("json should unmarshal to YAMLMapSlice", func(t *testing.T) {
				var data YAMLMapSlice
				require.NoError(t, json.Unmarshal(jazon, &data))

				t.Run("YAMLMapSlice should marshal to YAML bytes", func(t *testing.T) {
					text, err := data.MarshalYAML()
					require.NoError(t, err)

					buf, ok := text.([]byte)
					require.True(t, ok)

					// standard YAML used by [assert.YAMLEq] interprets YAML timestamp as [time.Time],
					// but in our context, we use string
					neutralizeTimestamp := strings.ReplaceAll(string(fixtureSpecTags), "default:", "default: !!str ")
					assert.YAMLEq(t, neutralizeTimestamp, string(buf))
				})
			})
		})
	})
}

func TestYAMLEdgeCases(t *testing.T) {
	t.Run("should never happen because never called in the context of arrays", func(t *testing.T) {
		_, err := yamlDocument(&yaml.Node{
			Content: []*yaml.Node{
				{},
				{},
			},
		})
		require.Error(t, err)
	})

	t.Run("should never happen unless the document representation is corrupted", func(t *testing.T) {
		_, err := yamlSequence(&yaml.Node{
			Content: []*yaml.Node{
				{
					Kind: yaml.Kind(99), // illegal kind
				},
			},
		})
		require.Error(t, err)
	})

	t.Run("should never happen unless the document cannot be marshaled", func(t *testing.T) {
		invalidType := func() {}
		_, err := format(invalidType)
		require.Error(t, err)

		_, err = transformData([]any{
			map[any]any{
				complex128(0): struct{}{},
			},
		})
		require.Error(t, err)
	})
}

func mustLoadFixture(name string) []byte {
	const msg = "wrong embedded FS configuration: %w"
	data, err := embeddedFixtures.ReadFile(path.Join("fixtures", name)) // "/" even on windows
	if err != nil {
		panic(fmt.Errorf(msg, err))
	}

	return data
}
