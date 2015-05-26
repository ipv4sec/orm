// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package test

import (
	"os"

	"github.com/issue9/assert"
	"github.com/issue9/conv"
	"github.com/issue9/orm"
	"github.com/issue9/orm/dialect"
	"github.com/issue9/orm/fetch"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	driver = "mysql"

	prefix = "prefix_"
	dsn    string
	d      orm.Dialect
)

// 删除数据库文件。
func closeDB(a *assert.Assertion) {
	if driver != "sqlite3" {
		return
	}

	if _, err := os.Stat(dsn); err == nil || os.IsExist(err) {
		a.NotError(os.Remove(dsn))
	}
}

func newDB(a *assert.Assertion) *orm.DB {
	switch driver {
	case "mysql":
		dsn = "root@/orm_bench?charset=utf8"
		d = &dialect.Mysql{}
	case "sqlite3":
		dsn = "./test.db"
		d = &dialect.Sqlite3{}
	case "postgres":
		dsn = "" // TODO
		d = &dialect.Postgres{}
	default:
		panic("仅支持mysql,sqlite3,postgres三种数据库测试")
	}

	db, err := orm.NewDB(driver, dsn, prefix, d)
	a.NotError(err).NotNil(db)
	return db
}

// table表中是否存在size条记录，若不是，则触发error
func hasCount(db *orm.DB, a *assert.Assertion, table string, size int) {
	rows, err := db.Query(true, "SELECT COUNT(*) as cnt FROM #"+table)
	a.NotError(err).NotNil(rows)
	defer func() {
		a.NotError(rows.Close())
	}()

	data, err := fetch.Map(true, rows)
	a.NotError(err).NotNil(data)
	a.Equal(conv.MustInt(data[0]["cnt"], -1), size)

}