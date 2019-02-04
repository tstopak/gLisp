package main

import (
	"golisp/interpreter"
	"golisp/reader"
)

func main() {
	tokenTree := reader.ReadInput()
	traverser := interpreter.Traverser{Tree: tokenTree}
	traverser.Initialize()
	for {
		traverser.Traverse()
		tokenTree := reader.ReadInput()
		traverser.Tree = tokenTree
	}
}
