package runtime

import (
	"strconv"

	"github.com/awwithro/goink/pkg/parser/types"
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
	case TagMode:
		return "Tag"
	default:
		return strconv.Itoa(int(m))
	}

}

const (
	None Mode = iota
	Str
	Eval
	TagMode
)

var _ types.Visitor = &Story{}
