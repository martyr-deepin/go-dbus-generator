package main

import "path"
import "os"
import "flag"
import "fmt"
import "os/exec"

import "pkg.linuxdeepin.com/lib/dbus/introspect"

func GetInterfaceInfo(inputdir string, ifc _Interface) introspect.InterfaceInfo {
	inFile := path.Join(inputdir, ifc.XMLFile)
	_, err := os.Stat(inFile)
	if err != nil {
		panic("open info file failed")
	}

	reader, err := os.Open(inFile)
	defer reader.Close()
	if err != nil {
		panic(err.Error() + "(File:" + inFile + ")")
	}

	obj, err := introspect.Parse(reader)
	if err != nil {
		panic(err.Error() + "(File:" + inFile + ")")
	}

	for _, ifcInfo := range obj.Interfaces {
		if ifcInfo.Name == ifc.Interface {
			return ifc.handlleBlackList(ifcInfo)
		}
	}
	panic("Not Found Interface " + ifc.Interface)
}

type _Interface struct {
	OutFile, XMLFile, Dest, ObjectPath, Interface, ObjectName, TestPath string
	BlackMethods                                                        []string
	BlackProperties                                                     []string
	BlackSignals                                                        []string
}

func (ifc _Interface) handlleBlackList(data introspect.InterfaceInfo) introspect.InterfaceInfo {
	for _, name := range ifc.BlackMethods {
		var methods []introspect.MethodInfo
		for _, mInfo := range data.Methods {
			if mInfo.Name != name {
				methods = append(methods, mInfo)
			}
		}
		data.Methods = methods
	}
	for _, name := range ifc.BlackProperties {
		var properties []introspect.PropertyInfo
		for _, pInfo := range data.Properties {
			if pInfo.Name != name {
				properties = append(properties, pInfo)
			}
		}
		data.Properties = properties
	}
	for _, name := range ifc.BlackSignals {
		var signals []introspect.SignalInfo
		for _, sInfo := range data.Signals {
			if sInfo.Name != name {
				signals = append(signals, sInfo)
			}
		}
		data.Signals = signals
	}
	return data
}

func renderedEnd(infos *Infos) {
	switch infos.Target() {
	case GoLang:
		exec.Command("gofmt", "-w", infos.OutputDir()).Run()
	case QML:
		renderQMLProject(infos.OutputDir(), infos)
	}
}

func geneateInit(infos *Infos) {
	for _, ifc := range infos.ListInterfaces() {
		writer, newOne := infos.GetWriter(ifc.OutFile)
		if newOne {
			renderInterfaceInit(writer, GetTemplate(infos.Target(), TemplateTypeInit), infos)
		}
	}

}

func generateMain(infos *Infos) {
	switch infos.Target() {
	case GoLang:
		writer, _ := infos.GetWriter("init")
		renderMain(writer, GetTemplate(GoLang, TemplateTypeGlobal), infos)
	case PyQt:
		writer, _ := infos.GetWriter("__init__")
		renderMain(writer, GetTemplate(PyQt, TemplateTypeGlobal), infos)
	case QML:
		writer, _ := infos.GetWriter("plugin.h")
		renderMain(writer, GetTemplate(QML, TemplateTypeGlobal), infos)
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
		renderInterface(infos.Target(), info, writer, ifc.Interface, ifc.ObjectName, infos)
	}
}

func main() {
	var outputPath, inputFile, target string
	flag.StringVar(&outputPath, "out", "", "the directory to save the generated code")
	flag.StringVar(&inputFile, "in", "dbus.in.json", "the config file path")
	flag.StringVar(&target, "target", "", "now support QML/PyQt/GoLang target")
	flag.Parse()

	infos, err := LoadInfos(inputFile, outputPath, target)
	if err != nil {
		fmt.Printf("Can't load info file (%q), please check it. \nError:%s\n\n", inputFile, err.Error())
		flag.Usage()
		return
	}

	geneateInit(infos)
	generateMain(infos)

	renderedEnd(infos)
}
