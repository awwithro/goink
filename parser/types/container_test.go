package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRootCotnainer(t *testing.T) {
	root := NewContainer("Foo", nil)
	c1 := NewContainer("C1", root)
	c2 := NewContainer("C2", c1)
	c3 := NewContainer("C3", c2)
	assert.Equal(t, root, c3.GetRoot())

}

func TestGetNamedContainer(t *testing.T) {
	assert := assert.New(t)
	foo := &Container{Name: "Foo"}
	baz := &Container{Name: "Baz"}
	bar := &Container{Name: "Bar", SubContainers: map[string]*Container{"Baz": baz}}

	c := &Container{
		Contents: []Acceptor{
			StringVal("String"),
			foo,
			bar,
		},
	}
	actual, err := c.GetNamedContainer("Foo")
	assert.NoError(err)
	assert.Equal(foo, actual)
	actual, err = bar.GetNamedContainer("Baz")
	assert.NoError(err)
	assert.Equal(baz, actual)
	_, err = c.GetNamedContainer("BLAH")
	assert.Error(err)
}

func TestResolvePath(t *testing.T) {
	root := NewContainer("Root", nil)
	c := NewContainer("Start", root)
	foo := NewContainer("Foo", c)
	baz := NewContainer("Baz", nil)
	bar := NewContainer("Bar", c)
	bar.SubContainers["Baz"] = baz
	baz.ParentContainer = bar
	c.Contents = []Acceptor{
		StringVal("String"),
		foo,
		bar,
	}
	root.Contents = []Acceptor{
		c,
	}

	testCases := []struct {
		desc             string
		path             Path
		rootContainer    *Container
		currentContainer *Container
		expected         *Container
		panics           bool
	}{
		{
			desc:             "Parent Lookup",
			path:             ".^.Foo",
			currentContainer: c,
			expected:         foo,
		},
		{
			desc:             "Address Lookup",
			path:             "0.1",
			currentContainer: c,
			expected:         foo,
		},
		{
			desc:             "Mixed address and Name Lookup",
			path:             "0.2.Baz",
			currentContainer: c,
			expected:         baz,
		},
		{
			desc:             "Failed lookup",
			path:             "1.9.Florb",
			currentContainer: c,
			expected:         nil,
			panics:           true,
		},
		{
			desc:             "Multiple Parent Refs",
			path:             ".^.^.^.Foo",
			currentContainer: baz,
			expected:         foo,
		},
		{
			desc:             "Multiple Parent Refs and addrs",
			path:             ".^.^.^.^.0",
			currentContainer: baz,
			expected:         c,
		},
	}
	for _, tC := range testCases {
		assert := assert.New(t)
		t.Run(tC.desc, func(t *testing.T) {
			if tC.panics {
				assert.Panics(func() {
					ResolvePath(tC.path, tC.currentContainer)
				})
			} else {
				actual := ResolvePath(tC.path, tC.currentContainer)
				assert.Equal(tC.expected.Name, actual.Name)
			}
		})
	}
}
