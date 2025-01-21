package types

import "strconv"

type Ink struct {
	InkVersion int       `json:"inkVersion"`
	Root       Container `json:"root"`
}
type BoolVal bool

func (b BoolVal) Accept(v Visitor) {
	v.VisitBoolVal(b)
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

type FloatVal float64

func (f FloatVal) Accept(v Visitor) {
	v.VisitFloatVal(f)
}
func (f FloatVal) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

type NumericVal interface {
	IsFloat() bool
	AsInt() int
	AsFloat() float64
	AsBool() bool
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
	Path        Path
	Conditional bool
}

func (f FunctionDivert) Accept(v Visitor) {
	//v.VisitFunction(d)
	panic("not implemented")
}

type TunnelDivert struct {
	Path        Path
	Conditional bool
}

func (t TunnelDivert) Accept(v Visitor) {
	//v.TunnelDivert(d)
	panic("not implemented")
}

type ExternalFunctionDivert struct {
	Path        Path
	Args        int
	Conditional bool
}

func (e ExternalFunctionDivert) Accept(v Visitor) {
	//v.ExternalFunctionDivert(d)
	panic("not implemented")
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
