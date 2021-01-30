package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func findConfig(config string) bool {
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return false
	}
	return true
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the field
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func writeText(file, text string) {
	config, _ := os.Create(file)
	defer config.Close()
	config.WriteString(text)
}

func readText(configFilePath string) string {
	file, err := os.Open(configFilePath)
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	delay := scanner.Text()
	//Convert the string read from config to int
	return string(delay)
}

func sendNotif(title, message, icon string, time int) {
	fmt.Println(strconv.Itoa(time))
	if icon != "" {
		exec.Command(CMD, title, message, "-i", icon, "-t", strconv.Itoa(time)).Run()
	} else {
		exec.Command(CMD, title, message, "-t", string(time)).Run()
	}

}

////////////////////////////Curl Library////////////////////////////
type NotiBody struct {
	Title string `json:"title"`
	Msg   string `json:"msg"`
	Icon  string `json:"icon"`
	Time  int    `json:"time"`
}

func curled(w http.ResponseWriter, r *http.Request) {
	//{"title":"iii","msg":"Hello","icon":"gnome-foot","time":19000}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// error handling
		w.Write([]byte("{\"status\":false}"))
		return
	}
	var req NotiBody
	errparse := json.Unmarshal(body, &req)
	if errparse != nil {
		// error handling
		w.Write([]byte("{\"status\":false}"))
		return
	}
	fmt.Println(req.Icon)
	if req.Title == "" {
		w.Write([]byte("{\"status\":false}"))
		return
	}
	sendNotif(req.Title, req.Msg, req.Icon, req.Time)
	w.Write([]byte("{\"status\":true}"))
	fmt.Println(string(body))
}
