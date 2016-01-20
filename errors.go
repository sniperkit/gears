package gears

// StatusError is an interface used for handling
// http errors with proper statuses
type StatusError interface {
	Error() string
	Status() int
}

// JSONError is an interface used for handling
// errors as JSON encoded bytes.
type JSONError interface {
	Error() string
	MarshalJSON() ([]byte, error)
}
