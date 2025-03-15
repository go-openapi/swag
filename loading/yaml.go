package loading

import (
	"encoding/json"
	"path/filepath"

	"github.com/go-openapi/swag/yamlutils"
)

// YAMLMatcher matches yaml
func YAMLMatcher(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".yaml" || ext == ".yml"
}

// YAMLDoc loads a yaml document from either http or a file and converts it to json
func YAMLDoc(path string) (json.RawMessage, error) {
	yamlDoc, err := YAMLData(path)
	if err != nil {
		return nil, err
	}

	data, err := yamlutils.YAMLToJSON(yamlDoc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// YAMLData loads a yaml document from either http or a file
func YAMLData(path string, opts ...Option) (interface{}, error) {
	data, err := LoadFromFileOrHTTP(path, opts...)
	if err != nil {
		return nil, err
	}

	return yamlutils.BytesToYAMLDoc(data)
}
