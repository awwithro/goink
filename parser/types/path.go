package types

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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

// Resolves to a container and the index of the contents
func ResolvePath(p Path, current *Container) (*Container, int) {
	var idx int
	root := current.GetRoot()
	segs := p.Segments()
	// starts with an address
	if segs[0].IsAddr || segs[0].Name != ParentContainer {
		cnt := root
		for _, seg := range segs {
			if seg.IsAddr {
				x := cnt.Contents[seg.Addr]
				if c, ok := x.(*Container); ok {
					cnt = c
					idx = 0
				} else {
					idx = seg.Addr
				}
			} else {
				c, err := cnt.GetNamedContainer(seg.Name)
				if err != nil {
					log.Panic(err)
				}
				cnt = c
				idx = 0
			}
		}
		return cnt, idx
		// starts with '^' ie a local ref
	} else {
		checkedCurrent := false
		for _, seg := range segs {
			if seg.IsAddr {
				ct := current.Contents[seg.Addr]
				if c, ok := ct.(*Container); ok {
					return c, 0
				} else {
					return current, seg.Addr
				}
			} else if seg.Name != ParentContainer {
				ct, err := current.GetNamedContainer(seg.Name)
				if err != nil {
					log.Panic(err)
				}
				return ct, 0
			} else {
				if !checkedCurrent {
					checkedCurrent = true
				} else {
					current = current.ParentContainer
				}
			}
		}

	}
	panic("shouldn't be here")
}
