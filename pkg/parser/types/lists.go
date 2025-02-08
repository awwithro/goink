package types

import (
	"math/rand"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"
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

type ListVal []*ListValItem

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
		newList := ListVal{}
		for itemName, itemVal := range list {
			lvi := &ListValItem{
				Name:   itemName,
				Value:  int(itemVal),
				Parent: &newList,
			}
			newList = append(newList, lvi)
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
	slices.SortFunc(l, func(a, b *ListValItem) int {
		if a.Value < b.Value {
			return -1
		} else if a.Value > b.Value {
			return 1
		}
		// equal values
		if a.Name < b.Name {
			return -1
		}
		return 1
	})
	return l
}
func (l ListVal) Get(name string) *ListValItem {
	for _, item := range l {
		if item.Name == name {
			return item
		}
	}
	return nil
}

func (l ListVal) Min() ListVal {
	if len(l) == 0 {
		return ListVal{}
	}
	var min *ListValItem
	for _, v := range l {
		if min == nil {
			min = v
		} else if v.Value < min.Value {
			min = v
		}
	}
	return ListVal{min}
}

func (l ListVal) Max() ListVal {
	if len(l) == 0{
		return ListVal{}
	}
	var max *ListValItem
	for _, v := range l {
		if max == nil {
			max = v
		} else if v.Value > max.Value {
			max = v
		}
	}
	return ListVal{max}
}

func (l ListVal) Random() ListVal {
	log.Debugf("picking from %v", l.All())
	if len(l) == 0 {
		log.Debug("Empty List")
	} else {
		i := rand.Intn(len(l))
		x := 0
		for _, val := range l {
			if x == i {
				return ListVal{val}
			}
			x++
		}
	}
	return ListVal{}
}

func (l ListVal) AsBool() bool {
	return len(l) > 0
}

func (l ListVal) All() (all ListVal) {
	lists := map[*ListVal]bool{}
	for _, item := range l {
		lists[item.Parent] = true
	}
	log.Debugf("Lists: %v", l)
	for list := range lists {
		for _, item := range list.AsList() {
			all = append(all, item)
		}
	}
	return all.AsList()
}

func (l ListVal) Count() int {
	return len(l)
}

func (l ListVal) String() string {
	keys := make([]string, 0, len(l))
	for _, k := range l.AsList() {
		keys = append(keys, k.Name)
	}
	return strings.Join(keys, ",")
}
func (l ListVal) GetValue(val int) (item *ListValItem) {
	for _, i := range l {
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
