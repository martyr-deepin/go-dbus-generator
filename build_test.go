package main

import "fmt"
import "testing"
import "os/exec"

func TestGolang(t *testing.T) {
	output := "tmp_golang__"
	infos := loadInfos("dbus.in.json")
	infos.normalize(output, "golang")

	geneateInit(infos)
	generateMain(infos)
	formatCode(infos)

	if out, err := exec.Command("bash", "-c", fmt.Sprintf("cd %s && ls && go build", output)).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		t.Fatal("Build:" + err.Error())
	}

	if err := exec.Command("rm", "-rf", "tmp_golang__").Run(); err != nil {
		t.Fatal("rm:" + err.Error())
	}
}

func TestQML(t *testing.T) {
	output := "tmp_qml__"
	infos := loadInfos("dbus.in.json")
	infos.normalize(output, "qml")

	geneateInit(infos)
	generateMain(infos)
	formatCode(infos)

	if out, err := exec.Command("bash", "-c", fmt.Sprintf("cd %s && ls && qmake && make", output)).CombinedOutput(); err != nil {
		fmt.Println(string(out))
		t.Fatal("Build:" + err.Error())
	}

	if err := exec.Command("rm", "-rf", output).Run(); err != nil {
		t.Fatal("rm:" + err.Error())
	}
}
