package main

import (
	"fmt"
)

func marshallStruct(direction string, size int) string {
	var s = "argument.beginStructure();\n"

	s = s + "argument"
	for i := 0; i < size; i++ {
		s = s + " " + direction + " "
		s = s + fmt.Sprintf("v.m%d", i+1)
	}
	s = s + ";\n"

	s = s + "argument.endStructure();\n"
	s = s + "return argument;\n"
	return s
}
func marshallArray(eleSig string) string {
	var s = fmt.Sprintf("argument.beginArray(getTypeId(\"%s\"));\n", eleSig)

	s = s + "for (int i=0; i < v.size(); ++i)\n"
	s = s + "    argument << v.at(i);\n"

	s = s + "argument.endArray();\n"
	s = s + "return argument;\n"
	return s
}
func unmarshallArray(eleSig string) string {
	var s = "argument.beginArray();\n"

	s = s + "while (!argument.atEnd()) {\n"
	s = s + fmt.Sprintf("    %s ele;\n", tTable.Get(eleSig))
	s = s + "    argument >> ele;\n"
	s = s + "    v.append(ele);\n"
	s = s + "}\n"

	s = s + "argument.endArray();\n"
	s = s + "return argument;\n"
	return s
}

func marshallDict(firstSig string, secondSig string) string {
	var s = fmt.Sprintf("argument.beginMap(getTypeId(\"%s\"), getTypeId(\"%s\"));\n", firstSig, secondSig)

	s = s + fmt.Sprintf("QList<%s > keys;\n", tTable.Get(firstSig))
	s = s + "for (int i=0; i < keys.size(); ++i) {\n"
	s = s + "    argument.beginMapEntry();\n"
	s = s + "    argument << keys[i] << v[keys[i]];\n"
	s = s + "    argument.endMapEntry();\n"
	s = s + "}\n"

	s = s + "argument.endMap();\n"
	s = s + "return argument;\n"
	return s
}
func unmarshallDict(firstSig string, secondSig string) string {
	var s = "argument.beginMap();\n"
	s = s + "v.clear();\n"

	s = s + "while (!argument.atEnd()) {\n"
	s = s + fmt.Sprintf("    %s key;\n", tTable.Get(firstSig))
	s = s + fmt.Sprintf("    %s value;\n", tTable.Get(secondSig))
	s = s + "    argument.beginMapEntry();\n"
	s = s + "    argument >> key >> value;\n"
	s = s + "    argument.endMapEntry();\n"
	s = s + "    v.insert(key, value);\n"
	s = s + "}\n"

	s = s + "argument.endMap();\n"
	s = s + "return argument;\n"
	return s
}
