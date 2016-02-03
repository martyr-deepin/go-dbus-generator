/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package main

import "testing"
import C "launchpad.net/gocheck"
import testdbus "./gen"
import "fmt"
import "time"
import "math/rand"

type testWrap struct{}

func Test(t *testing.T) { C.TestingT(t) }
func init() {
	C.Suite(&testWrap{})
}

func (*testWrap) TestMaxMatchRules(c *C.C) {
	fmt.Println("TestMaxMatchRules: should not print errors in this case")
	createLargeNumObjects()
	createAndDestroyLargeNumObjects()
	createAndDestroyLargeNumObjectsRandom()
	createAndDestroyLargeNumObjectsLazy()
}
func createLargeNumObjects() {
	for i := 0; i < 1000; i++ {
		testdbus.NewTestDBus(dbusDest, dbusPath)
	}
}
func createAndDestroyLargeNumObjects() {
	for i := 0; i < 1000; i++ {
		obj, _ := testdbus.NewTestDBus(dbusDest, dbusPath)
		testdbus.DestroyTestDBus(obj)
	}
}
func createAndDestroyLargeNumObjectsRandom() {
	for i := 0; i < 1000; i++ {
		obj, _ := testdbus.NewTestDBus(dbusDest, dbusPath)
		for j := 0; j <= rand.Intn(5); j++ {
			testdbus.DestroyTestDBus(obj)
		}
	}
}
func createAndDestroyLargeNumObjectsLazy() {
	objs := make([]*testdbus.TestDBus, 0)
	for i := 0; i < 1000; i++ {
		obj, _ := testdbus.NewTestDBus(dbusDest, dbusPath)
		objs = append(objs, obj)
	}
	for _, obj := range objs {
		testdbus.DestroyTestDBus(obj)
	}
}

func (*testWrap) TestDBusSignals(c *C.C) {
	fmt.Println("TestDBusSignals: catch dbus signals...")
	dbusobj, err := testdbus.NewTestDBus(dbusDest, dbusPath)
	if err != nil {
		c.Error(err)
	}

	var sigChan = make(chan bool)
	var propChan = make(chan bool)
	dbusobj.ConnectTimerSignal(func(msg string) {
		fmt.Println("catch dbus signal:", msg)
		sigChan <- true
	})
	dbusobj.TimerProp.ConnectChanged(func() {
		fmt.Println("catch dbus prop:", dbusobj.TimerProp.Get())
		propChan <- true
	})

	var destroyChann = make(chan bool)
	go func() {
		createAndDestroyLargeNumObjects()
		destroyChann <- true
	}()
	createAndDestroyLargeNumObjects()
	createAndDestroyLargeNumObjectsRandom()
	createAndDestroyLargeNumObjectsLazy()
	<-destroyChann

	// start dbus service
	go SetupDBus()

	// wait for result
	select {
	case <-time.After(2 * time.Second):
		c.Error("catch dbus siganl failed")
	case <-sigChan:
		// catch dbus siganl success
	}
	select {
	case <-time.After(2 * time.Second):
		c.Error("catch dbus properties changed failed")
	case <-propChan:
		// catch dbus property changed success
	}
}
