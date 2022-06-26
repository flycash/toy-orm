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
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	type User struct {
		Id        int64
		FirstName string
		Ctime     uint64
	}
	n := uint64(1000)
	u := &User{
		Id:        12,
		FirstName: "Tom",
		Ctime:     n,
	}
	u1 := &User{
		Id:        13,
		FirstName: "Jerry",
		Ctime:     n,
	}
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
		builder  QueryBuilder
		wantArgs []interface{}
		wantSql  string
		wantErr  error
	}{
		{
			name:    "no examples of values",
			builder: NewInserter[User](db).Values(),
			wantErr: errors.New("toy-orm: 插入0行"),
		},
		{
			name:     "single example of values",
			builder:  NewInserter[User](db).Values(u),
			wantSql:  "INSERT INTO `user`(`id`,`first_name`,`ctime`) VALUES(?,?,?);",
			wantArgs: []interface{}{int64(12), "Tom", n},
		},

		{
			name:     "multiple values of same type",
			builder:  NewInserter[User](db).Values(u, u1),
			wantSql:  "INSERT INTO `user`(`id`,`first_name`,`ctime`) VALUES(?,?,?),(?,?,?);",
			wantArgs: []interface{}{int64(12), "Tom", n, int64(13), "Jerry", n},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(tc.name, func(t *testing.T) {
			q, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, c.wantSql, q.SQL)
			assert.Equal(t, c.wantArgs, q.Args)
		})
	}
}

func TestInserter_Exec(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := newDB(mockDB)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec("INSERT .*").WillReturnResult(sqlmock.NewResult(1, 2))

	res := NewInserter[TestModel](db).Values(&TestModel{FirstName: "Tom"}).
		Exec(context.Background())
	id, err := res.LastInsertId()
	
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, id > 0)
}
