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
	"database/sql"
)

type Orm interface {
	// Get 给 SELECT 语句用
	// Get 首先查找 val 里面的主键，如果主键存在，就用主键作为 WHERE 条件
	// 否则，查找唯一索引，如果唯一索引存在，就用唯一索引作为 WHERE 条件
	Get(ctx context.Context, val interface{}) error
	// GetMulti val里面字段如果有值，就会被用来构造 WHERE 条件
	GetMulti(ctx context.Context, val interface{}) ([]interface{}, error)
	// GetByWhere 根据查询条件来筛选
	GetByWhere(ctx context.Context, where string, args ...any) error

	// 不断加方法

	Update(ctx context.Context, val interface{}) (sql.Result, error)
	Insert(ctx context.Context, val interface{}) (sql.Result, error)
	Delete(ctx context.Context, val interface{}) (sql.Result, error)
}

type Query[T any] interface {
	From(tbl string) Query[T]
	Where(where string, args ...any) Query[T]
	OrderBy(o string) Query[T]

	// Get 根据前面查询结果来
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)

	Update(ctx context.Context) (sql.Result, error)

	Values(ts ...*T) Query[T]
	Insert(ctx context.Context) (sql.Result, error)
	OnConflict() Query[T]

	Delete(ctx context.Context) (sql.Result, error)
}

type Selector[T any] interface {
	From(tbl string) Selector[T]
	Where(where string, args ...any) Selector[T]
	OrderBy(o string) Selector[T]

	// 以下终结方法

	// Get 根据前面查询结果来
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}
