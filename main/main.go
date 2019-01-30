package main

import (
	"golisp/interpreter"
	"golisp/reader"
)

func main() {
	tokentree := *reader.ReadInput()
	traverser := interpreter.Traverser{Tree: tokentree}
	traverser.Traverse()
	/*test := GlispError.GLispError{}
	testFuture := GlispError.Future{Contents: test}
	fmt.Println(testFuture.IsGLispError())
	*/
}
