// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"database/sql"
	"errors"

	"github.com/issue9/orm/core"
)

// InsertStmt 表示插入操作的 SQL 语句
type InsertStmt struct {
	table string
	cols  []string
	args  [][]interface{}
}

// Insert 声明一条插入语句
func Insert(table string) *InsertStmt {
	return &InsertStmt{
		table: table,
		args:  make([][]interface{}, 0, 10),
	}
}

// Columns 指定插入的列，多次指定，之前的会被覆盖。
func (stmt *InsertStmt) Columns(cols ...string) *InsertStmt {
	stmt.cols = cols
	return stmt
}

// Values 指定需要插入的值
func (stmt *InsertStmt) Values(vals ...interface{}) *InsertStmt {
	stmt.args = append(stmt.args, vals)
	return stmt
}

// Reset 重置语句
func (stmt *InsertStmt) Reset() {
	stmt.table = ""
	stmt.cols = stmt.cols[:0]
	stmt.args = stmt.args[:0]
}

// SQL 获取 SQL 的语句及参数部分
func (stmt *InsertStmt) SQL() (string, []interface{}, error) {
	if stmt.table == "" {
		return "", nil, ErrTableIsEmpty
	}

	if len(stmt.cols) == 0 {
		return "", nil, ErrColumnsIsEmpty
	}

	if len(stmt.args) == 0 {
		return "", nil, ErrValueIsEmpty
	}

	for _, vals := range stmt.args {
		if len(vals) != len(stmt.cols) {
			return "", nil, errors.New("数据与列数量不相同")
		}
	}

	buffer := core.NewStringBuilder("INSERT INTO ")
	buffer.WriteString(stmt.table)

	buffer.WriteByte('(')
	for _, col := range stmt.cols {
		buffer.WriteString(col)
		buffer.WriteByte(',')
	}
	buffer.TruncateLast(1)
	buffer.WriteByte(')')

	args := make([]interface{}, 0, len(stmt.cols)*len(stmt.args))
	buffer.WriteString(" VALUES ")
	for _, vals := range stmt.args {
		buffer.WriteByte('(')
		for _, v := range vals {
			if named, ok := v.(sql.NamedArg); ok && named.Name != "" {
				buffer.WriteByte('@')
				buffer.WriteString(named.Name)
			} else {
				buffer.WriteByte('?')
			}
			buffer.WriteByte(',')
			args = append(args, v)
		}
		buffer.TruncateLast(1) // 去掉最后的逗号
		buffer.WriteByte(')')
	}

	return buffer.String(), args, nil
}
