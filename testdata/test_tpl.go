package main

import "pkg.linuxdeepin.com/lib/dbus"

func init() {
	self := &TestFoo{}
	GetConfigFile().AddInterfaceConfig(self, DBusObject2InterfaceConfig(self))
}

func (*TestFoo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DestName,
		ObjPath,
		BaseIfcName + ".TestFoo",
	}
}

type TestFoo struct {
	Foo string
	Bar func()
}

func (*TestFoo) Foobar() {
}
