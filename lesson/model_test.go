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

func Test_registry_register(t *testing.T) {
	testCases := []struct {
		name    string
		input   any
		wantMi  *ModelInfo
		wantErr error
	}{
		{
			// 纯粹的 struct
			name:    "struct",
			input:   TestModel{},
			wantErr: errors.New("toy-orm: 非法类型"),
		},
		{
			name:  "test_model",
			input: &TestModel{},
			wantMi: &ModelInfo{
				tableName: "test_model",
				fieldMap: map[string]*FieldInfo{
					"Id": {
						columnName: "id",
					},
					"FirstName": {
						columnName: "first_name",
					},
					"Age": {
						columnName: "age",
					},
					"LastName": {
						columnName: "last_name",
					},
				},
			},
		},
	}
	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, err := r.register(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantMi, mi)
		})
	}
}
