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
