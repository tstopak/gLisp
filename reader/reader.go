package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TokenTree struct {
	Root *Token
}

type Token struct {
	Children []*Token
	Value    string
}

func ReadInput() (tree TokenTree) {
	var input string
	fmt.Print("Glisp>>")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	root, _ := parse(formatInput(input))
	tree = TokenTree{Root: root}
	return tree
}

//(+ (+ 2 3) 2 3 (- (+ 2 3) 3))
func formatInput(input string) []string {
	input = strings.Replace(input, "(", "( ", strings.Count(input, "("))
	input = strings.Replace(input, ")", " ) ", strings.Count(input, ")"))
	input = strings.Replace(input, "  ", " ", strings.Count(input, "  "))
	splitInput := strings.Split(input, " ")
	return splitInput
}

//(+ (+ 2 3) 2 3 (- 2 3))
func parse(input []string) (*Token, int) {
	token := Token{[]*Token{},
		input[0]}
	input = input[1:]
	popped := 0
	if token.Value == "(" || token.Value == "'(" || token.Value == ",(" {
		for input[0] != ")" {
			if input[0] == "(" || input[0] == "'(" || input[0] == ",(" {
				tVal, additionPop := parse(input)
				token.Children = append(token.Children, tVal)
				input = input[additionPop:]
				popped += additionPop
			} else {
				token.Children = append(token.Children, &Token{[]*Token{}, input[0]})
			}
			input = input[1:]
			popped++
		}
		input = input[1:]
		popped++
	}
	return &token, popped
}
