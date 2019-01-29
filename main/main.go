package main

import (
	"fmt"
	"golisp/GlispError"
)

func main() {
	/*tokentree := *reader.ReadInput()
	traverser := interpreter.Traverser{Tree: tokentree}
	traverser.Traverse()
	*/
	test := GlispError.GLispError{}
	testFuture := GlispError.Future{Contents: test}
	fmt.Println(testFuture.IsGLispError())
}
