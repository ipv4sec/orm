// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package test

type bench struct {
	ID   int    `orm:"name(id);ai"`
	Name string `orm:"name(name);len(20)"`
	Pass string `orm:"name(pass);len(32)"`
}

func (b *bench) Meta() string {
	return "name(bench)"
}

type user struct {
	ID       int    `orm:"name(id);ai;"`
	Username string `orm:"unique(unique_username);index(index_name);len(50)"`
	Password string `orm:"name(password);len(20)"`
	Regdate  int    `orm:"-"`
}

func (m *user) Meta() string {
	return "check(chk_name,id>0);engine(innodb);charset(utf-8);name(users)"
}

type userInfo struct {
	UID       int    `orm:"name(uid);pk"`
	FirstName string `orm:"name(firstName);unique(unique_name);len(20)"`
	LastName  string `orm:"name(lastName);unique(unique_name);len(20)"`
	Sex       string `orm:"name(sex);default(male);len(6)"`
}

func (m *userInfo) Meta() string {
	return "check(chk_name,uid>0);engine(innodb);charset(utf-8);name(user_info)"
}

type admin struct {
	user

	Email string `orm:"name(email);len(20);unique(unique_email)"`
	Group int    `orm:"name(group);"`
}

func (m *admin) Meta() string {
	return "check(chk_name,id>0);engine(innodb);charset(utf-8);name(administrators)"
}
