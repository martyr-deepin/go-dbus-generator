package main

import "pkg.linuxdeepin.com/lib/dbus"

func init() {
	self := &TestMethodCall{}
	GetConfigFile().AddInterfaceConfig(self, DBusObject2InterfaceConfig(self))
}
func (*TestMethodCall) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DestName,
		ObjPath,
		BaseIfcName + ".TestMethodCall",
	}
}

type testParamType map[string][]map[string]int32

type TestMethodCall struct {
}

func (*TestMethodCall) Test() {
}

func (*TestMethodCall) TestReturn() testParamType {
	return nil
}

func (*TestMethodCall) TestIn(a, b, c testParamType) {
}

func (*TestMethodCall) TestMultiReturn() (testParamType, testParamType) {
	return nil, nil
}

func (*TestMethodCall) TestError(a, b, c testParamType) error {
	return nil
}
func (*TestMethodCall) TestErrorWithMultiReturn(a, b, c testParamType) (testParamType, error) {
	return nil, nil
}
