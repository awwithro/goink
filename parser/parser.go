package parser

import (
	"encoding/json"

	"github.com/awwithro/goink/parser/types"
)

func Parse(rawJson string) types.Ink {
	i := &types.Ink{}
	if err := json.Unmarshal([]byte(rawJson), i); err != nil {
		panic(err)
	}
	return *i
}
