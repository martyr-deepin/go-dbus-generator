# DBus Generator

**Description**:
DBus Generator is static dbus binding tool for multipble languages.

Currently, it support *QDBus* target and golang target with *pkg.deepin.io/lib/dbus*.

## Dependencies

### Build Dependencies

- go >= 1.1
- pkg.linuxdeepin.com/lib/dbus (This will be removed when the branch feature/qt completed)
- pkg.linuxdeepin.com/lib/dbus/introspect (This will be removed when the branch feature/qt completed)

## Installation

```
go get pkg.linuxdeepin.com/lib/dbus
go get pkg.linuxdeepin.com/lib/dbus/introspect
go get pkg.deepin.io/dbus-generator
```

## Usage

```
Usage of dbus-generator:
  -in="dbus.in.json": the config file path
  -out="": output directory
  -target="": the target language to generate binding for [QML|PyQt|GoLang]
```

the dbus.in.json format looks like

```
{
	"Config": {
		"OutputDir": ".", # which directory to save the generated code.
		"InputDir": "../xml", # where to find the dbus introspect xml file.
		"BusType": "System",  # Session Bus or System Bus
		"DestName": "com.deepin.api.Device" # the DBus Service name
	},
	"Interfaces": [
		{
			"Interface": "com.deepin.api.Device",
			"OutFile": "device",
			"XMLFile": "com.deepin.api.Device.xml",
			"ObjectName": "Device"
		}
	]
}
```

[DBus Factory](https://github.com/linuxdeepin/dbus-factory) is the real example to show how to use dbus generator.

## Known issues

The pyqt target is no longger maintained by deepin team  due to lack of application using it.
The golang and qml target will be maintained and it has been verified by
deepin team in many projects.


## TODO

- Remove the build dependencies of pkg.deepin.io/lib/dbus
- Write more test code.
- Make the test.qml support user input.
- Refactor code for easier adding new target and make code more readable.

## Getting help

Any usage issues can ask for help via

* [Gitter](https://gitter.im/orgs/linuxdeepin/rooms)
* [IRC channel](https://webchat.freenode.net/?channels=deepin)
* [Forum](https://bbs.deepin.org)
* [WiKi](http://wiki.deepin.org/)

## Getting involved

We encourage you to report issues and contribute changes

* [Contirubtion guide for
users](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Users)
* [Contribution guide for developers](http://wiki.deepin.org/index.php?title=Contribution_Guidelines_for_Developers).

## License

DBus Geneator is licensed under [GPLv3](LICENSE).
