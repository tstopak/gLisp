package interpreter

import (
	"fmt"
	"golisp/reader"
)

type Traverser struct {
	Tree     reader.TokenTree
	dispatch Dispatcher
}

func (traverser *Traverser) Initialize() {
	traverser.dispatch = NewDispatcher(traverser)
}

func (traverser Traverser) Traverse() {
	rootNode := *traverser.Tree.Root
	fmt.Println(traverser.Interpret(rootNode))
}

func (traverser Traverser) Interpret(token reader.Token) (result string) {
	//Assume that the token is a list
	isList := true
	//Verify the assumption
	if len(token.Children) == 0 {
		isList = false
	}
	var isQuote bool
	if isList && string(token.Value) == "'(" {
		isQuote = true
	}
	if isQuote {
		result = "'( "
		if len(token.Children) != 0 {
			for _, child := range token.Children {
				addedQuote := false
				if child.Value == "(" {
					child.Value = "'("
					addedQuote = true
				}
				tempResult := traverser.Interpret(*child) + " "
				if addedQuote {
					tempResult = tempResult[1:]
				}
				result += tempResult
			}
			result += ")"
		} else {
			result = "'()"
		}
		return result
	} else if isList {
		isDefun := false
		if string(token.Children[0].Value) == "defun" {
			isDefun = true
		}
		if isDefun {
			name := token.Children[1].Value
			params := token.Children[2]
			body := token.Children[3]
			traverser.dispatch.Defun(name, params, body)
			return "true"
		}
		isIf := false
		if string(token.Children[0].Value) == "if" {
			isIf = true
		}
		if isIf {
			testResult := traverser.Interpret(*token.Children[1])
			if testResult == "true" {
				return traverser.Interpret(*token.Children[2])
			} else {
				return traverser.Interpret(*token.Children[3])
			}
		}
		isFunc := false
		//Create a variable to hold a list in the case this is a function call
		var funcCall []string = nil
		// Check if this list is a function call
		if string(token.Children[0].Value) != "" {
			isFunc = true
			// In the case it is instantiate the array
			funcCall = make([]string, 0, len(token.Children))
		}
		// For every child
		for _, child := range token.Children {
			// interpret this child recursively and save result
			var value string
			value = traverser.Interpret(*child)

			//If this list is a function call append all evaluated children to an array
			if isFunc {
				funcCall = append(funcCall, value)
			}
		}
		//If the list we just evaluated all of the children for was a function
		if isFunc {
			fmt.Print("Call trace: ")
			fmt.Println(funcCall)
			//Pass the constructed array to the dispatcher to handle the function call
			result = traverser.dispatch.Call(funcCall)
		}
	} else {
		//This was not a list and we can just pass back the literal value of the token
		if token.Value == "(" {
			return "()"
		}
		return token.Value
	}
	// return the result of the list evaluation
	return
}
