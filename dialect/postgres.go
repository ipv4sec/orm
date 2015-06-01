// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialect

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/issue9/orm/forward"
)

// 返回一个适配postgresql的forward.Dialect接口
func Postgres() forward.Dialect {
	return &postgres{}
}

type postgres struct{}

// implement forward.Dialect.SupportInsertMany()
func (p *postgres) SupportInsertMany() bool {
	return true
}

// implement forward.Dialect.QuoteTuple()
func (p *postgres) QuoteTuple() (byte, byte) {
	return '"', '"'
}

// implement forward.Dialect.Quote()
func (p *postgres) Quote(w *bytes.Buffer, name string) error {
	if err := w.WriteByte('"'); err != nil {
		return err
	}

	if _, err := w.WriteString(name); err != nil {
		return err
	}

	return w.WriteByte('"')
}

// implement forward.Dialect.ReplaceMarks()
// 在有?占位符的情况下，语句中不能包含$字符串
func (p *postgres) ReplaceMarks(sql *string) error {
	s := *sql
	if strings.IndexByte(s, '?') < 0 {
		return nil
	}

	num := 1
	ret := make([]rune, 0, len(s))
	for _, c := range s {
		switch c {
		case '?':
			ret = append(ret, '$')
			ret = append(ret, []rune(strconv.Itoa(num))...)
			num++
		case '$':
			return errors.New("语句中包含非法的字符串:$")
		default:
			ret = append(ret, c)
		}
	}

	*sql = string(ret)
	return nil
}

// implement forward.Dialect.LimitSQL()
func (p *postgres) LimitSQL(w *bytes.Buffer, limit int, offset ...int) ([]int, error) {
	return mysqlLimitSQL(w, limit, offset...)
}

// implement forward.Dialect.AIColSQL()
func (p *postgres) AIColSQL(w *bytes.Buffer, model *forward.Model) error {
	// Potgres的AI仅仅是类型不同，可直接使用NoAiColSQL输出
	return nil
}

// implement forward.Dialect.NoAIColSQL()
func (p *postgres) NoAIColSQL(w *bytes.Buffer, model *forward.Model) error {
	for _, col := range model.Cols {
		if err := createColSQL(p, w, col); err != nil {
			return err
		}
		w.WriteByte(',')
	}
	return nil
}

// implement forward.Dialect.ConstraintsSQL()
func (p *postgres) ConstraintsSQL(w *bytes.Buffer, model *forward.Model) error {
	if len(model.PK) > 0 {
		createPKSQL(p, w, model.PK, model.Name+pkName) // postgres主键名需要全局唯一？
		w.WriteByte(',')
	}

	createConstraints(p, w, model)
	return nil
}

// implement forward.Dialect.TruncateTableSQL()
func (p *postgres) TruncateTableSQL(w *bytes.Buffer, tableName, aiColumn string) error {
	w.WriteString("TRUNCATE TABLE ")
	w.WriteString(tableName)
	if len(aiColumn) == 0 {
		return nil
	}

	w.WriteString("; ALTER SEQUENCE ")
	w.WriteString(tableName)
	w.WriteByte('_')
	w.WriteString(aiColumn)
	_, err := w.WriteString("_seq RESTART WITH 1")
	return err
}

// implement base.sqlType
// 将col转换成sql类型，并写入buf中。
func (p *postgres) sqlType(buf *bytes.Buffer, col *forward.Column) error {
	if col == nil {
		return errors.New("sqlType:col参数是个空值")
	}

	if col.GoType == nil {
		return errors.New("sqlType:无效的col.GoType值")
	}

	switch col.GoType.Kind() {
	case reflect.Bool:
		buf.WriteString("BOOLEAN")
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
		if col.IsAI() {
			buf.WriteString("SERIAL")
		} else {
			buf.WriteString("SMALLINT")
		}
	case reflect.Int32, reflect.Uint32:
		if col.IsAI() {
			buf.WriteString("SERIAL")
		} else {
			buf.WriteString("INT")
		}
	case reflect.Int64, reflect.Int, reflect.Uint64, reflect.Uint:
		if col.IsAI() {
			buf.WriteString("BIGSERIAL")
		} else {
			buf.WriteString("BIGINT")
		}
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("DOUBLE(%d,%d)", col.Len1, col.Len2))
	case reflect.String:
		if col.Len1 < 65533 {
			buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
		} else {
			buf.WriteString("TEXT")
		}
	case reflect.Slice, reflect.Array: // []rune,[]byte当作字符串处理
		k := col.GoType.Elem().Kind()
		if (k != reflect.Uint8) && (k != reflect.Int32) {
			return errors.New("sqlType:不支持数组类型")
		}

		if col.Len1 < 65533 {
			buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
		} else {
			buf.WriteString("TEXT")
		}
	case reflect.Struct:
		switch col.GoType {
		case nullBool:
			buf.WriteString("BOOLEAN")
		case nullFloat64:
			buf.WriteString(fmt.Sprintf("DOUBLE(%d,%d)", col.Len1, col.Len2))
		case nullInt64:
			if col.IsAI() {
				buf.WriteString("BIGSERIAL")
			} else {
				buf.WriteString("BIGINT")
			}
		case nullString:
			if col.Len1 < 65533 {
				buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
			}
			buf.WriteString("TEXT")
		case timeType:
			buf.WriteString("TIME")
		}
	default:
		return fmt.Errorf("sqlType:不支持的类型:[%v]", col.GoType.Name())
	}

	return nil
}
