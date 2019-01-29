package GlispError

import (
	"reflect"
)

type Future struct {
	Contents interface{}
}

func (future Future) IsGLispError() bool {
	if reflect.TypeOf(future.Contents).String() == "GlispError.GLispError" {
		return true
	}
	return false
}
