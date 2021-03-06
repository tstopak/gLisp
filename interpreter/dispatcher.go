package interpreter

import (
	"glisp/GlispError"
	"glisp/reader"
)

type dispatcher struct {
	ns *GLispNamespace
}
type Dispatcher interface {
	Call(callForm []string) (result string)
	Defun(name string, params reader.Token, body reader.Token)
	Defparameter(name string, value *reader.Token)
	ResolveVar(name string) string
}

func NewDispatcher(traverser *Traverser) (disp Dispatcher) {
	gLispNS := NewGLispNamespace(traverser)
	disp = dispatcher{&gLispNS}
	return
}

func (disp dispatcher) Call(callForm []string) (result string) {
	funcName := callForm[0]
	args := callForm[1:]
	ns := *disp.ns
	value := ns.Invoke(funcName, args)
	if value.IsGLispError() {
		thisError := value.Contents.(GlispError.GLispError)
		thisError.PrintError()
		panic("")
	} else {
		result = value.Contents.(string)
	}
	return
}

func (disp dispatcher) Defun(name string, params reader.Token, body reader.Token) {
	ns := *disp.ns
	ns.Defun(name, params, body)
}

func (disp dispatcher) Defparameter(name string, value *reader.Token) {
	ns := *disp.ns
	ns.Defparameter(name, value)

}

func (disp dispatcher) ResolveVar(name string) string {
	ns := *disp.ns
	val := ns.ResolveVar(name)
	return val
}
