/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type BindingTarget string

const (
	QML    BindingTarget = "qml"
	GoLang BindingTarget = "golang"
	PyQt   BindingTarget = "pyqt"
)

type _Config struct {
	Target       BindingTarget
	NotExportBus bool
	OutputDir    string
	InputDir     string
	PkgName      string
	DestName     string
	BusType      string
}

type Infos struct {
	Interfaces []_Interface `json:"Interfaces"`
	Config     _Config      `json:"Config"`
	outputs    map[string]io.Writer
}

func LoadInfos(path string, outputDir string, target string) (*Infos, error) {
	infos := NewInfos()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&infos)
	if err != nil {
		return nil, fmt.Errorf(`failed parse "%s": %s`, path, err.Error())
	}
	err = infos.normalize(outputDir, target)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

const DefaultOutputDir = "out"

func (infos *Infos) normalize(outputDir string, target string) error {
	var output string
	if outputDir != "" {
		output = outputDir
	} else if infos.OutputDir() != "" {
		output = infos.OutputDir()
	} else {
		output = DefaultOutputDir
	}
	err := infos.SetOutputDir(output)
	if err != nil {
		return err
	}

	if target != "" {
		err := infos.SetTarget((target))
		if err != nil {
			return err
		}
	}

	switch infos.Target() {
	case GoLang:
		err := infos.SetBusType(upper(infos.BusType()))
		if err != nil {
			return err
		}
	case PyQt, QML:
		err := infos.SetBusType(lower(infos.BusType()))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Please set the correct target before invoke SetBusType: %s", infos.Target())
	}

	if infos.PackageName() == "" {
		guestName := getMember(infos.Target(), infos.DestName())
		if guestName == "" {
			return fmt.Errorf("What's the name of pacakge you want to generated? Please specify an PkgName or DestName field.")
		}
		err := infos.SetPackageName(guestName)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (i *Infos) BusType() string {
	return i.Config.BusType
}
func (i *Infos) PackageName() string {
	return i.Config.PkgName
}
func (i *Infos) SetPackageName(name string) error {
	i.Config.PkgName = name
	return nil
}
func (i *Infos) ListInterfaces() []_Interface {
	return i.Interfaces
}
func (i *Infos) DestName() string {
	return i.Config.DestName
}
func (i *Infos) OutputDir() string {
	return i.Config.OutputDir
}
func (i *Infos) InputDir() string {
	return i.Config.InputDir
}
func (i *Infos) SetOutputDir(dir string) error {
	if dir != "" {
		err := os.MkdirAll(dir, os.ModePerm|os.ModeDir)
		if err != nil {
			return err
		}
	}
	i.Config.OutputDir = dir
	return nil
}
func (i *Infos) Target() BindingTarget {
	return i.Config.Target
}
func (i *Infos) SetTarget(target string) error {
	t := BindingTarget(strings.ToLower(target))
	switch t {
	case GoLang, QML:
	case PyQt:
		log.Println("Warning: pyqt support is not completed.")
	default:
		return fmt.Errorf("Didn't supported target type: %s", target)
	}
	i.Config.Target = t
	return nil
}
func (i *Infos) SetBusType(bus string) error {
	_bus := strings.ToLower(bus)
	if _bus != "session" && _bus != "system" {
		return fmt.Errorf("Didn't support bus type %s", bus)
	}
	i.Config.BusType = bus

	return nil
}
func (i *Infos) GetWriter(name string) (io.Writer, bool) {
	f, ok := i.outputs[name]
	if ok {
		return f, false
	}

	var suffix string
	if path.Ext(name) == "" {
		switch i.Target() {
		case GoLang:
			suffix = ".go"
		case PyQt:
			suffix = ".py"
		case QML:
			suffix = ".h"
		}
	}

	f, err := os.Create(path.Join(i.OutputDir(), name+suffix))
	if err != nil {
		panic(err)
	}
	i.outputs[name] = f
	return f, true
}

func NewInfos() *Infos {
	return &Infos{
		outputs: make(map[string]io.Writer),
	}
}
