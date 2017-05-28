package osenv

import "reflect"

type invalidValueError struct {
	Type reflect.Type
}

func (e *invalidValueError) Error() string {
	return e.Type.String()
}
