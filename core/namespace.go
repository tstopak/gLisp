package core

import (
	"golisp/GlispError"
	"strconv"
)

type gLispNamespace struct {
	specialForms map[string]func([]string) string
}
type GLispNamespace interface {
	Invoke(funcName string, args []string) GlispError.Future
}

func NewGLispNamespace() (ns GLispNamespace) {
	localNs := gLispNamespace{
		specialForms: make(map[string]func([]string) string),
	}
	localNs.specialForms["+"] = add

	return localNs
}
func (gLispNs gLispNamespace) Invoke(funcName string, args []string) (result GlispError.Future) {
	if funcCall, exists := gLispNs.specialForms[funcName]; exists {
		result = GlispError.Future{Contents: funcCall(args)}
	} else {
		result = GlispError.Future{Contents: GlispError.GLispError{}}
	}
	return
}
func add(args []string) string {
	var result int64
	for _, arg := range args {
		number, _ := strconv.ParseInt(arg, 10, 64)
		result += number
	}
	return strconv.FormatInt(result, 10)
}
