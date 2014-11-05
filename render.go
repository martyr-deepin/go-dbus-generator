package main

import "text/template"
import "pkg.linuxdeepin.com/lib/dbus"
import "pkg.linuxdeepin.com/lib/dbus/introspect"
import "io"
import "log"
import "strings"

//TODO: try removing dependent on variable of infos

func renderMain(writer io.Writer, tpl string, infos *Infos) {
	template.Must(template.New("main").Funcs(template.FuncMap{
		"Lower":   lower,
		"Upper":   upper,
		"BusType": func() string { return infos.BusType() },
		"PkgName": func() string { return infos.PackageName() },
		"GetModules": func() map[string]string {
			r := make(map[string]string)
			for _, ifc := range infos.ListInterfaces() {
				r[ifc.OutFile] = ifc.OutFile
			}
			return r
		},
		"GetQtSignaturesType": func() map[string]string { return getQtSignaturesType(infos) },
	}).Parse(tpl)).Execute(writer, infos)
}

func renderInterfaceInit(writer io.Writer, tpl string, infos *Infos) {
	template.Must(template.New("IfcInit").Funcs(template.FuncMap{
		"BusType":    func() string { return infos.BusType() },
		"PkgName":    func() string { return infos.PackageName() },
		"HasSignals": func() bool { return true },
	}).Parse(tpl)).Execute(writer, nil)
}

func renderInterface(target BindingTarget, info introspect.InterfaceInfo, writer io.Writer, ifc_name, exportName string, infos *Infos) {
	//TODO: removing dependent on variable of target
	filterKeyWord(target, &info)

	log.Printf("Generate %q code for service:%q interface:%q ObjectName:%q", target, infos.DestName(), ifc_name, exportName)
	//TODO: move the common functions to the file of template_common.go
	funcs := template.FuncMap{
		"Lower":          lower,
		"Upper":          upper,
		"BusType":        func() string { return infos.BusType() },
		"PkgName":        func() string { return infos.PackageName() },
		"OBJ_NAME":       func() string { return "obj" },
		"TypeFor":        dbus.TypeFor,
		"getQType":       getQType,
		"DestName":       func() string { return infos.DestName() },
		"IfcName":        func() string { return ifc_name },
		"ExportName":     func() string { return exportName },
		"NormalizeQDBus": normalizeQDBus,
		"Normalize":      normalizeMethodName,
		"Ifc2Obj":        ifc2obj,
		"PropWritable":   func(prop introspect.PropertyInfo) bool { return prop.Access == "readwrite" },
		"GetOuts": func(args []introspect.ArgInfo) []introspect.ArgInfo {
			ret := make([]introspect.ArgInfo, 0)
			for _, a := range args {
				if a.Direction == "in" {
					ret = append(ret, a)
				}
			}
			return ret
		},
		"CalcArgNum": func(args []introspect.ArgInfo, direction string) (r int) {
			for _, arg := range args {
				if arg.Direction == direction {
					r++
				}
			}
			return
		},
		"Repeat": func(str string, sep string, times int) (r string) {
			for i := 0; i < times; i++ {
				if i != 0 {
					r += sep
				}
				r += str
			}
			return
		},
		"GetParamterNames": func(args []introspect.ArgInfo) (ret string) {
			for _, arg := range args {
				if arg.Direction == "in" {
					ret += ", "
					ret += arg.Name
				}
			}
			return
		},
		"GetParamterOuts": func(args []introspect.ArgInfo) (ret string) {
			var notFirst = false
			for _, arg := range args {
				if arg.Direction == "out" {
					if notFirst {
						ret += ","
					}
					notFirst = true
					ret += "&" + arg.Name
				}
			}
			return
		},
		"GetParamterOutsProto": func(args []introspect.ArgInfo) (ret string) {
			var notFirst = false
			for _, arg := range args {
				if arg.Direction == "out" {
					if notFirst {
						ret += ","
					}
					notFirst = true
					ret += arg.Name + " " + dbus.TypeFor(arg.Type)
				}
			}
			return
		},
		"GetParamterInsProto": func(args []introspect.ArgInfo) (ret string) {
			var notFirst = false
			for _, arg := range args {
				if arg.Direction == "in" {
					if notFirst {
						ret += ","
					}
					notFirst = true
					if strings.Contains(arg.Type, "(") {
						ret += arg.Name + " interface{}"
					} else {
						ret += arg.Name + " " + dbus.TypeFor(arg.Type)
					}
				}
			}
			return
		},
		"TryConvertObjectPath": func(prop introspect.PropertyInfo) string {
			if v := getObjectPathConvert("Property", prop.Annotations); v != "" {
				switch target {
				case GoLang:
					return tryConvertObjectPathGo(infos, prop.Type, v)
				case QML:
					return tryConvertObjectPathQML(prop.Type, v)
				}
			}
			return ""
		},
		"GetObjectPathType": func(prop introspect.PropertyInfo) (ret string) {
			if v := getObjectPathConvert("Property", prop.Annotations); v != "" {
				switch target {
				case GoLang:
					ret, _ = guessTypeGo(infos, prop.Type, v)
				case QML:
					ret, _ = guessTypeQML(prop.Type, v)
				}
				return
			}
			return dbus.TypeFor(prop.Type)
		},
	}
	templ := template.Must(template.New(exportName).Funcs(funcs).Parse(GetTemplate(target, TemplateTypeInterface)))
	templ.Execute(writer, info)
}

func renderTest(testPath, objName string, writer io.Writer, info introspect.InterfaceInfo, infos *Infos) {
	funcs := template.FuncMap{
		"TestPath": func() string { return testPath },
		"PkgName":  func() string { return infos.PackageName() },
		"ObjName":  func() string { return objName },
		/*"GetTestValue": func(args []dbus.ArgInfo) string {*/
		/*},*/
	}
	template.Must(template.New("testing").Funcs(funcs).Parse(__TEST_TEMPLATE)).Execute(writer, info)
}
