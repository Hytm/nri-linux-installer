package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func writeToFile(value, fileName string) (err error) {
	s := []byte(value)
	err = ioutil.WriteFile(tmpFileName, s, 0644)
	if err != nil {
		fmt.Printf("Write File Error: %s\n", err)
		return
	}
	err = runCmd(fmt.Sprintf("sudo mv %s %s", tmpFileName, fileName))
	if err != nil {
		fmt.Printf("Move File Error: %s\n", err)
	}
	return
}

func runCmd(cmdString string) (err error) {
	cmd := exec.Command("/bin/sh", "-c", cmdString)
	err = cmd.Run()
	return
}

func parseOSRelease(in string) (d, v string) {
	s := strings.Split(in, "\n")
	for i := 0; i < len(s); i++ {
		kv := strings.Split(s[i], "=")
		if kv[0] == "VERSION_ID" {
			v = strings.ReplaceAll(kv[1], "\"", "")
		}
		if kv[0] == "ID" {
			d = strings.ReplaceAll(kv[1], "\"", "")
		}
	}
	return
}
