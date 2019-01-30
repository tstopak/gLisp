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
	fmt.Println(traverser.interpret(rootNode))
}

func (traverser Traverser) interpret(token reader.Token) (result string) {
	//Assume that the token is a list
	isList := true
	//Verify the assumption
	if token.Children == nil {
		isList = false
	}
	//If the token is a list
	if isList {
		isFunc := false
		//Create a variable to hold a list in the case this is a function call
		var funcCall []string = nil
		// Check if this list is a function call
		if string(token.Children[0].Value[0]) != "'" {
			isFunc = true
			// In the case it is instantiate the array
			funcCall = make([]string, 0, len(token.Children))
		}
		// For every child
		for _, child := range token.Children {
			// interpret this child recursively and save result
			value := traverser.interpret(*child)
			//If this list is a function call append all evaluated children to an array
			if isFunc {
				funcCall = append(funcCall, value)
			}
		}
		//If the list we just evaluated all of the children for was a function
		if isFunc {
			fmt.Println(funcCall)
			//Pass the constructed array to the dispatcher to handle the function call
			result = traverser.dispatch.Call(funcCall)
		} else {
			// Create a list as the value of the token
			result += "("
			for _, child := range token.Children {
				// Still interpret any sub-values
				val := traverser.interpret(*child)
				//Reconstruct the list without making a top level function call
				result += val + " "
			}
			result += ")"
		}
	} else {
		//This was not a list and we can just pass back the literal value of the token
		return token.Value
	}
	// return the result of the list evaluation
	return
}
