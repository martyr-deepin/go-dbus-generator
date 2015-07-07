dbus-generate is a static dbus binding tool.

Currently, it support *QDBus* target and golang target with *pkg.deepin.io/lib/dbus*.

#Note
The pyqt target is not maintained, because I haven't use this target in real world.
The golang and qml target will be maintained long-term and it has been verified by
deepin OS team in many project.

#TODO
1. Improve generated code quality.
   Especially for qml target, the QDBusType hasn't handle very well currently.
2. Write more test code.
3. Make the test.qml support user input.
4. Refactor code for easier adding new target and make code more readable.

#Feture plan
C language with Gtk?
