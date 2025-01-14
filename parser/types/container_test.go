package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamedContainer(t *testing.T) {
	assert := assert.New(t)
	foo := &Container{Name: "Foo"}
	baz := &Container{Name: "Baz"}
	bar := &Container{Name: "Bar", SubContainers: map[string]*Container{"Baz": baz}}

	c := &Container{
		Contents: []any{
			"^String",
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

func Test(t *testing.T) {
	foo := &Container{Name: "Foo"}
	baz := &Container{Name: "Baz"}
	bar := &Container{Name: "Bar", SubContainers: map[string]*Container{"Baz": baz}}
	c := &Container{
		Contents: []any{
			"^String",
			foo,
			bar,
		},
	}
	root := &Container{
		Contents: []any{
			c,
		},
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
			rootContainer:    root,
			currentContainer: c,
			expected:         foo,
		},
		{
			desc:             "Address Lookup",
			path:             "0.1",
			rootContainer:    root,
			currentContainer: c,
			expected:         foo,
		},
		{
			desc:             "Mixed address and Name Lookup",
			path:             "0.2.Baz",
			rootContainer:    root,
			currentContainer: c,
			expected:         baz,
		},
		{
			desc:             "Failed lookup",
			path:             "1.9.Florb",
			rootContainer:    root,
			currentContainer: c,
			expected:         nil,
			panics:           true,
		},
	}
	for _, tC := range testCases {
		assert := assert.New(t)
		t.Run(tC.desc, func(t *testing.T) {
			if tC.panics {
				assert.Panics(func() {
					ResolvePath(tC.path, tC.rootContainer, tC.currentContainer)
				})
			} else {
				actual := ResolvePath(tC.path, tC.rootContainer, tC.currentContainer)
				assert.Equal(tC.expected, actual)
			}
		})
	}
}
