/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
// build +ignore
package main

type TemplateType int

const (
	TemplateTypeGlobal TemplateType = iota
	TemplateTypeInterface
	TemplateTypeInit
)

func GetTemplate(target BindingTarget, ttype TemplateType) string {
	var templs = map[string]string{
		"GLOBAL_pyqt":   __GLOBAL_TEMPLATE_PyQt,
		"GLOBAL_golang": __GLOBAL_TEMPLATE_GoLang,
		"GLOBAL_qml":    __GLOBAL_TEMPLATE_QML,

		"IFC_pyqt":   __IFC_TEMPLATE_PyQt,
		"IFC_golang": __IFC_TEMPLATE_GoLang,
		"IFC_qml":    __IFC_TEMPLATE_QML,

		"IFC_INIT_pyqt":   __IFC_TEMPLATE_INIT_PyQt,
		"IFC_INIT_golang": __IFC_TEMPLATE_INIT_GoLang,
		"IFC_INIT_qml":    __IFC_TEMPLATE_INIT_QML,
	}
	var name string
	switch target {
	case PyQt:
		name = "pyqt"
	case GoLang:
		name = "golang"
	case QML:
		name = "qml"
	default:
		panic("didn't support binding target" + target)
	}
	switch ttype {
	case TemplateTypeGlobal:
		return templs["GLOBAL_"+name]
	case TemplateTypeInterface:
		return templs["IFC_"+name]
	case TemplateTypeInit:
		return templs["IFC_INIT_"+name]
	default:
		panic("didn't support TemplateType")
	}
}
