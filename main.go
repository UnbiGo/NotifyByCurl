package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var CMD = "notify-send"
var ADDR = ":8000"

func main() {
	//Search in config path if there is the directory water-reminder
	OS := runtime.GOOS
	var configPath string
	home, _ := os.LookupEnv("HOME")
	if OS == "darwin" {
		configPath = filepath.Join(home, "Library/Application Support")
	} else {
		configPath = filepath.Join(home, ".config")
	}
	configDirPath := filepath.Join(configPath, "curl_notify")
	configCMD := filepath.Join(configDirPath, "command")
	configAddr := filepath.Join(configDirPath, "addr")
	if !findConfig(configDirPath) {
		//Create config directory
		os.Mkdir(configDirPath, 0700)
	}
	if !findConfig(configCMD) {
		//Create config directory
		writeText(configCMD, CMD)
	} else {
		CMD = readText(configCMD)
	}
	if !findConfig(configAddr) {
		//Create config directory
		writeText(configAddr, ADDR)
	} else {
		ADDR = readText(configAddr)
	}
	ADDR = strings.ReplaceAll(ADDR, "\n", "")
	CMD = strings.ReplaceAll(CMD, "\n", "")

	fmt.Println("Running at: " + ADDR)
	http.HandleFunc("/notify", curled)
	http.ListenAndServe(ADDR, nil)
}
