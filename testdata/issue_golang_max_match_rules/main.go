package main

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"time"
)

const (
	dbusDest      = "com.deepin.TestDBus"
	dbusPath      = "/com/deepin/TestDBus"
	dbusInterface = "com.deepin.TestDBus"
)

type testDBus struct {
	// property
	TimerProp string

	// signal
	TimerSignal func(string)
}

func (self *testDBus) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusInterface,
	}
}

func (self *testDBus) startTimerNotify() {
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)

			self.TimerProp = "Hello Prop"
			dbus.NotifyChange(self, "TimerProp")

			dbus.Emit(self, "TimerSignal", "Hello Signal")
		}
	}()
}

func SetupDBus() {
	obj := &testDBus{}

	err := dbus.InstallOnSession(obj)
	if err != nil {
		fmt.Println(err)
		return
	}

	obj.startTimerNotify()

	if err = dbus.Wait(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	SetupDBus()
}
