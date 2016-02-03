/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
// +build ignore

package main

import C "launchpad.net/gocheck"

func (*testWrap) TestSet(c *C.C) {
	infos := NewInfos()

	err := infos.SetBusType("session")
	c.Check(err, C.Equals, nil)
	c.Check(infos.BusType(), C.Equals, "session")

	err = infos.SetOutputDir("/dev/shm")
	c.Check(err, C.Equals, nil)
	c.Check(infos.OutputDir(), C.Equals, "/dev/shm")

	infos.SetTarget("gOlang")
	c.Check(err, C.Equals, nil)
	c.Check(infos.Target(), C.Equals, GoLang)

	err = infos.SetTarget("Olang")
	c.Check(err, C.NotNil)
}

func (*testWrap) TestUtils(c *C.C) {
	c.Check(lower("Deepin"), C.Equals, "deepin")
	c.Check(lower("D"), C.Equals, "d")
	c.Check(lower(""), C.Equals, "")

	c.Check(upper("deepin"), C.Equals, "Deepin")
	c.Check(upper("D"), C.Equals, "D")
	c.Check(upper("d"), C.Equals, "D")
	c.Check(upper(""), C.Equals, "")
}

func (*testWrap) TestLoadInfos(c *C.C) {
	infos := NewInfos()
	infos.SetOutputDir("/dev/shm")
	infos.SetBusType("session")
	infos.SetPackageName("dbus-generate.test")
	err := infos.normalize("out", "goLang")
	c.Check(err, C.Equals, nil)

	c.Check(infos.Target(), C.Equals, GoLang)
	c.Check(infos.BusType(), C.Equals, "Session")

	infos = NewInfos()
	infos.SetBusType("session")

	infos.SetTarget("golang")
	infos.normalize("", "")
	c.Check(infos.BusType(), C.Equals, "Session")

	infos.SetTarget("qml")
	infos.normalize("", "")
	c.Check(infos.BusType(), C.Equals, "session")

	infos.SetTarget("pyqt")
	infos.normalize("", "")
	c.Check(infos.BusType(), C.Equals, "session")
}
