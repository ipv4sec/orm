// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package model

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/orm/internal/modeltest"
)

// go1.10 BenchmarkNewModelNoCached-4   	  200000	      8161 ns/op
func BenchmarkNewModelNoCached(b *testing.B) {
	Clear()
	a := assert.New(b)

	for i := 0; i < b.N; i++ {
		m, err := New(&modeltest.User{})
		Clear()
		a.NotError(err).NotNil(m)
	}
}

// go1.10 BenchmarkNewModelCached-4     	10000000	       187 ns/op
func BenchmarkNewModelCached(b *testing.B) {
	Clear()
	a := assert.New(b)

	for i := 0; i < b.N; i++ {
		m, err := New(&modeltest.User{})
		a.NotError(err).NotNil(m)
	}
}
