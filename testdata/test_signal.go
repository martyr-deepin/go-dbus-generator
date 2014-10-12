package main

import "pkg.linuxdeepin.com/lib/dbus"

func init() {
	self := &TestSignal{}
	GetConfigFile().AddInterfaceConfig(self, DBusObject2InterfaceConfig(self))
}

func (*TestSignal) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DestName,
		ObjPath,
		BaseIfcName + ".TestSignal",
	}
}

type TestSignal struct {
	Byte func(byte)

	//Bool           func(bool) //TODO: fix FixBool
	FixBool func(bool)

	Int16          func(int16)
	Int32          func(int32)
	Int64          func(int64)
	Float64        func(float64)
	String         func(string)
	Signature      func(dbus.Signature)
	ObjectPathfunc (dbus.ObjectPath)
	Variant        func(dbus.Variant)

	MultiIn func(byte, bool, int16)
	Map     func(map[byte]byte)
}
