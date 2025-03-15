package loading

type loadingError string

const (
	// ErrLoader is an error raised by the file loader utility
	ErrLoader loadingError = "loader error"
)

func (e loadingError) Error() string {
	return string(e)
}
