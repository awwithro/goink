package parser

import (
	"encoding/json"

	"github.com/awwithro/goink/parser/types"
)

func Parse(rawJson []byte) types.Ink {
	c := types.NewContainer("", nil)
	i := &types.Ink{
		Root: *c,
	}
	if err := json.Unmarshal(rawJson, i); err != nil {
		panic(err)
	}
	return *i
}
