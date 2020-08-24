package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Test_Path(t *testing.T) {
	execpath, err := os.Executable() // 获得程序路径
	// handle err ...
	if err != nil {
		fmt.Printf("%v \n", err)
	}
	configfile := filepath.Join(filepath.Dir(execpath), ".conf/microservice.yml")
	log.Println(configfile)


	configfile = strings.Replace(configfile, "\\", "/", -1)
	log.Println(configfile)
}

func GetCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	s = strings.Replace(s, "\\", "/", -1)
	s = strings.Replace(s, "\\\\", "/", -1)
	i := strings.LastIndex(s, "/")
	path := string(s[0 : i+1])
	return path
}