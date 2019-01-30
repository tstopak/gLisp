package core

import (
	"golisp/GlispError"
)

type dispatcher struct {
	ns GLispNamespace
}
type Dispatcher interface {
	Call(callForm []string) (result string)
}

func NewDispatcher() (disp Dispatcher) {
	gLispNS := NewGLispNamespace()
	disp = dispatcher{gLispNS}
	return
}

func (disp dispatcher) Call(callForm []string) (result string) {
	funcName := callForm[0]
	args := callForm[1:]
	value := disp.ns.Invoke(funcName, args)
	if value.IsGLispError() {
		thisError := value.Contents.(GlispError.GLispError)
		thisError.PrintError()
		panic("")
	} else {
		result = value.Contents.(string)
	}
	return
}
