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

import (
	"context"
	"fmt"
	"strings"
)

type Selector[T any] struct {
	db   *DB
	sb   strings.Builder
	args []any
	mi   *ModelInfo

	tbl   string
	where []Predicate
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		t   T
		err error
	)
	s.mi, err = s.db.r.get(&t)
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT * FROM ")
	if s.tbl == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.mi.tableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.tbl)
	}

	// 构造 WHERE
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}

	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column:
		fi, ok := s.mi.fieldMap[exp.name]
		if !ok {
			return fmt.Errorf("toy-orm: 非法列名 %s", exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fi.columnName)
		s.sb.WriteByte('`')
	case value:
		s.addArg(exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	default:
		return fmt.Errorf("toy-web: 不支持的表达式 %v", exp)
	}
	return nil
}

func (s *Selector[T]) addArg(val interface{}) {
	s.sb.WriteByte('?')
	s.args = append(s.args, val)
}

func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.tbl = tbl
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	// TODO implement me
	panic("implement me")
}
