/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import (
	"flag"
)

const (
	DestName    = "dbus_generator.test"
	ObjPath     = "/dbus_generator/test"
	BaseIfcName = "dbus_generator.test"
)

var OutputDir = flag.String("output", "output", "the output directory for generated files")

func main() {
	flag.Parse()
	GetConfigFile().SetPrjInfo(PrjConfig{
		"testdata/" + *OutputDir,
		"Session",
		DestName,
	})
	GetConfigFile().Generate()
}
