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
	Parent   *Token
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
	input = doubleSpaceReplace(input)
	input = strings.Replace(input, "  ", " ", strings.Count(input, "  "))
	splitInput := strings.Split(input, " ")
	return splitInput
}

func doubleSpaceReplace(input string) string {
	prev := ""
	result := ""
	for _, char := range input {
		if (prev == " " || prev == "@!@REPLACESPACE@!@") && string(char) == " " {
			prev = "@!@REPLACESPACE@!@"
			result += "@!@REPLACESPACE@!@"
		} else {
			prev = string(char)
			result += string(char)
		}
	}
	return result
}

//(+ (+ 2 3) 2 3 (- 2 3))
func parse(input []string) (*Token, int) {
	token := Token{[]*Token{}, nil,
		input[0]}
	input = input[1:]
	popped := 0
	if token.Value == "(" || token.Value == "'(" || token.Value == ",(" {
		for input[0] != ")" {
			if input[0] == "(" || input[0] == "'(" || input[0] == ",(" || string(input[0][0]) == "\"" {
				tVal, additionPop := parse(input)
				token.Children = append(token.Children, tVal)
				input = input[additionPop:]
				popped += additionPop
			} else {
				token.Children = append(token.Children, &Token{[]*Token{}, &token, input[0]})
			}
			input = input[1:]
			popped++
		}
		input = input[1:]
		popped++
	} else if string(token.Value[0]) == "\"" {
		token.Value = strings.TrimSuffix(token.Value, "\n")
		result := token.Value
		tokenLast := string(token.Value[len(token.Value)-1])
		if tokenLast != "\"" {
			last := string(input[0][len(input[0])-1])
			for last != "\"" {
				result += " " + input[0]
				input = input[1:]
				popped++
				last = string(input[0][len(input[0])-1])
			}
			result += " " + input[0]
			input = input[1:]
			popped++
			token.Value = result
		}
	}
	return &token, popped
}
