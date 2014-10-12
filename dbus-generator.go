package main

import "path"
import "encoding/xml"
import "os"
import "flag"
import "os/exec"

//TODO: remove pkg.linuxdeepin.com/lib/dbus
import "pkg.linuxdeepin.com/lib/dbus"

func GetInterfaceInfo(inputdir string, ifc _Interface) dbus.InterfaceInfo {
	inFile := path.Join(inputdir, ifc.XMLFile)
	_, err := os.Stat(inFile)
	if err != nil {
		panic("open info file failed")
	}

	reader, err := os.Open(inFile)
	if err != nil {
		panic(err.Error() + "(File:" + inFile + ")")
	}

	decoder := xml.NewDecoder(reader)
	obj := dbus.NodeInfo{}
	decoder.Decode(&obj)
	for _, ifcInfo := range obj.Interfaces {
		if ifcInfo.Name == ifc.Interface {
			return ifc.handlleBlackList(ifcInfo)
		}
	}
	reader.Close()
	panic("Not Found Interface " + ifc.Interface)
}

type _Interface struct {
	OutFile, XMLFile, Dest, ObjectPath, Interface, ObjectName, TestPath string
	BlackMethods                                                        []string
	BlackProperties                                                     []string
	BlackSignals                                                        []string
}

func (ifc _Interface) handlleBlackList(data dbus.InterfaceInfo) dbus.InterfaceInfo {
	for _, name := range ifc.BlackMethods {
		var methods []dbus.MethodInfo
		for _, mInfo := range data.Methods {
			if mInfo.Name != name {
				methods = append(methods, mInfo)
			}
		}
		data.Methods = methods
	}
	for _, name := range ifc.BlackProperties {
		var properties []dbus.PropertyInfo
		for _, pInfo := range data.Properties {
			if pInfo.Name != name {
				properties = append(properties, pInfo)
			}
		}
		data.Properties = properties
	}
	for _, name := range ifc.BlackSignals {
		var signals []dbus.SignalInfo
		for _, sInfo := range data.Signals {
			if sInfo.Name != name {
				signals = append(signals, sInfo)
			}
		}
		data.Signals = signals
	}
	return data
}

func formatCode(infos *Infos) {
	switch BindingTarget(infos.Config.Target) {
	case GoLang:
		exec.Command("gofmt", "-w", infos.OutputDir()).Run()
	}
}

func geneateInit(infos *Infos) {
	for _, ifc := range infos.ListInterfaces() {
		writer, newOne := infos.GetWriter(ifc.OutFile)
		if newOne {
			renderInterfaceInit(writer, GetTemplate(BindingTarget(infos.Target()), TemplateTypeInit), infos)
		}
	}

	if BindingTarget(infos.Target()) == QML {
		renderQMLProject(infos.OutputDir(), infos)
	}
}

func generateMain(infos *Infos) {
	target := BindingTarget(infos.Target())

	switch target {
	case GoLang:
		writer, _ := infos.GetWriter("init")
		renderMain(writer, GetTemplate(target, TemplateTypeGlobal), infos)
	case PyQt:
		writer, _ := infos.GetWriter("__init__")
		renderMain(writer, GetTemplate(target, TemplateTypeGlobal), infos)
	case QML:
		writer, _ := infos.GetWriter("plugin.h")
		renderMain(writer, GetTemplate(target, TemplateTypeGlobal), infos)
	}

	for _, ifc := range infos.ListInterfaces() {
		writer, _ := infos.GetWriter(ifc.OutFile)

		inFile := path.Join(infos.InputDir(), ifc.XMLFile)
		_, err := os.Stat(inFile)
		if err != nil {
			panic(inFile + " dind't exists")
			//TODO: try find automatically the interface information
		}

		info := GetInterfaceInfo(infos.InputDir(), ifc)
		renderInterface(target, info, writer, ifc.Interface, ifc.ObjectName, infos)
	}
}

func main() {
	var outputPath, inputFile, target string
	flag.StringVar(&outputPath, "out", "out", "the directory to save the generated code")
	flag.StringVar(&inputFile, "in", "dbus.in.json", "the config file path")
	flag.StringVar(&target, "target", "", "now support QML/PyQt/GoLang target")
	flag.Parse()

	infos := loadInfos(inputFile)
	infos.normalize(outputPath, target)

	geneateInit(infos)
	generateMain(infos)

	formatCode(infos)
}
