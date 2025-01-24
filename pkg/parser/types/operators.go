package types

type Operator int

var operatorMap = map[string]Operator{
	"+":   Plus,
	"-":   Minus,
	"/":   Divide,
	"*":   Multiply,
	"%":   Modulus,
	"_":   Negate,
	"==":  Equal,
	">":   GreaterThan,
	"<":   LessThan,
	">=":  GreaterThanEqual,
	"<=":  LessThanEqual,
	"!=":  NotEqual,
	"!":   Not,
	"&&":  And,
	"||":  Or,
	"MIN": Min,
	"MAX": Max,
}

func (o Operator) IsUnary() bool {
	return o == Negate || o == Not
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
)

func IsOperator(str string) (Operator, bool) {
	c, ok := operatorMap[str]
	return c, ok
}

func (o Operator) Accept(v Visitor) {
	v.VisitOperator(o)
}
