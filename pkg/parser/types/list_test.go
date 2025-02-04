package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDef(t *testing.T) {
	assert := assert.New(t)
	ld := ListDefs{
		"one": listDef{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		"two": listDef{
			"two":   2,
			"three": 3,
			"four":  4,
		},
	}
	dupes := ld.GetDuplicatedKeys()
	assert.False(dupes["one"])
	assert.True(dupes["two"])
	assert.True(dupes["three"])
	assert.False(dupes["four"])

}
