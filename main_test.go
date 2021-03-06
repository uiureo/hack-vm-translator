package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const emulatorPath = "/Users/zat/Downloads/nand2tetris/tools/CPUEmulator.sh"

func testFile(t *testing.T, input string, bootstrap bool) {
	name := strings.Split(filepath.Base(input), ".")[0]
	cmd := exec.Command("go", "run", "main.go", input)
	if !bootstrap {
		cmd.Env = append(os.Environ(), "BOOTSTRAP=false")
	}
	asmCode, err := cmd.Output()

	if err != nil {
		t.Error(err)
	}

	dir := filepath.Dir(input) + "/"
	ioutil.WriteFile(dir+name+".asm", asmCode, os.ModePerm)

	out, err := exec.Command(emulatorPath, dir+name+".tst").Output()
	if err != nil {
		t.Error(name, string(out))
	}
}

func TestMain(t *testing.T) {
	noInitFiles, err := filepath.Glob("fixtures/no-init/*.vm")
	if err != nil {
		t.Error(err)
	}

	for _, input := range noInitFiles {
		testFile(t, input, false)
	}

	files, err := filepath.Glob("fixtures/*.vm")
	if err != nil {
		t.Error(err)
	}

	for _, input := range files {
		testFile(t, input, true)
	}
}

func TestStaticsTest(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "fixtures/StaticsTest/Class1.vm", "fixtures/StaticsTest/Class2.vm", "fixtures/StaticsTest/Sys.vm")

	asmCode, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}

	dir := "fixtures/StaticsTest/"
	ioutil.WriteFile(dir+"StaticsTest.asm", asmCode, os.ModePerm)

	out, err := exec.Command(emulatorPath, dir+"StaticsTest.tst").Output()
	if err != nil {
		t.Error(string(out))
	}
}
