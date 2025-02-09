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

func TestListValEquality(t *testing.T) {
	assert := assert.New(t)
	l1 := ListValItem{Name: "test", Parent: nil, Value: 1}
	l2 := ListValItem{Name: "test", Parent: nil, Value: 1}
	assert.True(l1.Equals(l2))

	l1 = ListValItem{Name: "test", Parent: nil, Value: 1}
	l2 = ListValItem{Name: "test", Parent: nil, Value: 2}

	assert.False(l1.Equals(l2))

}

func TestRange(t *testing.T) {
	assert := assert.New(t)
	primes := []*ListValItem{
		{Name: "two", Value: 2},
		{Name: "three", Value: 3},
		{Name: "five", Value: 5},
		{Name: "seven", Value: 7},
		{Name: "eleven", Value: 11},
		{Name: "thirteen", Value: 13},
		{Name: "seventeen", Value: 17},
		{Name: "nineteen", Value: 19},
		{Name: "twentythree", Value: 23},
	}
	list := NewListVal(primes...)
	actual := list.Range(10, 20)
	assert.Equal("eleven,thirteen,seventeen,nineteen", actual.String())

}
