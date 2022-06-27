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

type Tx struct {
	tx *sql.Tx
	r  *registry
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t *Tx) query(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, sql, args...)
}

func (t *Tx) exec(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, sql, args...)
}

func (t *Tx) registry() *registry {
	return t.r
}
