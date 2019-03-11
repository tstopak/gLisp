package main

import (
	"glisp/interpreter"
	"glisp/reader"
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
