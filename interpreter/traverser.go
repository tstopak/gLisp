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
		return traverser.generateQuote(token)
	} else if isList {
		result = traverser.evaluateList(token)
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
	illegalStartChars := []string{"\"", " "}
	if value == "" {
		panic("Unexpected space in list")
	}
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

func (traverser Traverser) evaluateList(token reader.Token) string {
	switch keyword := string(token.Children[0].Value); keyword {
	case "defun":
		name := token.Children[1].Value
		params := token.Children[2]
		body := token.Children[3]
		traverser.dispatch.Defun(name, *params, *body)
		return name
	case "defparameter":
		name := token.Children[1].Value
		value := token.Children[2]
		traverser.dispatch.Defparameter(name, value)
		return name

	case "if":
		testResult := traverser.Interpret(*token.Children[1])
		if testResult == "true" {
			return traverser.Interpret(*token.Children[2])
		} else {
			return traverser.Interpret(*token.Children[3])
		}
	case "progn":
		if len(token.Children) > 1 {
			for i := 1; i < len(token.Children)-1; i++ {
				traverser.Interpret(*token.Children[i])
			}
			return traverser.Interpret(*token.Children[len(token.Children)-1])
		} else {
			return "()"
		}
	case "car":
		return token.Children[1].Children[0].Value
	default:
		funcCall := make([]string, 0, len(token.Children))
		for index, child := range token.Children {
			// interpret this child recursively and save result
			var value string
			value = traverser.Interpret(*child)
			if index == 0 {
				value = child.Value
			}
			funcCall = append(funcCall, value)
		}
		return traverser.dispatch.Call(funcCall)
	}

}

func (traverser Traverser) generateQuote(token reader.Token) (result string) {
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
}
