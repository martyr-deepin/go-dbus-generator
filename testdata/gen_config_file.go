package main

import "encoding/json"
import "encoding/xml"
import "log"
import "os"
import "path"
import "pkg.linuxdeepin.com/lib/dbus"
import "pkg.linuxdeepin.com/lib/dbus/introspect"
import "sort"
import "strings"

const (
	XMLDir = "xml"
)

type PrjConfig struct {
	InputDir string
	BusType  string
	DestName string
}

type InterfaceConfig struct {
	Interface  string
	OutFile    string
	XMLFile    string
	ObjectName string
}
type InterfaceConfigs []InterfaceConfig

func (ifs InterfaceConfigs) Len() int {
	return len(ifs)
}
func (ifs InterfaceConfigs) Swap(i, j int) {
	ifs[i], ifs[j] = ifs[j], ifs[i]
}
func (ifs InterfaceConfigs) Less(i, j int) bool {
	return ifs[i].Interface < ifs[j].Interface
}

type ConfigFile struct {
	Config     PrjConfig
	Interfaces []InterfaceConfig
	objs       map[InterfaceConfig]dbus.DBusObject
}

var __configFile__ *ConfigFile

func DBusObject2InterfaceConfig(obj dbus.DBusObject) InterfaceConfig {
	filed := strings.Split(obj.GetDBusInfo().Interface, ".")
	name := filed[len(filed)-1]
	return InterfaceConfig{
		Interface:  obj.GetDBusInfo().Interface,
		OutFile:    name,
		XMLFile:    name + ".xml",
		ObjectName: name,
	}
}

func (cfg *ConfigFile) AddInterfaceConfig(obj dbus.DBusObject, conf InterfaceConfig) {
	cfg.objs[conf] = obj
}
func (cfg *ConfigFile) SetPrjInfo(conf PrjConfig) {
	cfg.Config = conf
}

func GetConfigFile() *ConfigFile {
	if __configFile__ == nil {
		__configFile__ = &ConfigFile{
			objs: make(map[InterfaceConfig]dbus.DBusObject),
		}
	}
	return __configFile__
}

type XMLInterface struct {
	Interface *introspect.InterfaceInfo `xml:"interface"`
}

func (cfg *ConfigFile) Generate() {
	os.Mkdir(*OutputDir, os.ModePerm|os.ModeDir)

	for ifcInfo, obj := range cfg.objs {
		cfg.Interfaces = append(cfg.Interfaces, ifcInfo)
		xmlInfo := dbus.BuildInterfaceInfo(obj)
		xmlInfo.Name = ifcInfo.Interface
		f, err := os.Create(path.Join(*OutputDir, ifcInfo.XMLFile))
		if err != nil {
			log.Fatal("failed create xml file", ifcInfo.XMLFile, err)
		}

		type Interface *introspect.InterfaceInfo
		bytes, err := xml.MarshalIndent(introspect.NodeInfo{
			Interfaces: []introspect.InterfaceInfo{*xmlInfo},
		}, "", "  ")
		if err != nil {
			log.Fatal("failed marshl interface:", err)
		}

		f.Write(bytes)
		f.Close()
		log.Println("Generated...", path.Join(*OutputDir, ifcInfo.XMLFile))
	}

	//avoid random position otherwise it will bother git diff.
	sort.Sort(InterfaceConfigs(cfg.Interfaces))

	f, err := os.Create(path.Join(*OutputDir, "dbus.in.json"))
	if err != nil {
		log.Fatal("failed create dbus.in.json", err)
	}
	defer f.Close()
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatal("marsh ConfigFile failed:", err)
	}
	f.Write(bytes)
	log.Println("Generated...", path.Join(*OutputDir, "dbus.in.json"))

	cfg.Interfaces = nil
}
