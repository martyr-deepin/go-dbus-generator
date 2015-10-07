package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "path"

import "./introspect"

func BuildNodeInfo(fpath string) *introspect.NodeInfo {
	reader, err := os.Open(fpath)
	if err != nil {
		panic(err.Error())
	}
	ninfo, err := introspect.Parse(reader)
	if err != nil {
		fmt.Println(err.Error() + ":" + fpath)
		return nil
	}
	return ninfo
}

func TouchFiles(base string, list []string) {
	for _, fpath := range list {
		fullPath := path.Join(base, fpath)
		err := os.MkdirAll(path.Dir(fullPath), os.ModePerm|os.ModeDir)
		if err != nil {
			fmt.Println("E:", err)
		}
		f, err := os.Create(fullPath)
		if err != nil {
			fmt.Println("W:", err)
			continue
		}
		f.Close()
	}
}

func MergeNodes(ninfos ...introspect.NodeInfo) []introspect.InterfaceInfo {
	if ninfos == nil {
		return nil
	}
	cache := make(map[string]introspect.InterfaceInfo)
	for _, ninfo := range ninfos {
		for _, iinfo := range ninfo.Interfaces {
			cache[iinfo.Name] = iinfo
		}
		for _, iinfo := range MergeNodes(ninfo.Children...) {
			cache[iinfo.Name] = iinfo
		}
	}
	var r []introspect.InterfaceInfo
	for _, v := range cache {
		r = append(r, v)
	}
	return r
}

var outFile = flag.String("o", "test/dbus.h", "the file to store generated code")

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Please specify the xml files")
		return
	}

	var ii []introspect.NodeInfo
	for _, xml := range flag.Args() {
		i := BuildNodeInfo(xml)
		if i == nil {
			fmt.Println("E: Build at: ", xml)
			continue
		}
		fmt.Println("Build:", xml)
		ii = append(ii, *i)
	}
	infos := MergeNodes(ii...)
	cpp := CppBackend{}
	fmt.Println("Try generating codes to ", *outFile)
	ioutil.WriteFile(*outFile, ([]byte)(cpp.Stage3(infos)), os.ModePerm)
}
