package interpreter

import (
	"fmt"
	"golisp/reader"
	"unicode"
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
		result = "( "
		if len(token.Children) != 0 {
			for _, child := range token.Children {
				//addedQuote := false
				if child.Value == "(" {
					child.Value = "'("
					//addedQuote = true
				} else if child.Value == ",(" {
					child.Value = "("
				}
				tempResult := traverser.Interpret(*child) + " "
				/*if addedQuote {
					tempResult = tempResult[1:]
				}*/
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
			traverser.dispatch.Defun(name, *params, *body)
			return name
		}
		isDefParameter := false
		if string(token.Children[0].Value) == "defparameter" {
			isDefParameter = true
		}
		if isDefParameter {
			name := token.Children[1].Value
			value := token.Children[2]
			traverser.dispatch.Defparameter(name, value)
			return name
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
		isProgn := false
		if string(token.Children[0].Value) == "progn" {
			isProgn = true
		}
		if isProgn {
			if len(token.Children) > 1 {
				for i := 1; i < len(token.Children)-1; i++ {
					traverser.Interpret(*token.Children[i])
				}
				return traverser.Interpret(*token.Children[len(token.Children)-1])
			} else {
				return "()"
			}

		}
		isCar := false
		if string(token.Children[0].Value) == "car" {
			isCar = true
		}
		if isCar {
			return token.Children[1].Children[0].Value

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
		for index, child := range token.Children {
			// interpret this child recursively and save result
			var value string
			value = traverser.Interpret(*child)
			if index == 0 {
				value = child.Value
			}

			//If this list is a function call append all evaluated children to an array
			if isFunc {
				funcCall = append(funcCall, value)
			}
		}
		//If the list we just evaluated all of the children for was a function
		if isFunc {
			//.Print("Call trace: ")
			//fmt.Println(funcCall)
			//Pass the constructed array to the dispatcher to handle the function call
			result = traverser.dispatch.Call(funcCall)
		}
	} else {
		//This was not a list and we can just pass back the literal value of the token
		if token.Value == "(" {
			return "()"
		}
		if traverser.isSymbol(token.Value) {
			return traverser.dispatch.ResolveVar(token.Value)
		}
		return token.Value
	}
	// return the result of the list evaluation
	return
}

func (traverser Traverser) isSymbol(value string) bool {
	illegalStartChars := []string{"\""}
	for _, illegal := range illegalStartChars {
		if string(value[0]) == illegal {
			return false
		}
	}
	isAllNumbers := true
	for _, char := range value {
		if unicode.IsDigit(char) != true {
			isAllNumbers = false
		}
	}
	if isAllNumbers {
		return false
	}
	return true
}
