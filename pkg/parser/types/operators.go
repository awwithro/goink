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
	"LIST_ALL":    ListAll,
	"listInt":     ListInt,
	"INT":         Int,
	"FLOOR":       Floor,
	"FLOAT":       Float,
	"rnd":         Random,
}

func (o Operator) IsUnary() bool {
	return o == Negate || o == Not || o == ListValue ||
		o == ListMin || o == ListMax || o == ListRandom || o == ListCount ||
		o == Int || o == Floor || o == Float || o ==ListAll
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
	ListAll
	Int
	Floor
	Float
	Random
)

func IsOperator(str string) (Operator, bool) {
	c, ok := operatorMap[str]
	return c, ok
}

func (o Operator) Accept(v Visitor) {
	v.VisitOperator(o)
}
