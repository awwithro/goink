package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	testCases := []struct {
		desc       string
		path       Path
		assertions func(a *assert.Assertions, seg []Segment)
	}{
		{
			desc: "Simple Path",
			path: "0.c-1",
			assertions: func(a *assert.Assertions, seg []Segment) {
				a.True(seg[0].IsAddr, "First segment should be an address")
				a.Equal(seg[0].Addr, 0, "First segment should Equal 0")
				a.False(seg[1].IsAddr, "Second segment should be a Name")
				a.Equal(seg[1].Name, "c-1", "Second segment Name should be 'c-1'")
			},
		},
		{
			desc: "Relative Path",
			path: ".^.$r1",
			assertions: func(a *assert.Assertions, seg []Segment) {
				a.False(seg[0].IsAddr, "First segment should not be an address")
				a.Equal(seg[0].Name, "^", "First segment should Equal '^")
				a.False(seg[1].IsAddr, "Second segment should be a Name")
				a.Equal(seg[1].Name, "$r1", "Second segment Name should be '$r1'")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			seg := tC.path.Segments()
			tC.assertions(assert, seg)
		})
	}
}
