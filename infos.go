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
	GoLang               = "golang"
	PyQt                 = "pyqt"
)

type _Config struct {
	Target       string
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

func loadInfos(path string) *Infos {
	infos := NewInfos()

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&infos)
	if err != nil {
		panic(err)
	}
	return infos
}
func (infos *Infos) normalize(outputDir string, target string) {
	if outputDir != "out" {
		infos.SetOutputDir(outputDir)
	} else if len(infos.OutputDir()) == 0 {
		infos.SetOutputDir(outputDir)
	}

	if target != "" {
		infos.SetTarget(target)
	}

	os.MkdirAll(infos.OutputDir(), 0755)

	infos.SetBusType(infos.BusType())
	infos.SetTarget(infos.Target())

	if infos.PackageName() == "" {
		name := getMember(BindingTarget(infos.Target()), infos.DestName())
		infos.SetPackageName(name)
		if infos.PackageName() == "" {
			log.Fatal("Didn't specify an PkgName and can't calclus an valid PkgName by DestName:" + infos.DestName())
		}
	}
}

func (i *Infos) BusType() string {
	return i.Config.BusType
}
func (i *Infos) PackageName() string {
	return i.Config.PkgName
}
func (i *Infos) SetPackageName(name string) {
	i.Config.PkgName = name
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
func (i *Infos) SetOutputDir(dir string) {
	i.Config.OutputDir = dir
}
func (i *Infos) Target() string {
	return i.Config.Target
}
func (i *Infos) SetTarget(target string) {
	i.Config.Target = target
}
func (i *Infos) SetBusType(bus string) error {
	bus = strings.ToLower(bus)
	if bus != "session" && bus != "system" {
		return fmt.Errorf("Didn't support bus type %s", bus)
	}

	switch BindingTarget(i.Config.Target) {
	case GoLang:
		i.Config.BusType = upper(bus)
	case PyQt:
		i.Config.BusType = lower(bus)
	case QML:
		i.Config.BusType = lower(bus)
	}
	return nil
}
func (i *Infos) GetWriter(name string) (io.Writer, bool) {
	f, ok := i.outputs[name]
	if ok {
		return f, false
	}

	var suffix string
	if path.Ext(name) == "" {
		switch BindingTarget(i.Target()) {
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
	i := Infos{
		outputs: make(map[string]io.Writer),
	}
	return &i
}
