package main

import (
	"fmt"
	"strings"
)

func normalizeSignature(s string) string {
	s = strings.Replace(s, "a{", "e_", -1)
	s = strings.Replace(s, "}", "_", -1)
	s = strings.Replace(s, "(", "r_", -1)
	s = strings.Replace(s, ")", "_", -1)
	return s
}

type QtTypeConvert struct {
}

func (QtTypeConvert) ToArray(element string) string {
	if element == "QString" {
		return "QStringList "
	} else if element == "uchar" {
		return "QByteArray"
	}
	return "QList<" + element + " >"
}

func (QtTypeConvert) ToDict(first string, second string) string {
	return fmt.Sprintf("QMap<%s, %s >", first, second)
}
func (QtTypeConvert) ToStruct(elements ...string) string {
	var r = "dbus::types::BaseStruct<"
	for i, element := range elements {
		if i != 0 {
			r = r + ", "
		}
		r = r + element
	}
	return r + " >"
}
func (QtTypeConvert) ToPrimitive(s string) string {
	switch s {
	case "y":
		return "uchar"
	case "b":
		return "bool"
	case "n":
		return "short"
	case "q":
		return "ushort"
	case "i":
		return "int"
	case "u":
		return "uint"
	case "x":
		return "qlonglong"
	case "t":
		return "qulonglong"
	case "d":
		return "double"
	case "h":
		return "int" //UNIX_FD
	case "s":
		return "QString"
	case "o":
		return "QDBusObjectPath"
	case "g":
		return "QDBusSignature"
	case "v":
		return "QDBusVariant"
	default:
		panic(fmt.Sprintf("%q is not a primitive signature", s))
	}
}
