package main

import "testing"
import "./introspect"

func TestSigTable(t *testing.T) {
	data := map[string]string{
		"(nnqq)":    "dbus::types::BaseStruct<short, short, ushort, ushort >",
		"(ss(iy)o)": "dbus::types::BaseStruct<QString, QString, dbus::types::BaseStruct<int, uchar >, QDBusObjectPath >",
		"a(nq)":     "QList<dbus::types::BaseStruct<short, ushort > >",
		"a{nq}":     "QMap<short, ushort >",
	}
	table := introspect.NewTypesTable(QtTypeConvert{})
	for k, v := range data {
		if table.Get(k) != v {
			t.Fatalf("%q conveted to %q\n", k, table.Get(k))
		}
	}
}
