package main

import "pkg.linuxdeepin.com/lib/dbus"

func init() {
	self := &TestArray{}
	GetConfigFile().AddInterfaceConfig(self, DBusObject2InterfaceConfig(self))
}

func (*TestArray) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DestName,
		ObjPath,
		BaseIfcName + ".TestArray",
	}
}

type TestArray struct {
	Byte       []byte
	Bool       []bool
	Int16      []int16
	Int32      []int32
	Int64      []int64
	Float64    []float64
	String     []string
	Signature  []dbus.Signature
	ObjectPath []dbus.ObjectPath
	Variant    []dbus.Variant

	Array [][]bool
	//Map   []map[bool]bool
}
