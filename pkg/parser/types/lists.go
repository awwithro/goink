package types

import (
	"maps"
	"math/rand"
	"slices"
	"strings"
)

var _ Comparable[*ListValItem] = &ListValItem{}

// Initial representation
type ListDefs map[string]listDef

type listDef map[string]int

type ListValItem struct {
	Name   string
	Value  int
	Parent *ListVal
	Next   *ListValItem
}

// Creates a new list
type ListInit struct {
	List listDef
	// Reference to the named lists in ListDefs
	Origins []string
}

type ListVal map[string]*ListValItem

// returns a map of list keys. If the key is duplicated, the value is true
func (l *ListDefs) GetDuplicatedKeys() map[string]bool {
	result := make(map[string]bool)
	for _, lst := range *l {
		for k := range lst {
			if _, seen := result[k]; seen {
				result[k] = true
			} else {
				result[k] = false
			}
		}
	}
	return result
}

func (l *ListDefs) GetListValItems() map[string]ListVal {
	result := map[string]ListVal{}
	for listName, list := range *l {
		newList := ListVal{}
		for itemName, itemVal := range list {
			lvi := &ListValItem{
				Name:  itemName,
				Value: int(itemVal),
				Parent: &newList,
			}
			newList[itemName] = lvi
		}
		sorted := newList.AsList()
		for x, item := range sorted {
			if x+1 < len(sorted) {
				item.Next = sorted[x+1]
			}
		}
		result[listName] = newList
	}
	return result
}

func (l ListVal) AsList() []*ListValItem {
	return slices.SortedFunc(maps.Values(l), func(a, b *ListValItem) int {
		if a.Value < b.Value {
			return -1
		}
		return 1
	})
}

func (l ListVal) Min() (min *ListValItem) {
	for _, v := range l {
		if min == nil {
			min = v
		} else if v.Value < min.Value {
			min = v
		}
	}
	return min
}

func (l ListVal) Max() (max *ListValItem) {
	for _, v := range  l{
		if max == nil {
			max = v
		} else if v.Value > max.Value {
			max = v
		}
	}
	return max
}

func (l ListVal) Random() *ListValItem {
	i := rand.Intn(len(l))
	x := 0
	for _, val := range l {
		if x == i {
			return val
		}
		x++
	}
	return nil
}

func (l ListVal) Count() int {
	return len(l)
}

func (l ListVal) String() string {
	keys := make([]string, 0, len(l))
	for _,k := range l.AsList() {
		keys = append(keys, k.Name)
	}
	return strings.Join(keys, ", ")
}

func (l *ListValItem) Equals(other *ListValItem) bool {
	return l.Name == other.Name && l.Value == other.Value && l.Parent == other.Parent
}

func (l *ListValItem) NotEquals(other *ListValItem) bool {
	return !l.Equals(other)
}

func (l *ListValItem) GT(other *ListValItem) bool {
	return l.Value > other.Value
}

func (l *ListValItem) GTE(other *ListValItem) bool {
	return l.Value >= other.Value
}

func (l *ListValItem) LT(other *ListValItem) bool {
	return l.Value < other.Value
}

func (l *ListValItem) LTE(other *ListValItem) bool {
	return l.Value <= other.Value
}

func (l ListValItem) String() string {
	return l.Name
}

func (l *ListValItem) Accept(v Visitor) {
	v.VisitListValItem(l)
}

func (l ListInit) Accept(v Visitor) {
	v.VisitListInit(l)
}
