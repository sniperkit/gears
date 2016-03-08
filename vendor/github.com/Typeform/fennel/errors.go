package fennel

import "fmt"

type MultiError interface {
	error
	Errors() []error
}

type BatchGathererError struct {
	errors []error
}

func NewBatchGathererError(errors []error) *BatchGathererError {
	return &BatchGathererError{errors: errors}
}

func (e *BatchGathererError) Error() string {

	s := ""
	for _, err := range e.errors {
		s += fmt.Sprint("%s\n", err.Error())
	}

	if len(s) > 0 {
		s += "\n"
	}

	return s
}

func (e *BatchGathererError) Errors() []error {
	return e.errors
}
