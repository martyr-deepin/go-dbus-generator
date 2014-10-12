package main

import "encoding/json"
import "encoding/xml"
import "fmt"
import "log"
import "os"
import "path"
import "pkg.linuxdeepin.com/lib/dbus"
import "strings"
import "time"

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

func (cfg *ConfigFile) Generate2() {
	con, err := dbus.SessionBus()
	if err != nil {
		log.Fatal("Can't connect to session bus")
	}
	os.Mkdir(*OutputDir, os.ModePerm|os.ModeDir)
	fmt.Println("HH")

	for ifcInfo, obj := range cfg.objs {
		cfg.Interfaces = append(cfg.Interfaces, ifcInfo)
		fmt.Println("1")
		switch strings.ToLower(cfg.Config.BusType) {
		case "session":
			err := dbus.InstallOnSession(obj)
			if err != nil {
				log.Fatal("failed install on session", obj.GetDBusInfo(), err)
			}
		case "system":
			err := dbus.InstallOnSystem(obj)
			log.Fatal("failed install on session", obj.GetDBusInfo(), err)
		default:
			log.Fatal("didn't support bus type", cfg.Config)
		}
		fmt.Println("2")

		var xml string
		<-time.After(time.Second * 1)
		dbusobj := con.Object(obj.GetDBusInfo().Dest, dbus.ObjectPath(obj.GetDBusInfo().ObjectPath))
		dbusobj.Call("org.freedesktop.DBus.Introspectable.Introspect", dbus.FlagNoAutoStart).Store(&xml)
		fmt.Println("3")
		f, err := os.Create(path.Join(*OutputDir, ifcInfo.XMLFile))
		if err != nil {
			log.Fatal("failed create xml file", ifcInfo.XMLFile, err)
		}
		f.WriteString(xml)
		f.Close()
		log.Println("Generated...", path.Join(*OutputDir, ifcInfo.XMLFile))

		//dbus.UnInstallObject(obj)
	}

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
}

type XMLInterface struct {
	Interface *dbus.InterfaceInfo `xml:"interface"`
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

		type Interface *dbus.InterfaceInfo
		bytes, err := xml.MarshalIndent(dbus.NodeInfo{
			Interfaces: []dbus.InterfaceInfo{*xmlInfo},
		}, "", "  ")
		if err != nil {
			log.Fatal("failed marshl interface:", err)
		}

		fmt.Println("HH:", string(bytes))

		f.Write(bytes)
		f.Close()
		log.Println("Generated...", path.Join(*OutputDir, ifcInfo.XMLFile))
	}

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
