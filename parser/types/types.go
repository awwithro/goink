package types

type Ink struct {
	InkVersion int       `json:"inkVersion"`
	Root       Container `json:"root"`
}

type StringVal string

func (s StringVal) String() string {
	return string(s)
}

type IntVal int
type FloatVal float64

type NumericVal interface {
	IsFloat() bool
	AsInt() int
	AsFloat() float64
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

func (f FloatVal) IsFloat() bool {
	return true
}
func (f FloatVal) AsInt() int {
	return int(f)
}
func (f FloatVal) AsFloat() float64 {
	return float64(f)
}

type DivertTarget Path

type Divert struct {
	Path        Path
	Conditional bool
}

type VariableDivert struct {
	Name        string
	Conditional bool
}

type FunctionDivert struct {
	Path        Path
	Conditional bool
}

type TunnelDivert struct {
	Path        Path
	Conditional bool
}

type ExternalFunctionDivert struct {
	Path        Path
	Args        int
	Conditional bool
}

type ChoicePoint struct {
	Path Path
	Flag byte
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
func (c *ChoicePoint) IsInvisible() bool {
	return c.Flag&0x8 == 0x8
}
func (c *ChoicePoint) OnceOnly() bool {
	return c.Flag&0x10 == 0x10
}

type ReadCount string

type VarRef string

type GlobalVar struct {
	Name     string
	ReAssign bool
}

type TempVar struct {
	Name     string
	ReAssign bool
}

type VariablePointer struct {
	Name         string
	ContextIndex int
}

func NewVariablePointer(name string) VariablePointer {
	return VariablePointer{
		Name:         name,
		ContextIndex: -1,
	}
}
