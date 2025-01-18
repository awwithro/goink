package runtime

import (
	"math"

	"github.com/awwithro/goink/parser/types"
	log "github.com/sirupsen/logrus"
)

func (s *Story) VisitOperator(op types.Operator) {
	log.Debugf("Visiting Operator: %d", op)
	if op.IsUnary() {
		val := mustPopNumeric(s.evaluationStack)
		log.Debug("Operating on ", val)
		switch op {
		case types.Negate:
			s.evaluationStack.Push(unaryOperator(val, negate))
		case types.Not:
			s.evaluationStack.Push(unaryOperator(val, not))
		}
	} else {
		val2 := mustPopNumeric(s.evaluationStack)
		val1 := mustPopNumeric(s.evaluationStack)
		log.Debug("Operating on ", val1, val2)
		switch op {
		case types.Plus:
			s.evaluationStack.Push(binaryOperator(val1, val2, add))
		case types.Minus:
			s.evaluationStack.Push(binaryOperator(val1, val2, sub))
		case types.Multiply:
			s.evaluationStack.Push(binaryOperator(val1, val2, mult))
		case types.Divide:
			s.evaluationStack.Push(binaryOperator(val1, val2, div))
		case types.Modulus:
			s.evaluationStack.Push(binaryOperator(val1, val2, mod))
		case types.Equal:
			s.evaluationStack.Push(binaryOperator(val1, val2, eq))
		case types.NotEqual:
			s.evaluationStack.Push(binaryOperator(val1, val2, neq))
		case types.LessThan:
			s.evaluationStack.Push(binaryOperator(val1, val2, lt))
		case types.LessThanEqual:
			s.evaluationStack.Push(binaryOperator(val1, val2, lte))
		case types.GreaterThan:
			s.evaluationStack.Push(binaryOperator(val1, val2, gt))
		case types.GreaterThanEqual:
			s.evaluationStack.Push(binaryOperator(val1, val2, gte))
		case types.And:
			s.evaluationStack.Push(binaryOperator(val1, val2, and))
		case types.Or:
			s.evaluationStack.Push(binaryOperator(val1, val2, or))
		case types.Min:
			s.evaluationStack.Push(binaryOperator(val1, val2, min))
		case types.Max:
			s.evaluationStack.Push(binaryOperator(val1, val2, max))
		}

	}
	s.currentAddress.Increment()
}

func binaryOperator(x, y types.NumericVal, f func(x, y float64) float64) types.NumericVal {
	res := f(x.AsFloat(), y.AsFloat())
	if x.IsFloat() || y.IsFloat() {
		return types.FloatVal(res)
	}
	return types.IntVal(res)
}

func unaryOperator(x types.NumericVal, f func(x float64) float64) types.NumericVal {
	res := f(x.AsFloat())
	if x.IsFloat() {
		return types.FloatVal(res)
	}
	return types.IntVal(res)
}
func mult(x float64, y float64) float64 {
	return x * y
}
func div(x float64, y float64) float64 {
	return x / y
}
func add(x, y float64) float64 {
	return x + y
}
func sub(x, y float64) float64 {
	return x - y
}
func mod(x, y float64) float64 {
	return float64(int(x) % int(y))
}
func eq(x, y float64) float64 {
	if x == y {
		return 1
	}
	return 0
}
func neq(x, y float64) float64 {
	if x != y {
		return 1
	}
	return 0
}
func lt(x, y float64) float64 {
	if x < y {
		return 1
	}
	return 0
}
func lte(x, y float64) float64 {
	if x <= y {
		return 1
	}
	return 0
}
func gt(x, y float64) float64 {
	if x > y {
		return 1
	}
	return 0
}

func gte(x, y float64) float64 {
	if x >= y {
		return 1
	}
	return 0
}

func and(x, y float64) float64 {
	var boolX, boolY bool
	if x != 0 {
		boolX = true
	}
	if y != 0 {
		boolY = true
	}
	if boolX && boolY {
		return 1
	}
	return 0
}
func or(x, y float64) float64 {
	var boolX, boolY bool
	if x != 0 {
		boolX = true
	}
	if y != 0 {
		boolY = true
	}
	if boolX || boolY {
		return 1
	}
	return 0
}
func not(x float64) float64 {
	if x == 0 {
		return 0
	}
	return 1
}
func min(x, y float64) float64 {
	return math.Min(x, y)
}
func max(x, y float64) float64 {
	return math.Max(x, y)
}
func negate(x float64) float64 {
	return x * -1
}
