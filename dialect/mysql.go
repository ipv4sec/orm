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

	"github.com/issue9/orm"
)

type Mysql struct{}

// implement orm.Dialect.QuoteTuple()
func (m *Mysql) QuoteTuple() (byte, byte) {
	return '`', '`'
}

// implement orm.Dialect.Quote
func (m *Mysql) Quote(w *bytes.Buffer, name string) error {
	if err := w.WriteByte('`'); err != nil {
		return err
	}

	if _, err := w.WriteString(name); err != nil {
		return err
	}

	return w.WriteByte('`')
}

// implement orm.Dialect.Limit()
func (m *Mysql) LimitSQL(w *bytes.Buffer, limit int, offset ...int) ([]int, error) {
	return mysqlLimitSQL(w, limit, offset...)
}

// implement orm.Dialect.CreateTableSQL()
func (m *Mysql) CreateTableSQL(model *orm.Model) (string, error) {
	buf := bytes.NewBufferString("CREATE TABLE IF NOT EXISTS ")
	buf.Grow(300)

	buf.WriteString(model.Name)
	buf.WriteByte('(')

	// 写入字段信息
	for _, col := range model.Cols {
		if err := createColSQL(m, buf, col); err != nil {
			return "", err
		}

		if col.IsAI() {
			buf.WriteString(" PRIMARY KEY AUTO_INCREMENT")
		}
		buf.WriteByte(',')
	}

	// PK，若有自增，则已经在上面指定
	if len(model.PK) > 0 && !model.PK[0].IsAI() {
		createPKSQL(m, buf, model.PK, pkName)
		buf.WriteByte(',')
	}

	createConstraints(m, buf, model)

	// key index不存在CONSTRAINT形式的语句
	if len(model.KeyIndexes) == 0 {
		for name, index := range model.KeyIndexes {
			buf.WriteString("INDEX ")
			buf.WriteString(name)
			buf.WriteByte('(')
			for _, col := range index {
				buf.WriteString(col.Name)
				buf.WriteByte(',')
			}
			buf.Truncate(buf.Len() - 1) // 去掉最后的逗号
			buf.WriteString("),")
		}
	}

	buf.Truncate(buf.Len() - 1) // 去掉最后的逗号
	buf.WriteByte(')')          // end CreateTable

	return buf.String(), nil
}

// implement orm.Dialect.TruncateTableSQL()
func (m *Mysql) TruncateTableSQL(tableName string) string {
	return "TRUNCATE TABLE " + tableName
}

// implement base.sqlType()
func (m *Mysql) sqlType(buf *bytes.Buffer, col *orm.Column) error {
	if col == nil {
		return errors.New("sqlType:col参数是个空值")
	}

	if col.GoType == nil {
		return errors.New("sqlType:无效的col.GoType值")
	}

	addIntLen := func() {
		if col.Len1 > 0 {
			buf.WriteByte('(')
			buf.WriteString(strconv.Itoa(col.Len1))
			buf.WriteByte(')')
		}
	}

	switch col.GoType.Kind() {
	case reflect.Bool:
		buf.WriteString("BOOLEAN")
	case reflect.Int8:
		buf.WriteString("SMALLINT")
		addIntLen()
	case reflect.Int16:
		buf.WriteString("MEDIUMINT")
		addIntLen()
	case reflect.Int32:
		buf.WriteString("INT")
		addIntLen()
	case reflect.Int64, reflect.Int: // reflect.Int大小未知，都当作是BIGINT处理
		buf.WriteString("BIGINT")
		addIntLen()
	case reflect.Uint8:
		buf.WriteString("SMALLINT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Uint16:
		buf.WriteString("MEDIUMINT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Uint32:
		buf.WriteString("INT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Uint64, reflect.Uint, reflect.Uintptr:
		buf.WriteString("BIGINT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("DOUBLE(%d,%d)", col.Len1, col.Len2))
	case reflect.String:
		if col.Len1 < 65533 {
			buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
		} else {
			buf.WriteString("LONGTEXT")
		}
	case reflect.Slice, reflect.Array: // []rune,[]byte当作字符串处理
		k := col.GoType.Elem().Kind()
		if (k != reflect.Uint8) && (k != reflect.Int32) {
			return fmt.Errorf("sqlType:不支持[%v]类型的数组", k)
		}

		if col.Len1 < 65533 {
			buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
		} else {
			buf.WriteString("LONGTEXT")
		}
	case reflect.Struct:
		switch col.GoType {
		case nullBool:
			buf.WriteString("BOOLEAN")
		case nullFloat64:
			buf.WriteString(fmt.Sprintf("DOUBLE(%d,%d)", col.Len1, col.Len2))
		case nullInt64:
			buf.WriteString("BIGINT")
			addIntLen()
		case nullString:
			if col.Len1 < 65533 {
				buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
			} else {
				buf.WriteString("LONGTEXT")
			}
		case timeType:
			buf.WriteString("DATETIME")
		}
	default:
		return fmt.Errorf("sqlType:不支持的类型:[%v]", col.GoType.Name())
	}

	return nil
}
