package core

import "glisp/reader"

type Function struct {
	Param reader.Token
	Body  reader.Token
}
