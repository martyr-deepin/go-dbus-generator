package main

import "bytes"
import "text/template"
import "./introspect"

type CppBackend struct {
}

var globalTemplate = `
#ifndef __AUTO_GENERATED_DBUS__
#define __AUTO_GENERATED_DBUS__
#include <QtDBus>

namespace dbus {
	namespace common {
		{{ .Common }}
	}
	namespace types {
		{{ .Structs}}
	}
	namespace objects {
		{{ .Objects }}
	}
}

{{ .DeclareMetaTypes }}
#endif
`

func (b CppBackend) Stage3(infos []introspect.InterfaceInfo) string {
	var buffer = bytes.NewBuffer(nil)

	var info = struct {
		Common           string
		Structs          string
		Objects          string
		DeclareMetaTypes string
	}{}
	info.Structs, info.DeclareMetaTypes = b.Structs(infos)
	info.Common = s_common
	info.Objects = b.Objects(infos)

	template.Must(template.New("glboals").Funcs(template.FuncMap{}).Parse(globalTemplate)).Execute(buffer, info)
	return buffer.String()
}

var structsTemplate = `
	template<typename T1, typename T2=char, typename T3=char, typename T4=char, typename T5=char, typename T6=char, typename T7=char, typename T8=char, typename T9=char, typename T10=char, typename T11=char, typename T12=char, typename T13=char, typename T14=char>
	struct BaseStruct {
		T1 m1;
		T2 m2;
		T3 m3;
		T4 m4;
		T5 m5;
		T6 m6;
		T7 m7;
		T8 m8;
		T9 m9;
		T10 m10;
		T11 m11;
		T12 m12;
		T13 m13;
		T14 m14;
	};
	int getTypeId(const QString& s);

{{ range $i, $type := .Types }}
	typedef {{ nativeType $type}} {{ normalizeType $type}};
{{end}}

{{ range $i, $type := .Types }}
	inline QDBusArgument& operator<<(QDBusArgument &argument, const dbus::types::{{ normalizeType $type}}& v)
	{
		{{ marshall $type }}
	}
	inline const QDBusArgument& operator>>(const QDBusArgument &argument, dbus::types::{{ normalizeType $type}}& v)
	{
		{{ unmarshall $type }}
	} {{ end }}

	inline int getTypeId(const QString& s) {
	if (0) { 
	} {{ range $i, $type := .Types }} else if (s == "{{ $type.Signature }}") {
		return qDBusRegisterMetaType<dbus::types::{{ normalizeType $type}}>();
	}{{end}}
	}
`

func (CppBackend) Structs(infos []introspect.InterfaceInfo) (string, string) {
	var types []introspect.TypeMeta
	for _, sig := range introspect.MergeSignatures(infos...) {
		t := introspect.NewTypeMeta(sig)
		if t.Type == introspect.PrimitiveTypeId {
			continue
		}
		types = append(types, t)
	}

	var buffer = bytes.NewBuffer(nil)
	template.Must(template.New("types").Funcs(template.FuncMap{
		"marshall": func(s introspect.TypeMeta) string {
			switch s.Type {
			case introspect.StructTypeId:
				tt := introspect.NewStructTypeMeta(s.Signature)
				return marshallStruct("<<", len(tt.Elements))
			case introspect.ArrayTypeId:
				tt := introspect.NewArrayTypeMeta(s.Signature)
				return marshallArray(tt.Element.Signature)
			case introspect.DictTypeId:
				tt := introspect.NewDictTypeMeta(s.Signature)
				return marshallDict(tt.First.Signature, tt.Second.Signature)
			}

			return "--"
		},
		"unmarshall": func(s introspect.TypeMeta) string {
			switch s.Type {
			case introspect.StructTypeId:
				tt := introspect.NewStructTypeMeta(s.Signature)
				return marshallStruct(">>", len(tt.Elements))
			case introspect.ArrayTypeId:
				tt := introspect.NewArrayTypeMeta(s.Signature)
				return unmarshallArray(tt.Element.Signature)
			case introspect.DictTypeId:
				tt := introspect.NewDictTypeMeta(s.Signature)
				return unmarshallDict(tt.First.Signature, tt.Second.Signature)
			}
			return "--"
		},
		"nativeType": func(s introspect.TypeMeta) string {
			return tTable.Get(s.Signature)
		},
		"normalizeType": func(s introspect.TypeMeta) string {
			return normalizeSignature(s.Signature)
		},
	}).Parse(structsTemplate)).Execute(buffer, struct{ Types []introspect.TypeMeta }{types})

	var declareMetaTypes = ""
	for _, t := range types {
		declareMetaTypes = declareMetaTypes + "Q_DECLARE_METATYPE(dbus::types::" + normalizeSignature(t.Signature) + ");\n"
	}
	return buffer.String(), declareMetaTypes
}

//预处理nodeinfo信息
func (CppBackend) Stage1(ninfo *introspect.NodeInfo) *introspect.NodeInfo {
	return ninfo
}

//touch所有的文件
func (CppBackend) Stage2(ninfo *introspect.NodeInfo) []string {
	var r []string

	r = append(r, "common.h", "common.cpp")

	for _, ifc := range ninfo.Interfaces {
		r = append(r, ifc.Name+".h", ifc.Name+".cpp")
	}
	return r
}
