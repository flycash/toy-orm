// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lesson

type op string

const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opAND = "AND"
	opOR  = "OR"
	opNOT = "NOT"
)

func (o op) String() string {
	return string(o)
}

// Expression 代表语句，或者语句的部分
// 暂时没想好怎么设计方法，所以直接做成标记接口
type Expression interface {
	expr()
}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return valueOf(exp)
	}
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (Predicate) expr() {}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAND,
		right: r,
	}
}

func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOR,
		right: right,
	}
}

// type Predicate struct {
// 	Column string
// 	Op     string
// 	Arg    any
// }
//
// func Eq(column string, arg any) Predicate {
// 	return Predicate{
// 		Column: column,
// 		Op:     "=",
// 		Arg:    arg,
// 	}
// }
