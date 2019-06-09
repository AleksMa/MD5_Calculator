package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Task struct {
	url    string
	hash   [16]byte
	status int
}

const (
	notExist = http.StatusNotFound
	running  = http.StatusAccepted
	done     = http.StatusOK
	error    = http.StatusInternalServerError
)

var (
	tasks = make(map[string]Task)
)

func stringStatus(status int) string {
	if status == running {
		return "running"
	} else if status == done {
		return "done"
	}
	return "error"
}

func SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================POST==========================")
	r.ParseForm()
	params := r.Form
	fmt.Println(params)
	url := params["url"][0]

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body2, _ := ioutil.ReadAll(resp.Body)
	outputStream := string(body2)

	hash := md5.Sum([]byte(outputStream))
	fmt.Println(hash)

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	id := string(out)[:len(string(out))-1]
	fmt.Printf("%s\n", id)

	w.WriteHeader(http.StatusOK)
	w.Write(out)

	tasks[id] = Task{url, hash, done}
	fmt.Println("================END=OF=POST=======================")
}

func CheckRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================GET==========================")
	r.ParseForm()
	params := r.Form
	id := params["id"][0]
	fmt.Println(id)
	task, ok := tasks[id]
	fmt.Println(task)
	if ok {
		fmt.Println("NOT ERROR")
		w.WriteHeader(task.status)
		fmt.Fprintln(w, stringStatus(task.status))
	} else {
		fmt.Println("ERROR")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "not exist")
	}
	fmt.Println("================END=OF=GET=======================")
}
func main() {

	http.HandleFunc("/submit", SubmitResponseHandler)
	http.HandleFunc("/check", CheckRouterHandler)

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}

}
