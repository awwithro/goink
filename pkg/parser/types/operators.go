package types

type Operator int

var operatorMap = map[string]Operator{
	"+":           Plus,
	"-":           Minus,
	"/":           Divide,
	"*":           Multiply,
	"%":           Modulus,
	"_":           Negate,
	"==":          Equal,
	">":           GreaterThan,
	"<":           LessThan,
	">=":          GreaterThanEqual,
	"<=":          LessThanEqual,
	"!=":          NotEqual,
	"!":           Not,
	"&&":          And,
	"||":          Or,
	"MIN":         Min,
	"MAX":         Max,
	"LIST_VALUE":  ListValue,
	"LIST_MIN":    ListMin,
	"LIST_MAX":    ListMax,
	"LIST_COUNT":  ListCount,
	"LIST_RANDOM": ListRandom,
	"listInt":     ListInt,
}

func (o Operator) IsUnary() bool {
	return o == Negate || o == Not || o == ListValue || o == ListMin || o == ListMax || o == ListRandom || o == ListCount
}

const (
	Plus Operator = iota
	Minus
	Divide
	Multiply
	Modulus
	Negate
	Equal
	GreaterThan
	LessThan
	GreaterThanEqual
	LessThanEqual
	NotEqual
	Not
	And
	Or
	Min
	Max
	ListValue
	ListInt
	ListMin
	ListMax
	ListRandom
	ListCount
)

func IsOperator(str string) (Operator, bool) {
	c, ok := operatorMap[str]
	return c, ok
}

func (o Operator) Accept(v Visitor) {
	v.VisitOperator(o)
}
