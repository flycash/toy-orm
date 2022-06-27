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
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

type Inserter[T any] struct {
	sess   Session
	values []*T
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return &Query{}, errors.New("toy-orm: 插入0行")
	}
	meta, err := i.sess.registry().get(i.values[0])
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	sb.WriteString("INSERT INTO `")
	sb.WriteString(meta.tableName)
	sb.WriteString("`(")
	for index, fd := range meta.fields {
		if index > 0 {
			sb.WriteByte(',')
		}
		cm, _ := meta.fieldMap[fd]
		sb.WriteByte('`')
		sb.WriteString(cm.columnName)
		sb.WriteByte('`')
	}
	sb.WriteString(")")
	sb.WriteString(" VALUES")
	args := make([]any, 0, len(i.values)*len(meta.fields))
	for index, val := range i.values {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('(')
		refVal := reflect.ValueOf(val).Elem()
		for j, v := range meta.fields {
			if j > 0 {
				sb.WriteByte(',')
			}
			fdVal := refVal.FieldByName(v)
			sb.WriteByte('?')
			args = append(args, fdVal.Interface())
		}
		sb.WriteByte(')')
	}
	sb.WriteByte(';')
	return &Query{SQL: sb.String(), Args: args}, nil
}

func (i *Inserter[T]) Exec(ctx context.Context) sql.Result {
	q, err := i.Build()
	if err != nil {
		return Result{
			err: err,
		}
	}
	res, err := i.sess.exec(ctx, q.SQL, q.Args...)
	return Result{
		err: err,
		res: res,
	}
}

func NewInserter[T any](sess Session) *Inserter[T] {
	return &Inserter[T]{sess: sess}
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}
