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
	"LIST_INVERT": ListInvert,
	"lrnd":        ListRandom,
	"LIST_ALL":    ListAll,
	"listInt":     ListInt,
	"INT":         Int,
	"FLOOR":       Floor,
	"FLOAT":       Float,
	"rnd":         Random,
	"range":       ListRange,
	"L^":          ListIntersect,
	"srnd":        SeedRandom,
	"?":           Contains,
	"!?":          NotContains,
}

func (o Operator) IsUnary() bool {
	return o == Negate || o == Not || o == ListValue ||
		o == ListMin || o == ListMax || o == ListRandom || o == ListCount ||
		o == Int || o == Floor || o == Float || o == ListAll || o == ListInvert || o == SeedRandom
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
	ListInvert
	ListAll
	Int
	Floor
	Float
	Random
	ListRange
	ListIntersect
	SeedRandom
	Contains
	NotContains
)

func IsOperator(str string) (Operator, bool) {
	c, ok := operatorMap[str]
	return c, ok
}

func (o Operator) Accept(v Visitor) {
	v.VisitOperator(o)
}
