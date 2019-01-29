package reader

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
)

type TokenTree struct {
	Root *Token
}

type Token struct {
	Children []*Token
	Value    string
}

func ReadInput() (tree *TokenTree) {
	var input string
	fmt.Println("Input an sexpr:")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	root := parseList(input)
	tree = &TokenTree{root}
	return
}

//(print (+ 2 3) (+ 4 5))
func parseList(input string) *Token {
	var root = Token{}
	tokenStrings := SplitInput(input)
	root.Children = make([]*Token, 0, len(tokenStrings))
	for _, tokenVal := range tokenStrings {
		firstChar := string(tokenVal[0])
		if firstChar == "(" {
			root.Children = append(root.Children, parseList(tokenVal))
		} else {
			root.Children = append(root.Children, &Token{nil, tokenVal})
		}
	}
	return &root
}

func SplitInput(input string) (tokenList []string) {
	newlineRemover, _ := regexp.Compile(`\n`)
	input = newlineRemover.ReplaceAllString(input, "")
	if string(input[0]) == "(" && string(input[len(input)-1]) == ")" {
		input = input[1 : len(input)-1]
	}
	syntax, _ :=
		regexp.Compile(`(;+)(.*?)(;+)|(("*?)|('*?))(\(+).*?(\)+)"*|".*?"|[[:graph:]]+`)
	tokens := syntax.FindAllStringIndex(input, math.MaxInt64)
	for _, location := range tokens {
		tokenList = append(tokenList, input[location[0]:location[1]])
	}
	return
}
