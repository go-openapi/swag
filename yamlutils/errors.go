package yamlutils

type yamlError string

const (
	// ErrYAML is an error raised by YAML utilities
	ErrYAML yamlError = "yaml error"
)

func (e yamlError) Error() string {
	return string(e)
}
