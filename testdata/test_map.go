package main

import "pkg.linuxdeepin.com/lib/dbus"

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
