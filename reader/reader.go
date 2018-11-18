package reader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type TokenList struct{
	input	string
	tokens	[]string
}

type Token struct {
	tokenType string
	terminal bool
	children []*Token
}

func ReadInput()(input string){
	fmt.Println("Input an sexpr:")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	parseList(input)
	return
}

func parseList(input string){
	input = correctSpacing(input)
	fmt.Println("Input after correcting spacing: " + input)
	input = input[1:len(input)-2]
	strings.Split(input, " ")
	fmt.Println(input)

}

func correctSpacing(input string)(output string){
	internalspacer, _ := regexp.Compile(" +")
	leftparenSpacer, _:= regexp.Compile(`\(\s`)
	rightparenSpacer, _:= regexp.Compile(`\s\)`)
	output = internalspacer.ReplaceAllString(input, " ")
	output = leftparenSpacer.ReplaceAllString(output, "(")
	output = rightparenSpacer.ReplaceAllString(output, ")")
	return output
}