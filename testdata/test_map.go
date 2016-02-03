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
	self := &TestMap{}
	GetConfigFile().AddInterfaceConfig(self, DBusObject2InterfaceConfig(self))
}

func (*TestMap) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DestName,
		ObjPath,
		BaseIfcName + ".TestMap",
	}
}

type TestMap struct {
	Byte       map[string]byte
	Bool       map[string]bool
	Int16      map[string]int16
	Int32      map[string]int32
	Int64      map[string]int64
	Float64    map[string]float64
	String     map[string]string
	Signature  map[string]dbus.Signature
	ObjectPath map[string]dbus.ObjectPath
	Variant    map[string]dbus.Variant

	Map    map[dbus.ObjectPath]bool
	MapMap map[string]map[string]bool
}
