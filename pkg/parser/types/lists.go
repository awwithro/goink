package types

import "strings"

var _ Comparable[ListValItem] = ListValItem{}

type ListDefs map[string]List
type List map[string]int

// Creates a new list
type ListInit struct {
	List List
	// Reference to the named lists in ListDefs
	Origins []string
}

type ListVal map[string]ListValItem

func (l ListVal) String() string {
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
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

type ListValItem struct {
	Name   string
	Value  int
	Parent string
}

func (l ListValItem) String() string {
	return l.Name
}

func (l ListValItem) Accept(v Visitor) {
	panic("not implemented")
	// v.VisitListValItem(l)
}

func (l ListInit) Accept(v Visitor) {
	v.VisitListInit(l)
}
