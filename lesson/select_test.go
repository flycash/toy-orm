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
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Fatal(db)
	}
	testCases := []struct {
		name     string
		q        QueryBuilder
		wantSQL  string
		wantArgs []any
		wantErr  error
	}{
		{
			// From 都不调用
			name:    "no from",
			q:       NewSelector[TestModel](db),
			wantSQL: "SELECT * FROM `test_model`;",
		},
		{
			// 调用 FROM，但是传入空字符串
			name:    "empty from",
			q:       NewSelector[TestModel](db).From(""),
			wantSQL: "SELECT * FROM `test_model`;",
		},
		{
			// 调用 FROM
			name:    "with from",
			q:       NewSelector[TestModel](db).From("`test_model`"),
			wantSQL: "SELECT * FROM `test_model`;",
		},
		{
			// 调用 FROM，传入带 db 的
			name:    "from db.tbl",
			q:       NewSelector[TestModel](db).From("test_db.test_model"),
			wantSQL: "SELECT * FROM test_db.test_model;",
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q: NewSelector[TestModel](db).From("test_db.test_model").
				Where(C("Id").EQ(1)),
			wantSQL:  "SELECT * FROM test_db.test_model WHERE `id` = ?;",
			wantArgs: []any{1},
		},

		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel](db).From("test_db.test_model").
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantSQL:  "SELECT * FROM test_db.test_model WHERE (`age` > ?) AND (`age` < ?);",
			wantArgs: []any{18, 35},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[TestModel](db).From("test_db.test_model").
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantSQL:  "SELECT * FROM test_db.test_model WHERE (`age` > ?) AND (`age` < ?);",
			wantArgs: []any{18, 35},
		},
		{
			// 使用 OR
			name: "or",
			q: NewSelector[TestModel](db).From("test_db.test_model").
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantSQL:  "SELECT * FROM test_db.test_model WHERE (`age` > ?) OR (`age` < ?);",
			wantArgs: []any{18, 35},
		},
		{
			// 使用 NOT
			name: "not",
			q: NewSelector[TestModel](db).From("test_db.test_model").
				Where(Not(C("Age").GT(18))),
			wantSQL:  "SELECT * FROM test_db.test_model WHERE  NOT (`age` > ?);",
			wantArgs: []any{18},
		},
		{
			// 使用非法列名
			name: "invalid column",
			q: NewSelector[TestModel](db).From("test_db.test_model").
				Where(Not(C("invalid_column").GT(18))),
			wantErr: errors.New("toy-orm: 非法列名 invalid_column"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSQL, q.SQL)
			assert.Equal(t, tc.wantArgs, q.Args)
		})
	}
}
