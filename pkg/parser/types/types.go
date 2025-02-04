package types

import (
	"strconv"
)

var _ Truthy = IntVal(0)
var _ Truthy = FloatVal(0)
var _ Truthy = BoolVal(false)
var _ NumericVal = IntVal(0)
var _ NumericVal = FloatVal(0)
var _ NumericVal = BoolVal(false)

const (
	GlobalVarKey = "global decl"
)

type Truthy interface {
	AsBool() bool
}

type NumericVal interface {
	IsFloat() bool
	AsInt() int
	AsFloat() float64
	Truthy
	Comparable[NumericVal]
}

type Comparable[T any] interface {
	Equals(T) bool
	NotEquals(T) bool
	LT(T) bool
	LTE(T) bool
	GT(T) bool
	GTE(T) bool
}

type Ink struct {
	InkVersion int       `json:"inkVersion"`
	Root       Container `json:"root"`
	ListDefs   ListDefs  `json:"listDefs"`
}

type VoidVal struct{}

func (v VoidVal) String() string {
	return ""
}

func (v VoidVal) Accept(vs Visitor) {
	vs.VisitVoidVal(v)
}

type BoolVal bool

func (b BoolVal) AsBool() bool {
	return bool(b)
}

func (b BoolVal)String()string{
	if b{
		return "true"
	}
	return "false"
}

func (b BoolVal) Accept(v Visitor) {
	v.VisitBoolVal(b)
}

func (b BoolVal) Equals(other NumericVal) bool {
	return b.AsBool() == other.AsBool()
}

func (b BoolVal) NotEquals(other NumericVal) bool {
	return b.AsBool() == other.AsBool()
}
func (b BoolVal) LT(other NumericVal) bool {
	return b.AsInt() < other.AsInt()
}

func (b BoolVal) LTE(other NumericVal) bool {
	return b.AsInt() <= other.AsInt()
}
func (b BoolVal) GT(other NumericVal) bool {
	return b.AsInt() > other.AsInt()
}

func (b BoolVal) GTE(other NumericVal) bool {
	return b.AsInt() >= other.AsInt()
}

func (b BoolVal) IsFloat() bool {
	return false
}
func (b BoolVal) AsFloat() float64 {
	if b {
		return 1
	}
	return 0
}
func (b BoolVal) AsInt() int {
	if b {
		return 1
	}
	return 0
}

type StringVal string

func (s StringVal) String() string {
	return string(s)
}
func (s StringVal) Accept(v Visitor) {
	v.VisitString(s)
}

type IntVal int

func (i IntVal) Accept(v Visitor) {
	v.VisitIntVal(i)
}
func (i IntVal) String() string {
	return strconv.Itoa(int(i))
}
func (i IntVal) IsFloat() bool {
	return false
}
func (i IntVal) AsInt() int {
	return int(i)
}
func (i IntVal) AsFloat() float64 {
	return float64(i)
}
func (i IntVal) AsBool() bool {
	return i != 0
}

func (i IntVal) Equals(other NumericVal) bool {
	return i.AsInt() == other.AsInt()
}

func (i IntVal) NotEquals(other NumericVal) bool {
	return i.AsInt() == other.AsInt()
}
func (i IntVal) LT(other NumericVal) bool {
	return i.AsInt() < other.AsInt()
}

func (i IntVal) LTE(other NumericVal) bool {
	return i.AsInt() <= other.AsInt()
}
func (i IntVal) GT(other NumericVal) bool {
	return i.AsInt() > other.AsInt()
}

func (i IntVal) GTE(other NumericVal) bool {
	return i.AsInt() >= other.AsInt()
}

type FloatVal float64

func (f FloatVal) Accept(v Visitor) {
	v.VisitFloatVal(f)
}
func (f FloatVal) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (f FloatVal) IsFloat() bool {
	return true
}
func (f FloatVal) AsInt() int {
	return int(f)
}
func (f FloatVal) AsFloat() float64 {
	return float64(f)
}
func (i FloatVal) AsBool() bool {
	return i != 0
}

func (f FloatVal) Equals(other NumericVal) bool {
	return f.AsFloat() == other.AsFloat()
}

func (f FloatVal) NotEquals(other NumericVal) bool {
	return f.AsFloat() == other.AsFloat()
}

func (f FloatVal) LT(other NumericVal) bool {
	return f.AsFloat() < other.AsFloat()
}

func (f FloatVal) LTE(other NumericVal) bool {
	return f.AsFloat() <= other.AsFloat()
}
func (f FloatVal) GT(other NumericVal) bool {
	return f.AsFloat() > other.AsFloat()
}

func (f FloatVal) GTE(other NumericVal) bool {
	return f.AsFloat() >= other.AsFloat()
}

type DivertTarget Path

func (d DivertTarget) Accept(v Visitor) {
	v.VisitDivertTarget(d)
}

type Divert struct {
	Path        Path
	Conditional bool
}

func (d Divert) Accept(v Visitor) {
	v.VisitDivert(d)
}

type VariableDivert struct {
	Name        string
	Conditional bool
}

func (d VariableDivert) Accept(v Visitor) {
	v.VisitVariableDivert(d)
}

type FunctionDivert struct {
	Divert
}

func (f FunctionDivert) Accept(v Visitor) {
	v.VisitFunctionDivert(f)
}

type TunnelDivert struct {
	Divert
}

func (t TunnelDivert) Accept(v Visitor) {
	v.VisitTunnelDivert(t)
}

type ExternalFunctionDivert struct {
	Divert
	Args int
}

func (e ExternalFunctionDivert) Accept(v Visitor) {
	v.VisitExternalFunctionDivert(e)
}

type ChoicePoint struct {
	Path Path
	Flag byte
}

func (c ChoicePoint) Accept(v Visitor) {
	v.VisitChoicePoint(c)
}

func (c *ChoicePoint) HasCondition() bool {
	return c.Flag&0x1 == 0x1
}
func (c *ChoicePoint) HasStartContent() bool {
	return c.Flag&0x2 == 0x2
}
func (c *ChoicePoint) HasChoiceOnly() bool {
	return c.Flag&0x4 == 0x4
}
func (c *ChoicePoint) IsInvisibleDefault() bool {
	return c.Flag&0x8 == 0x8
}
func (c *ChoicePoint) OnceOnly() bool {
	return c.Flag&0x10 == 0x10
}

type ReadCount Path

func (r ReadCount) Accept(v Visitor) {
	v.VisitReadCount(r)
}

type VarRef string

func (vr VarRef) Accept(v Visitor) {
	v.VisitVarRef(vr)
}

type GlobalVar struct {
	Name     string
	ReAssign bool
}

func (g GlobalVar) Accept(v Visitor) {
	v.VisitGlobalVar(g)
}

type TempVar struct {
	Name     string
	ReAssign bool
}

func (t TempVar) Accept(v Visitor) {
	v.VisitTmpVar(t)
}

type VariablePointer struct {
	Name         string
	ContextIndex int
}

func (p VariablePointer) Accept(v Visitor) {
	v.VisitVariablePointer(p)
}

func NewVariablePointer(name string) VariablePointer {
	return VariablePointer{
		Name:         name,
		ContextIndex: -1,
	}
}
