package runtime

import (
	"strconv"

	"github.com/awwithro/goink/parser/types"
)

type Mode int

func (m Mode) String() string {
	switch m {
	case None:
		return "None"
	case Str:
		return "Str"
	case Eval:
		return "Eval"
	default:
		return strconv.Itoa(int(m))
	}

}

const (
	None Mode = iota
	Str
	Eval
)

var _ types.Visitor = &Story{}
