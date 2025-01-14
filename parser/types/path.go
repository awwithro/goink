package types

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	ParentContainer string = "^"
)

type Path string

type Segment struct {
	Addr   int
	Name   string
	IsAddr bool // If the path is an Addr or a Name
}

// Returns the elements in the path, ints or strings
func (p Path) Segments() []Segment {
	segments := []Segment{}
	s := strings.Split(string(p), ".")
	for _, x := range s {
		// current container identifiers start with .^ which creates the empty string
		if x == "" {
			continue
		}
		if i, err := strconv.Atoi(x); err != nil {
			segments = append(segments, Segment{Name: x, IsAddr: false})
		} else {
			segments = append(segments, Segment{Addr: i, IsAddr: true})
		}
	}
	return segments
}

func ResolvePath(p Path, root, current *Container) *Container {
	segs := p.Segments()
	// starts with an address
	if segs[0].IsAddr {
		cnt := root
		for _, seg := range segs {
			if seg.IsAddr {
				x := cnt.Contents[seg.Addr]
				if c, ok := x.(*Container); ok {
					cnt = c
				} else {
					logrus.Panicf("Address container non container element %v", reflect.TypeOf(x))
				}
			} else {
				c, err := cnt.GetNamedContainer(seg.Name)
				if err != nil {
					logrus.Panic(err)
				}
				cnt = c
			}
		}
		return cnt
		// starts with '^' ie a local ref
	} else if segs[0].Name == ParentContainer {
		ct, err := current.GetNamedContainer(segs[1].Name)
		if err != nil {
			panic(err)
		}
		return ct
	}
	return nil
}
