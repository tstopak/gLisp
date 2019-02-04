package GlispError

import "fmt"

type GLispError struct {
	name  string
	cause string
}

func (gle GLispError) PrintError() {
	fmt.Print("Error Thrown! \n Name: \n Cause: ")
}

func (gle *GLispError) SetName(name string) *GLispError {
	gle.name = name
	return gle
}

func (gle *GLispError) SetCause(cause string) *GLispError {
	gle.cause = cause
	return gle
}
