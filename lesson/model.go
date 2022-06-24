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
	"reflect"
	"sync"
	"unicode"
)

type ModelInfo struct {
	tableName string
	fieldMap  map[string]*FieldInfo
}

type FieldInfo struct {
	columnName string
}

type registry struct {
	models sync.Map
}

func (r *registry) register(val any) (*ModelInfo, error) {
	// 这里我们假设必然使用结构体指针
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("toy-orm: 非法类型")
	}
	typ = typ.Elem()

	numField := typ.NumField()
	fdInfos := make(map[string]*FieldInfo, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		cn := fd.Name
		fdInfos[cn] = &FieldInfo{
			columnName: underscoreName(cn),
		}
	}

	mi := &ModelInfo{
		tableName: underscoreName(typ.Name()),
		fieldMap:  fdInfos,
	}
	r.models.Store(reflect.TypeOf(val), mi)
	return mi, nil
}

func (r *registry) get(val any) (*ModelInfo, error) {
	typ := reflect.TypeOf(val)
	mi, ok := r.models.Load(typ)
	if ok {
		return mi.(*ModelInfo), nil
	}
	return r.register(val)
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

// func (d *DB) NewSelector[T any]() Selector[T] {
// 	return &Selector[T]{
// 		db: d,
// 	}
// }

// 例如 models[reflect.TypeOf(&TestModel{})]=&ModelInfo{}
// var models = map[reflect.Type]*ModelInfo{}
//
// var defaultRegistry = &registry{
// 	models: make(map[reflect.Type]*ModelInfo, 16),
// }
//
// type registry struct {
// 	models map[reflect.Type]*ModelInfo
// }
//
// func Register(val any) error {
// 	defaultRegistry.models[reflect.TypeOf(val)] = &ModelInfo{}
// 	return nil
// }
