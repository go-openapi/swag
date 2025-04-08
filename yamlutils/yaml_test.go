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

package yamlutils

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
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

func TestJSONToYAML(t *testing.T) {
	t.Run("should convert JSON to YAML", func(t *testing.T) {
		t.Run("with string values", func(t *testing.T) {
			const (
				input    = `{"1":"the int key value","name":"a string value","y":"some value"}`
				expected = `"1": the int key value
name: a string value
y: some value
`
			)

			var data YAMLMapSlice
			require.NoError(t, json.Unmarshal([]byte(input), &data))

			y, err := data.MarshalYAML()
			require.NoError(t, err)
			assert.Equal(t, expected, string(y.([]byte)))
		})

		t.Run("with nested object", func(t *testing.T) {
			const (
				input    = `{"1":"the int key value","name":"a string value","y":"some value","tag":{"name":"tag name"}}`
				expected = `"1": the int key value
name: a string value
y: some value
tag:
    name: tag name
`
			)

			var data YAMLMapSlice
			require.NoError(t, json.Unmarshal([]byte(input), &data))
			ny, err := data.MarshalYAML()
			require.NoError(t, err)
			assert.Equal(t, expected, string(ny.([]byte)))
		})
	})

	t.Run("with complete doc", func(t *testing.T) {
		t.Run("should convert bytes to YAML doc", func(t *testing.T) {
			ydoc, err := BytesToYAMLDoc(fixture2224)
			require.NoError(t, err)

			t.Run("should convert YAML doc to JSON", func(t *testing.T) {
				buf, err := YAMLToJSON(ydoc)
				require.NoError(t, err)

				t.Run("should unmarshal JSON into YAMLMapSlice", func(t *testing.T) {
					var data YAMLMapSlice
					require.NoError(t, json.Unmarshal(buf, &data))

					t.Run("should marshal YAMLMapSlice into the original doc", func(t *testing.T) {
						reconstructed, err := data.MarshalYAML()
						require.NoError(t, err)

						text, ok := reconstructed.([]byte)
						require.True(t, ok)

						assert.YAMLEq(t, string(fixture2224), string(text))
					})
				})
			})
		})
	})
}

func TestJSONToYAMLWithNull(t *testing.T) {
	const (
		jazon    = `{"1":"the int key value","name":null,"y":"some value"}`
		expected = `"1": the int key value
name: null
y: some value
`
	)
	var data YAMLMapSlice
	require.NoError(t, json.Unmarshal([]byte(jazon), &data))
	ny, err := data.MarshalYAML()
	require.NoError(t, err)
	assert.Equal(t, expected, string(ny.([]byte)))
}

func TestMarshalYAML(t *testing.T) {
	t.Run("marshalYAML should be deterministic", func(t *testing.T) {
		const (
			jazon    = `{"1":"x","2":null,"3":{"a":1.1,"b":2.2,"c":3.3}}`
			expected = `"1": x
"2": null
"3":
    a: 1.1
    b: 2.2
    c: 3.3
`
		)
		const iterations = 10
		for n := 0; n < iterations; n++ {
			var data YAMLMapSlice
			require.NoError(t, json.Unmarshal([]byte(jazon), &data))
			ny, err := data.MarshalYAML()
			require.NoError(t, err)
			assert.Equal(t, expected, string(ny.([]byte)))
		}
	})
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
			assert.JSONEq(t, `{"1":"the int key value","name":"a string value","y":"some value"}`, string(d))
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
			assert.JSONEq(t, `{"1":"the int key value","name":"a string value","y":"some value","tag":{"name":"tag name"}}`, string(d))
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
			var lst []interface{}
			lst = append(lst, "hello")

			d, err := YAMLToJSON(lst)
			require.NoError(t, err)
			require.NotNil(t, d)
			assert.JSONEq(t, `["hello"]`, string(d))

			t.Run("should convert object appended to array to JSON", func(t *testing.T) {
				lst = append(lst, data)

				d, err = YAMLToJSON(lst)
				require.NoError(t, err)
				require.NotEmpty(t, d)
				assert.JSONEq(t, `["hello",{"1":"the int key value","name":"a string value","y":"some value"}]`, string(d))
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
				assert.JSONEq(t, `{"description":"object created"}`, string(d))
			})
		})
	})
}

func TestWithYKey(t *testing.T) {
	doc, err := BytesToYAMLDoc([]byte(fixtureWithYKey))
	require.NoError(t, err)

	_, err = YAMLToJSON(doc)
	require.NoError(t, err)

	doc, err = BytesToYAMLDoc([]byte(fixtureWithQuotedYKey))
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
	dm := map[interface{}]interface{}{
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

func mustLoadFixture(name string) []byte {
	const msg = "wrong embedded FS configuration: %w"
	data, err := embeddedFixtures.ReadFile(path.Join("fixtures", name)) // "/" even on windows
	if err != nil {
		panic(fmt.Errorf(msg, err))
	}

	return data
}
