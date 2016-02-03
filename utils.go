/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import "strings"
import "pkg.deepin.io/lib/dbus/introspect"

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
func ifc2obj(ifc string) string { return "/" + strings.Replace(ifc, ".", "/", -1) }

func isValidInterface(s string) bool {
	if len(s) == 0 || len(s) > 255 || s[0] == '.' {
		return false
	}
	elem := strings.Split(s, ".")
	if len(elem) < 2 {
		return false
	}
	for _, v := range elem {
		if len(v) == 0 {
			return false
		}
		if v[0] >= '0' && v[0] <= '9' {
			return false
		}
		for _, c := range v {
			if !isMemberChar(c) {
				return false
			}
		}
	}
	return true
}
func isValidMember(s string) bool {
	if len(s) == 0 || len(s) > 255 {
		return false
	}
	i := strings.Index(s, ".")
	if i != -1 {
		return false
	}
	if s[0] >= '0' && s[0] <= '9' {
		return false
	}
	for _, c := range s {
		if !isMemberChar(c) {
			return false
		}
	}
	return true
}

func isMemberChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') || c == '_'
}

func getMember(target BindingTarget, str string) string {
	fields := strings.Split(str, ".")
	if len(fields) < 1 {
		return ""
	}
	r := fields[len(fields)-1]
	if isValidMember(r) {
		if target == GoLang {
			r = strings.ToLower(r)
		}
		return r
	} else {
		return ""
	}
}
func getQMLPkgName(str string) (r string) {
	fields := strings.Split(str, ".")
	for i, field := range fields {
		if i != 0 {
			r += "."
		}
		r += upper(field)
	}
	return
}

// com.deepin.DBus.ObjectPathConvert.Property for properties
// com.deepin.DBus.ObjectPathConvert.Out1..  for methods and signals
// com.deepin.DBus.ObjectPathConvert.Arg1.. for methods input
func getObjectPathConvert(suffix string, annotations []introspect.AnnotationInfo) string {
	for _, v := range annotations {
		if v.Name == "com.deepin.DBus.ObjectPathConvert."+suffix {
			return v.Value
		}
	}
	return ""
}
