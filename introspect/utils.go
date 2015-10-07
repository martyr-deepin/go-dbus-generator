package introspect

import (
	"sort"
	"strings"
)

func QueryServiceName(info InterfaceInfo) string {
	v := info.GetAnnotation(AnnotationConfigurationDefaultServiceName)
	if v == "" {
		return info.Name
	}
	return v
}
func QueryObjectPath(info InterfaceInfo) string {
	v := info.GetAnnotation(AnnotationConfigurationDefaultObjectPath)
	if v == "" {
		v = "/" + strings.Replace(info.Name, ".", "/", -1)
	}

	return v
}

type TypeSet []string

func (s TypeSet) Len() int {
	return len(s)
}
func (s TypeSet) Less(i, j int) bool {
	l1, l2 := len(s[i]), len(s[j])
	return l1 < l2
}
func (s TypeSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TypeSet) Append(i string) TypeSet {
	for _, t := range s {
		if t == i {
			return s
		}
	}
	return append(s, i)
}

func MergeSignatures(infos ...InterfaceInfo) []string {
	var r TypeSet
	//分析XML中出现过的signatures
	for _, info := range infos {
		for _, m := range info.Methods {
			for _, a := range m.Args {
				r = r.Append(a.Type)
			}
		}
		for _, p := range info.Properties {
			r = r.Append(p.Type)
		}

		for _, s := range info.Signals {
			for _, a := range s.Args {
				r = r.Append(a.Type)
			}
		}

	}

	//分解这些signatures得到隐含的合法子类型
	tt := NewTypesTable(DummyConvert{})
	for _, sig := range r {
		tt.Get(sig)
	}

	for _, v := range tt.GetCompoundTypeSignatures() {
		r = r.Append(v)
	}

	sort.Sort(r)
	return r
}
