// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sqlbuilder 提供一套通过字符串拼接来构成 SQL 语句的工具。
package sqlbuilder

import (
	"bytes"
	"errors"
)

var (
	// ErrTableIsEmpty 未指定表名，任何 SQL 语句中，
	// 若未指定表名时，会返回此错误
	ErrTableIsEmpty = errors.New("表名为空")

	// ErrValueIsEmpty 在 Update 和 Insert 语句中，
	// 若未指定任何值，则返回此错误
	ErrValueIsEmpty = errors.New("值为空")

	// ErrColumnsIsEmpty 在 Insert 和 Select 语句中，
	// 若未指定任何列表，则返回此错误
	ErrColumnsIsEmpty = errors.New("未指定列")

	// ErrDupColumn 在 Update 中可能存在重复设置的列名。
	ErrDupColumn = errors.New("重复的列名")

	// ErrArgsNotMatch 在生成的 SQL 语句中，传递的参数与语句的占位符数量不匹配。
	ErrArgsNotMatch = errors.New("列与值的数量不匹配")
)

// SQLBuilder 对 bytes.Buffer 的一个简单封装。
// 当 Write* 系列函数出错时，直接 panic。
type SQLBuilder bytes.Buffer

// New 声明一个新的 SQLBuilder 实例
func New(str string) *SQLBuilder {
	return (*SQLBuilder)(bytes.NewBufferString(str))
}

func (b *SQLBuilder) buffer() *bytes.Buffer {
	return (*bytes.Buffer)(b)
}

// WriteString 写入一字符串
func (b *SQLBuilder) WriteString(str string) *SQLBuilder {
	if _, err := b.buffer().WriteString(str); err != nil {
		panic(err)
	}

	return b
}

// WriteByte 写入一字符
func (b *SQLBuilder) WriteByte(c byte) *SQLBuilder {
	if err := b.buffer().WriteByte(c); err != nil {
		panic(err)
	}

	return b
}

// Reset 重置内容
func (b *SQLBuilder) Reset() *SQLBuilder {
	b.buffer().Reset()
	return b
}

// TruncateLast 去掉最后几个字符
func (b *SQLBuilder) TruncateLast(n int) *SQLBuilder {
	b.buffer().Truncate(b.Len() - n)

	return b
}

// String 获取表示的字符串
func (b *SQLBuilder) String() string {
	return b.buffer().String()
}

// Bytes 获取表示的字符串
func (b *SQLBuilder) Bytes() []byte {
	return b.buffer().Bytes()
}

// Len 获取长度
func (b *SQLBuilder) Len() int {
	return b.buffer().Len()
}
