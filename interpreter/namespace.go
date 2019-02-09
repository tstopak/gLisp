package interpreter

import (
	"fmt"
	"golisp/GlispError"
	"golisp/core"
	"golisp/reader"
	"strconv"
)

type gLispNamespace struct {
	specialForms        map[string]func([]string) string
	globalFunctionSpace map[string]core.Function
	globalConstantSpace map[string]string
	Traverser           *Traverser
}
type GLispNamespace interface {
	Invoke(funcName string, args []string) GlispError.Future
	Defun(name string, params reader.Token, body reader.Token)
	Defparameter(name string, value *reader.Token)
	ResolveVar(name string) string
}

func NewGLispNamespace(traverser *Traverser) GLispNamespace {
	localNs := gLispNamespace{
		specialForms:        make(map[string]func([]string) string),
		globalFunctionSpace: make(map[string]core.Function),
		globalConstantSpace: make(map[string]string),
		Traverser:           traverser,
	}
	localNs.specialForms["+"] = add
	localNs.specialForms["="] = equals
	localNs.specialForms["-"] = sub
	localNs.specialForms["/"] = div
	localNs.specialForms["*"] = mul
	localNs.specialForms[">"] = gt
	localNs.specialForms["<"] = lt
	localNs.specialForms["!="] = ne
	localNs.specialForms[">="] = gte
	localNs.specialForms["<="] = lte
	localNs.specialForms["and"] = and
	localNs.specialForms["or"] = or
	localNs.specialForms["println"] = println
	return &localNs
}
func (gLispNs gLispNamespace) Invoke(funcName string, args []string) (result GlispError.Future) {
	if funcCall, exists := gLispNs.specialForms[funcName]; exists {
		result = GlispError.Future{Contents: funcCall(args)}
	} else if funcCall, exists := gLispNs.globalFunctionSpace[funcName]; exists {
		fun := deepCopyFunction(funcCall)
		result = GlispError.Future{Contents: gLispNs.evalUserFunc(fun, args)}
	} else {
		gle := GlispError.GLispError{}
		(&gle).SetName("Function Not Found Error")
		(&gle).SetCause("An attempt was made to call a function that doesn't exist")
		result = GlispError.Future{Contents: gle}
	}
	return
}
func (gLispNs gLispNamespace) ResolveVar(name string) string {
	value, exist := gLispNs.globalConstantSpace[name]
	if exist == false {
		return name
	}
	return value
}
func (gLispNs gLispNamespace) evalUserFunc(fun core.Function, args []string) string {
	requiredParams := fun.Param.Children
	body := fun.Body
	for index, param := range requiredParams {
		value := args[index]
		insertParam(body, param.Value, value)
	}
	return gLispNs.Traverser.Interpret(body)

}
func insertParam(token reader.Token, param string, value string) reader.Token {
	for _, child := range token.Children {
		if child.Value == param {
			child.Value = value
		} else if child.Value == "(" {
			insertParam(*child, param, value)
		}
	}
	return token
}
func (gLispNs *gLispNamespace) Defun(name string, params reader.Token, body reader.Token) {
	gLispNs.globalFunctionSpace[name] = core.Function{params, body}
}
func (gLispNs *gLispNamespace) Defparameter(name string, value *reader.Token) {
	constantVal := gLispNs.Traverser.Interpret(*value)
	gLispNs.globalConstantSpace[name] = constantVal
}
func add(args []string) string {
	var result int64
	result = 0
	for _, arg := range args {
		number, _ := strconv.ParseInt(arg, 10, 64)
		result += number
	}
	return strconv.FormatInt(result, 10)
}
func mul(args []string) string {
	var result int64
	result = 1
	for _, arg := range args {
		number, _ := strconv.ParseInt(arg, 10, 64)
		result *= number
	}
	return strconv.FormatInt(result, 10)
}

func sub(args []string) string {
	var result int64
	result, _ = strconv.ParseInt(args[0], 10, 64)
	args = args[1:]
	for _, arg := range args {
		number, _ := strconv.ParseInt(arg, 10, 64)
		result -= number
	}
	return strconv.FormatInt(result, 10)
}

func div(args []string) string {
	var result int64
	result, _ = strconv.ParseInt(args[0], 10, 64)
	args = args[1:]
	for _, arg := range args {
		number, _ := strconv.ParseInt(arg, 10, 64)
		result /= number
	}
	return strconv.FormatInt(result, 10)
}

func equals(args []string) string {
	num1, _ := strconv.ParseInt(args[0], 10, 64)
	num2, _ := strconv.ParseInt(args[1], 10, 64)

	if num1 == num2 {
		return "true"
	}
	return "false"
}

func gt(args []string) string {
	num1, _ := strconv.ParseInt(args[0], 10, 64)
	num2, _ := strconv.ParseInt(args[1], 10, 64)

	if num1 > num2 {
		return "true"
	}
	return "false"
}

func lt(args []string) string {
	num1, _ := strconv.ParseInt(args[0], 10, 64)
	num2, _ := strconv.ParseInt(args[1], 10, 64)

	if num1 < num2 {
		return "true"
	}
	return "false"
}
func gte(args []string) string {
	num1, _ := strconv.ParseInt(args[0], 10, 64)
	num2, _ := strconv.ParseInt(args[1], 10, 64)

	if num1 >= num2 {
		return "true"
	}
	return "false"
}

func lte(args []string) string {
	num1, _ := strconv.ParseInt(args[0], 10, 64)
	num2, _ := strconv.ParseInt(args[1], 10, 64)

	if num1 <= num2 {
		return "true"
	}
	return "false"
}

func ne(args []string) string {
	num1, _ := strconv.ParseInt(args[0], 10, 64)
	num2, _ := strconv.ParseInt(args[1], 10, 64)

	if num1 != num2 {
		return "true"
	}
	return "false"
}

func and(args []string) string {
	for _, child := range args {
		if child == "false" {
			return "false"
		}
	}
	return "true"
}

func or(args []string) string {
	for _, child := range args {
		if child == "true" {
			return "true"
		}
	}
	return "false"
}

func println(args []string) string {
	for _, arg := range args {
		fmt.Println(arg)
	}
	return args[len(args)-1]
}

func deepCopyFunction(fun core.Function) core.Function {
	copyFun := core.Function{}
	copyFun.Param = *deepCopyToken(fun.Param)
	copyFun.Body = *deepCopyToken(fun.Body)
	return copyFun
}

func deepCopyToken(token reader.Token) *reader.Token {
	copiedToken := reader.Token{}
	if len(token.Children) != 0 {
		for _, child := range token.Children {
			copiedToken.Children = append(copiedToken.Children, deepCopyToken(*child))
		}
		copiedToken.Value = token.Value
		return &copiedToken
	} else {
		copiedToken.Value = token.Value
		return &copiedToken
	}
}
