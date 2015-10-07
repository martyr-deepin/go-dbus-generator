package main

import "./introspect"
import "bytes"
import "fmt"
import "strings"
import "text/template"

var tplObject = `
{{ beginNamespace .}}
{{ typedef .}}
class {{ baseName .Name }} : public dbus::common::DBusObject
{
	Q_OBJECT
	private:
	static const char *defaultService() { return "{{ serviceName . }}";}
	static const QDBusObjectPath defaultPath() { return QDBusObjectPath("{{ objectPath . }}");}
	public:
        {{ baseName .Name }}(QString addr="session", QObject* parent=0)
        :DBusObject(parent, defaultService(), defaultPath().path(), "{{.Name}}", addr)
        {
        }
	{{ baseName .Name }}(QString addr, QString service, QString path, QObject* parent=0)
	:DBusObject(parent, service, path, "{{.Name}}", addr)
	{
	}
	~{{ baseName .Name }}(){}

	{{range $i, $property := .Properties}}
	Q_PROPERTY(dbus::common::R<{{goodType .Type}} > {{ upper .Name }} READ {{ lower .Name }} NOTIFY {{ lower .Name}}Changed)
	dbus::common::R<{{goodType .Type}} > {{lower .Name}} () {
		QDBusPendingReply<> call = fetchProperty("{{.Name}}");
		return dbus::common::R<{{goodType .Type}} >(call, dbus::common::PropertyConverter);
	}
	{{end}}


	{{range $i, $method := .Methods}}
	{{ $outDeclare := outTypes $method }}
	{{ if eq $outDeclare "" }}
	dbus::common::R<void> {{.Name}} ({{ inTypes .}}) {
		QList<QVariant> argumentList;
		{{ inTypesDo .}}
		QDBusPendingReply<> call = asyncCallWithArgumentList(QLatin1String("{{.Name}}"), argumentList);
                return dbus::common::R<void>(call);
	}
	{{else}}
	dbus::common::R<{{ $outDeclare }}> {{.Name}} ({{ inTypes .}}) {
		QList<QVariant> argumentList;
		{{ inTypesDo .}}
		QDBusPendingReply<> call = asyncCallWithArgumentList(QLatin1String("{{.Name}}"), argumentList);
		return dbus::common::R<{{ $outDeclare }}>(call);
	}
	{{end}}

	{{end}}

	Q_SIGNALS:
	{{range $i, $signal := .Signals}}
	void {{.Name}}({{ sigTypes .}}); 
	{{end}}

	{{range $i, $property := .Properties}}
	void {{ lower .Name}}Changed (); {{end}}

};
{{ endNamespace .}}
`

func (CppBackend) Object(iinfo introspect.InterfaceInfo) string {

	var buffer = bytes.NewBuffer(nil)
	e := template.Must(template.New("objects").Funcs(template.FuncMap{
		"serviceName": introspect.QueryServiceName,
		"objectPath":  introspect.QueryObjectPath,
		"lower":       lower,
		"upper":       upper,
		"goodType":    goodType,

		"beginNamespace": func(iinfo introspect.InterfaceInfo) string {
			var r = ""
			ss := strings.Split(iinfo.Name, ".")
			for _, f := range ss[0 : len(ss)-1] {
				r = r + "namespace " + lower(f) + " {"
			}
			return r
		},
		"endNamespace": func(iinfo introspect.InterfaceInfo) string {
			return strings.Repeat("}", len(strings.Split(iinfo.Name, "."))-1)
		},
		"baseName": func(name string) string {
			ss := strings.Split(name, ".")
			return ss[len(ss)-1]
		},
		"typedef": func(iinfo introspect.InterfaceInfo) string {
			return ""
		},
		"outTypes": func(m introspect.MethodInfo) string {
			var sigs []string
			for _, arg := range m.Args {
				if arg.Direction == "out" {
					sigs = append(sigs, arg.Type)
				}
			}
			return declareType(sigs, "")
		},
		"sigTypes": func(m introspect.SignalInfo) string {
			return declareArgs(m.Args)
		},
		"inTypes": func(m introspect.MethodInfo) string {
			var sigs []string
			for _, arg := range m.Args {
				if arg.Direction != "out" {
					sigs = append(sigs, arg.Type)
				}
			}
			return declareType(sigs, "arg")
		},
		"inTypesDo": func(m introspect.MethodInfo) string {
			var r = ""
			var sigs []string
			for _, arg := range m.Args {
				if arg.Direction != "out" {
					sigs = append(sigs, arg.Type)
				}
			}
			if len(sigs) != 0 {
				r = "argumentList"
			}
			for i := 0; i < len(sigs); i++ {
				if i != len(sigs) {
					r = r + " << "
				}
				r = r + fmt.Sprintf("QVariant::fromValue(arg%d)", i)
			}
			return r + ";"
		},
	}).Parse(tplObject)).Execute(buffer, iinfo)
	e = e
	return buffer.String()
}
func (c CppBackend) Objects(infos []introspect.InterfaceInfo) string {
	var r = ""
	for _, iinfo := range infos {
		r = r + c.Object(iinfo)
	}
	return r
}

func lower(str string) string {
	if str == "" {
		return ""
	}

	return strings.ToLower(str[:1]) + str[1:]
}
func upper(str string) string {
	if str == "" {
		return ""
	}

	return strings.ToUpper(str[:1]) + str[1:]
}

func declareArgs(args []introspect.ArgInfo) string {
	var r = ""
	for i, arg := range args {
		if i != 0 {
			r += ", "
		}
		r = r + goodType(arg.Type)
		name := arg.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		r = r + fmt.Sprintf(" %s", name)
	}
	return r
}
