package types

import (
	"math/rand"
	"slices"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	log "github.com/sirupsen/logrus"
)

var _ Comparable[ListValItem] = &ListValItem{}
var _ Inty = &ListVal{}

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

type ListVal struct {
	mapset.Set[*ListValItem]
}

func NewListVal(items ...*ListValItem) ListVal {
	return ListVal{
		Set: mapset.NewSet(items...),
	}
}

// returns a map of list value names. If the key is duplicated, the value is true
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
		newList := NewListVal()
		for itemName, itemVal := range list {
			lvi := &ListValItem{
				Name:   itemName,
				Value:  int(itemVal),
				Parent: &newList,
			}
			newList.Add(lvi)
		}
		sorted := newList.ToSortedSlice()
		for x, item := range sorted {
			if x+1 < len(sorted) {
				item.Next = sorted[x+1]
			}
		}
		result[listName] = newList
	}
	return result
}

func (l ListVal) ToSortedSlice() []*ListValItem {
	items := l.ToSlice()
	slices.SortFunc(items, func(a, b *ListValItem) int {
		if a.Value < b.Value {
			return -1
		}
		if a.Value > b.Value {
			return 1
		}
		if a.Value == b.Value {
			if a.Name < b.Name {
				return -1
			}
		}
		return 1
	})
	return items
}

func (l ListVal) Get(name string) *ListValItem {
	for _, item := range l.ToSlice() {
		if item.Name == name {
			return item
		}
	}
	return nil
}

func (l ListVal) Min() ListVal {
	if l.Count() == 0 {
		return NewListVal()
	}
	var min *ListValItem
	for _, v := range l.ToSlice() {
		if min == nil {
			min = v
		} else if v.Value < min.Value {
			min = v
		}
	}
	return NewListVal(min)
}

func (l ListVal) Max() ListVal {
	if l.Count() == 0 {
		return NewListVal()
	}
	var max *ListValItem
	for _, v := range l.ToSlice() {
		if max == nil {
			max = v
		} else if v.Value > max.Value {
			max = v
		}
	}
	return NewListVal(max)
}

func (l ListVal) Random() ListVal {
	log.Debugf("picking from %v", l.All())
	if l.Count() == 0 {
		log.Debug("Empty List")
	} else {
		i := rand.Intn(l.Count())
		x := 0
		for _, val := range l.ToSlice() {
			if x == i {
				return NewListVal(val)
			}
			x++
		}
	}
	return NewListVal()
}

func (l ListVal) AsBool() bool {
	return l.Count() > 0
}

func (l ListVal) All() ListVal {
	all := NewListVal()
	for _, item := range l.ToSlice() {
		all.Set = all.Union(item.Parent.Set)
	}
	return all
}

func (l ListVal) Count() int {
	return len(l.ToSlice())
}

func (l ListVal) AsInt() int {
	if l.Count() > 1 {
		log.Panic("can't turn a multi-value list to an int")
	}
	return l.ToSlice()[0].Value
}

func (l ListVal) Range(min, max int) ListVal {
	res := NewListVal()
	for _, val := range l.ToSortedSlice() {
		if val.Value >= min && val.Value <= max {
			res.Add(val)
		}
	}
	return res
}

func (l ListVal) String() string {
	keys := make([]string, 0, l.Count())
	for _, k := range l.ToSortedSlice() {
		keys = append(keys, k.Name)
	}
	return strings.Join(keys, ",")
}
func (l ListVal) GetValue(val int) (item *ListValItem) {
	for _, i := range l.ToSortedSlice() {
		if val == i.Value {
			item = i
			break
		}
	}
	if item == nil {
		log.Panicf("No item with value %d", val)
	}
	return item
}

func (l *ListValItem) AsBool() bool {
	// can this ever be false
	return l != nil
}
func (l ListValItem) Equals(other ListValItem) bool {
	return l.Name == other.Name && l.Value == other.Value && l.Parent == other.Parent
}

func (l ListValItem) NotEquals(other ListValItem) bool {
	return !l.Equals(other)
}

func (l ListValItem) GT(other ListValItem) bool {
	return l.Value > other.Value
}

func (l ListValItem) GTE(other ListValItem) bool {
	return l.Value >= other.Value
}

func (l ListValItem) LT(other ListValItem) bool {
	return l.Value < other.Value
}

func (l ListValItem) LTE(other ListValItem) bool {
	return l.Value <= other.Value
}

func (l ListValItem) String() string {
	return l.Name
}

func (l ListInit) Accept(v Visitor) {
	v.VisitListInit(l)
}
