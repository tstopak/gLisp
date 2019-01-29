package interpreter

import (
	"fmt"
	"golisp/core"
	"golisp/reader"
)

type Traverser struct {
	Tree     reader.TokenTree
	dispatch core.Dispatcher
}

func (traverser Traverser) Traverse() {
	rootNode := *traverser.Tree.Root
	traverser.dispatch = core.NewDispatcher()
	traverser.interpret(rootNode)
}

func (traverser Traverser) interpret(token reader.Token) (result string) {
	isList := true
	if token.Children == nil {
		isList = false
	}
	if isList {
		isFunc := false
		var funcCall []string = nil
		if string(token.Children[0].Value[0]) != "'" {
			isFunc = true
			funcCall = make([]string, 0, len(token.Children))
		}
		for _, child := range token.Children {
			value := traverser.interpret(*child)
			if isFunc {
				funcCall = append(funcCall, value)
			}
		}
		if isFunc {
			fmt.Println(funcCall)
			traverser.dispatch.Call(funcCall)
		}
	} else {
		return token.Value
	}
	return "Traversed"
}
