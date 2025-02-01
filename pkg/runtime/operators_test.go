package runtime

import (
	"testing"

	"github.com/awwithro/goink/pkg/parser/types"
	"github.com/emirpasic/gods/v2/stacks/arraystack"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc     string
		stack    []types.Acceptor
		op       types.Operator
		expected types.Acceptor
	}{
		{
			desc:     "Test Not",
			stack:    []types.Acceptor{types.FloatVal(0)},
			op:       types.Not,
			expected: types.FloatVal(1),
		},
		{
			desc:     "Test Negate",
			stack:    []types.Acceptor{types.FloatVal(10)},
			op:       types.Negate,
			expected: types.FloatVal(-10),
		},
		{
			desc:     "Test And True",
			stack:    []types.Acceptor{types.FloatVal(1), types.FloatVal(1)},
			op:       types.And,
			expected: types.BoolVal(true),
		},
		{
			desc:     "Test And False",
			stack:    []types.Acceptor{types.FloatVal(0), types.FloatVal(1)},
			op:       types.And,
			expected: types.BoolVal(false),
		},
		{
			desc:     "Test Or True",
			stack:    []types.Acceptor{types.FloatVal(0), types.FloatVal(1)},
			op:       types.Or,
			expected: types.BoolVal(true),
		},
		{
			desc:     "Test Or False",
			stack:    []types.Acceptor{types.FloatVal(0), types.FloatVal(0)},
			op:       types.Or,
			expected: types.BoolVal(false),
		},
		{
			desc:     "Test Min",
			stack:    []types.Acceptor{types.FloatVal(1), types.FloatVal(5)},
			op:       types.Min,
			expected: types.FloatVal(1),
		},
		{
			desc:     "Test Max",
			stack:    []types.Acceptor{types.FloatVal(1), types.FloatVal(5)},
			op:       types.Max,
			expected: types.FloatVal(5),
		},
		{
			desc:     "Test ListValueItems",
			stack:    []types.Acceptor{types.ListValItem{Name: "test", Parent: "test", Value: 1}, types.ListValItem{Name: "test", Parent: "test", Value: 1}},
			op:       types.Equal,
			expected: types.BoolVal(true),
		},
		{
			desc:     "Test ListValueItems",
			stack:    []types.Acceptor{types.ListValItem{Name: "test", Parent: "test", Value: 1}, types.ListValItem{Name: "test", Parent: "other", Value: 1}},
			op:       types.Equal,
			expected: types.BoolVal(false),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			s := Story{evaluationStack: arraystack.New[any]()}
			for _, item := range tC.stack {
				s.evaluationStack.Push(item)
			}
			s.VisitOperator(tC.op)
			actual, ok := s.evaluationStack.Pop()
			assert.True(ok)
			assert.Equal(tC.expected, actual)
		})
	}
}
