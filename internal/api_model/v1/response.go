package v1

// Response is a version 1 API response object.
type Response[T any] struct {
	Error   bool
	Message string
	Data    []T
}
