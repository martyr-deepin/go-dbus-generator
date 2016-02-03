/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import "testing"

import C "launchpad.net/gocheck"

type testWrap struct{}

func Test(t *testing.T) { C.TestingT(t) }
func init() {
	C.Suite(&testWrap{})
}

func (*testWrap) TestQtType(c *C.C) {
	c.Check(getQType("u"), C.Equals, "uint")
	c.Check(getQType("ah"), C.Equals, "QList<quint32 >")
	c.Check(getQType("au"), C.Equals, "QList<uint >")
	c.Check(getQType("ao"), C.Equals, "QList<QDBusObjectPath >")
	c.Check(getQType("as"), C.Equals, "QList<QString >")
	c.Check(getQType("av"), C.Equals, "QList<QDBusVariant >")
	c.Check(getQType("a{ss}"), C.Equals, "QMap<QString, QString >")
}
