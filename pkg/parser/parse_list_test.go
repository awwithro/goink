package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	assert := assert.New(t)
	js, err := os.ReadFile("../../examples/list.json")
	assert.NoError(err)
	ink := Parse(js)
	l, ok := ink.ListDefs["kettleState"]
	assert.True(ok)
	assert.Equal(1, l["cold"])
}
