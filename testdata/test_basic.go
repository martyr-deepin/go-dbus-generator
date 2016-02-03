/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import "pkg.deepin.io/lib/dbus"

func init() {
	self := &TestBasic{}
	GetConfigFile().AddInterfaceConfig(self, DBusObject2InterfaceConfig(self))
}

func (*TestBasic) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DestName,
		ObjPath,
		BaseIfcName + ".TestBasic",
	}
}

type TestBasic struct {
	Byte       byte
	Bool       bool
	Int16      int16
	Int32      int32
	Int64      int64
	Float64    float64
	String     string
	Signature  dbus.Signature
	ObjectPath dbus.ObjectPath
	Variant    dbus.Variant
}
