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

//go:build demo
// +build demo

package lesson

import (
	"context"
	"fmt"
	"testing"
)

func TestOrm(t *testing.T) {
	var orm Orm
	tm := &TestModel{Id: 1, FirstName: "Deng"}
	// SELECT * FROM `test_model` WHERE id = 1
	err := orm.Get(context.Background(), &tm)
	if err != nil {
		t.Fatal(err)
	}

	// SELECT * FROM `test_model` WHERE `first_name` = 'Deng'
	orm.GetMulti(context.Background(), &TestModel{FirstName: "Deng"})
	// SELECT * FROM `test_model` WHERE `age` > ?
	orm.GetByWhere(context.Background(), "`age` > ?", 18)
}

func TestSelector(t *testing.T) {
	var s Selector[TestModel]
	// SELECT * FROM `test_model` WHERE id = 1
	tm, err := s.From("test_model").Where("id = ?", 1).Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tm)
}

func TestQuery(t *testing.T) {
	var q Query[TestModel]
	tm, err := q.From("test_model").Where("id = ?", 1).Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tm)

	var qi Query[TestModel]
	qi.Values(&TestModel{Id: 1}, &TestModel{Id: 2})
}
