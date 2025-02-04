package runtime

import (
	"math"

	"github.com/awwithro/goink/pkg/parser/types"
	log "github.com/sirupsen/logrus"
)

func (s *Story) VisitOperator(op types.Operator) {
	log.Debugf("Visiting Operator: %d", op)
	if op.IsUnary() {
		val := mustPopStack[any](s.evaluationStack)
		log.Debugf("Operating on %T: %v", val, val)
		switch v := val.(type) {
		case types.NumericVal:
			switch op {
			case types.Negate:
				s.evaluationStack.Push(unaryOperator(v, negate))
			case types.Not:
				s.evaluationStack.Push(unaryOperator(v, not))
			case types.Int:
				s.evaluationStack.Push(types.IntVal(v.AsInt()))
			case types.Float:
				s.evaluationStack.Push(types.FloatVal(v.AsFloat()))
			case types.Floor:
				s.evaluationStack.Push(types.IntVal(int(math.Floor(v.AsFloat()))))
			default:
				s.Panicf("Unimplemented Operator: %d for %T", op, val)
			}
		case types.ListVal:
			switch op {
			case types.ListMin:
				s.evaluationStack.Push(v.Min())
			case types.ListMax:
				s.evaluationStack.Push(v.Max())
			case types.ListCount:
				s.evaluationStack.Push(types.IntVal(v.Count()))
			case types.ListRandom:
				s.evaluationStack.Push(v.Random())
			default:
				s.Panicf("Unimplemented Operator: %d for %T", op, val)
			}
		case *types.ListValItem:
			switch op {
			case types.ListMin:
				s.evaluationStack.Push(types.StringVal(v.Name))
			case types.ListValue:
				s.evaluationStack.Push(types.IntVal(v.Value))
			default:
				s.Panicf("Unimplemented Operator: %d for %T", op, val)
			}
		default:
			s.Panicf("no unary operation implemented for %T", val)
		}

	} else {
		val2 := mustPopStack[any](s.evaluationStack)
		val1 := mustPopStack[any](s.evaluationStack)
		log.Debug("Operating on ", val1, val2)
		switch v1 := val1.(type) {
		case types.NumericVal:
			v2, ok := val2.(types.NumericVal)
			if !ok {
				panicInvalidStackType[types.NumericVal](val2, s)
			}
			switch op {
			case types.Plus:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, add))
			case types.Minus:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, sub))
			case types.Multiply:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, mult))
			case types.Divide:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, div))
			case types.Modulus:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, mod))
			// Can Min Max work on lists?
			case types.Min:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, min))
			case types.Max:
				s.evaluationStack.Push(binaryNumericOperator(v1, v2, max))
			case types.Equal:
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, eq))
			case types.NotEqual:
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, neq))
			case types.LessThan:
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, lt))
			case types.LessThanEqual:
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, lte))
			case types.GreaterThan:
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, gt))
			case types.GreaterThanEqual:
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, gte))
			case types.And:
				s.evaluationStack.Push(binaryBoolOperator(v1, v2, and))
			case types.Or:
				s.evaluationStack.Push(binaryBoolOperator(v1, v2, or))
			default:
				s.Panicf("Unimplemented Operator: %d", op)
			}

		case types.Truthy:
			v2, ok := val2.(types.Truthy)
			if !ok {
				panicInvalidStackType[types.Truthy](val2, s)
			}
			switch op {
			case types.And:
				s.evaluationStack.Push(binaryBoolOperator(v1, v2, and))
			case types.Or:
				s.evaluationStack.Push(binaryBoolOperator(v1, v2, or))
			default:
				s.Panicf("Unimplemented Operator: %d", op)
			}

		case *types.ListValItem:
			switch op {
			case types.Equal:
				v2 := val2.(*types.ListValItem)
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, eq))
			case types.NotEqual:
				v2 := val2.(*types.ListValItem)
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, neq))
			case types.LessThan:
				v2 := val2.(*types.ListValItem)
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, lt))
			case types.LessThanEqual:
				v2 := val2.(*types.ListValItem)
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, lte))
			case types.GreaterThan:
				v2 := val2.(*types.ListValItem)
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, gt))
			case types.GreaterThanEqual:
				v2 := val2.(*types.ListValItem)
				s.evaluationStack.Push(binaryComparableOperator(v1, v2, gte))
			case types.Plus:
				v2 := val2.(types.IntVal)
				for range v2.AsInt() {
					v1 = v1.Next
				}
				s.evaluationStack.Push(v1)
			default:
				s.Panicf("Unimplemented Operator: %d", op)
				// case Add, Sub
				// Need to use the add operator to increment the list by a number
			}
		// Odd case, a string and int are used by listInt to get the position from a list
		// don't know why a VAR? operator isn't used. NEEDS TO GET THE ORIGINAL GLOBAL DEF
		case types.StringVal:
			v2, ok := val2.(types.IntVal)
			if !ok {
				panicInvalidStackType[types.IntVal](val2, s)
			}
			switch op {
			case types.ListInt:
				lst := s.computedLists[v1.String()]
				s.evaluationStack.Push(lst.AsList()[v2.AsInt()-1]) // Not Zero Indexed!

			default:
				s.Panicf("Unimplemented Operator: %d for %T and %T", op, val1, val2)
			}

		default:
			s.Panicf("Unimplemented Type: %T for Operation: %d", val1, op)
		}

	}
	s.currentAddress.Increment()
}

func binaryNumericOperator(x, y types.NumericVal, f func(x, y float64) float64) types.NumericVal {
	res := f(x.AsFloat(), y.AsFloat())
	if x.IsFloat() || y.IsFloat() {
		return types.FloatVal(res)
	}
	return types.IntVal(res)
}

func binaryComparableOperator[T any](x, y types.Comparable[T], f func(x, y types.Comparable[T]) bool) types.BoolVal {
	res := f(x, y)
	return types.BoolVal(res)
}

func binaryBoolOperator(x, y types.Truthy, f func(x, y bool) bool) types.BoolVal {
	res := f(x.AsBool(), y.AsBool())
	return types.BoolVal(res)
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
func eq[T any](x, y types.Comparable[T]) bool {
	return x.Equals(y.(T))
}
func neq[T any](x, y types.Comparable[T]) bool {
	return x.NotEquals(y.(T))
}
func lt[T any](x, y types.Comparable[T]) bool {
	return x.LT(y.(T))
}
func lte[T any](x, y types.Comparable[T]) bool {
	return x.LTE(y.(T))
}

func gt[T any](x, y types.Comparable[T]) bool {
	return x.GT(y.(T))
}

func gte[T any](x, y types.Comparable[T]) bool {
	return x.GTE(y.(T))
}

func and(x, y bool) bool {
	return x && y
}
func or(x, y bool) bool {
	return x || y
}
func not(x float64) float64 {
	if x == 0 {
		return 1
	}
	return 0
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
