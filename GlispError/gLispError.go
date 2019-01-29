package GlispError

import "fmt"

type GLispError struct {
	name  string
	cause string
}

func (gle GLispError) PrintError() {
	fmt.Print("Error Thrown! \n Name: \n Cause: ")
}
