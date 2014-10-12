#What's this?
This tool is used to generate test xml files for pkg.linuxdeepin.com/go-dbus-generator.
It's mainly used by go-dbus-generator to run go test currently.

If you feel it's useful in other ways, e.g. to generate xml files for you dbus object,
welcome to make it more general.


#How use it?
1. go build && ./testdata (or other program if you are in different directory).
2. the default output directory is "output", you can change this by pass "output"
   parameter.
3. use go-dbus-generator to test the generated xml files


#TODO
1. write more test\_*special\_type*.go to cover more situation.
2. Fix all TODOs in test\_*special\_type*.go


#How can I add test data.
This is a simple tool, you can read and understand it quickly.
The simple way write test data is copy test\_tpl.go and adjust it.


#NOTE:
You can run "go test" in go-dbus-generator to verify whether you have find a new issue.
Don't forget run "go build && ./testdata" to generate new xml files after you have added
new test data.
