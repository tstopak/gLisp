package core

import "golisp/reader"

type Function struct {
	Param reader.Token
	Body  reader.Token
}
