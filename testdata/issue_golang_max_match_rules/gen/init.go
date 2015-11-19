package testdbus

import "pkg.deepin.io/lib/dbus"
import "fmt"
import "sync"

var __conn *dbus.Conn = nil
var __connLock sync.Mutex

var __objCounter map[string]int = nil
var __objCounterLock sync.Mutex

func getBus() *dbus.Conn {
	__connLock.Lock()
	defer __connLock.Unlock()
	if __conn == nil {
		var err error
		__conn, err = dbus.SessionBus()
		if err != nil {
			panic(err)
		}
	}
	return __conn
}

func getObjCounter() map[string]int {
	__objCounterLock.Lock()
	defer __objCounterLock.Unlock()
	if __objCounter == nil {
		__objCounter = make(map[string]int)
	}
	return __objCounter
}

func dbusCall(method string, flags dbus.Flags, args ...interface{}) (err error) {
	err = getBus().BusObject().Call(method, flags, args...).Err
	if err != nil {
		fmt.Println(err)
	}
	return
}

func dbusAddMatch(rule string) (err error) {
	return dbusCall("org.freedesktop.DBus.AddMatch", 0, rule)
}

func dbusRemoveMatch(rule string) (err error) {
	return dbusCall("org.freedesktop.DBus.RemoveMatch", 0, rule)
}

func incObjCount(objName string) {
	objCounter := getObjCounter()

	__objCounterLock.Lock()
	defer __objCounterLock.Unlock()
	objCounter[objName]++
}

func decObjCount(objName string) (cleanRules bool) {
	objCounter := getObjCounter()

	__objCounterLock.Lock()
	defer __objCounterLock.Unlock()
	if _, ok := objCounter[objName]; !ok {
		return false
	}
	objCounter[objName]--
	if objCounter[objName] == 0 {
		delete(objCounter, objName)
		return true
	}
	return false
}
