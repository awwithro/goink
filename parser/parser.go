package parser

import (
	"encoding/json"

	"github.com/awwithro/goink/parser/types"
	"github.com/sirupsen/logrus"
)

func Parse(rawJson []byte) types.Ink {
	i := &types.Ink{}
	if err := json.Unmarshal(rawJson, i); err != nil {
		logrus.Panic(err)
	}
	return *i
}
