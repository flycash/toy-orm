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
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := newDB(mockDB)
	if err != nil {
		t.Fatal(err)
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

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := newDB(mockDB)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			// 查询返回错误
			name:    "query error",
			mockErr: errors.New("invalid query"),
			wantErr: errors.New("invalid query"),
			query:   "SELECT .*",
		},
		{
			name:     "no row",
			wantErr:  errors.New("toy-orm: 未找到数据"),
			query:    "SELECT .*",
			mockRows: sqlmock.NewRows([]string{"id"}),
		},
		{
			name:    "too many column",
			wantErr: errors.New("toy-orm: 列过多"),
			query:   "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "extra_column"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"), []byte("nothing"))
				return res
			}(),
		},
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Da",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
		},
	}

	for _, tc := range testCases {
		exp := mock.ExpectQuery(tc.query)
		if tc.mockErr != nil {
			exp.WillReturnError(tc.mockErr)
		} else {
			exp.WillReturnRows(tc.mockRows)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}
