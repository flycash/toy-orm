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
)

type DBOption func(*DB)

type DB struct {
	db *sql.DB
	r  *registry
}

func NewDB(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return newDB(db, opts...)
}

func newDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		db: db,
		r:  &registry{},
	}
	for _, o := range opts {
		o(res)
	}
	return res, nil
}

func (db *DB) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		r:  db.r,
		tx: tx,
	}, nil
}

func (db *DB) query(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, sql, args...)
}

func (db *DB) exec(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, sql, args...)
}

func (db *DB) registry() *registry {
	return db.r
}

type Session interface {
	query(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	exec(ctx context.Context, sql string, args ...any) (sql.Result, error)
	registry() *registry
}
