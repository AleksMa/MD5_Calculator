package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

const (
	notExist      = "not exist"
	running       = "running"
	done          = "done"
	internalError = "error"
)

type Task struct {
	Hash   string `json:"md5"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

var (
	tasks = make(map[string]*Task)
)

func intStatus(status string) int {
	if status == running {
		return http.StatusAccepted
	} else if status == done {
		return http.StatusOK
	} else if status == notExist {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func hashError(id string) {
	tasks[id].Status = internalError
}

func makeHash(url string, id string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Cannot download resource")
		hashError(id)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read body")
		hashError(id)
		return
	}
	log.Println(string(body)[:100])

	hash := md5.Sum(body)
	log.Printf("Hash: %v\n", hex.EncodeToString(hash[:16]))

	task := tasks[id]
	task.Status = done
	task.Hash = hex.EncodeToString(hash[:16])
}

func SubmitRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================POST==========================")

	// Получаем URL
	r.ParseForm()
	url := r.FormValue("url")
	fmt.Println(url)

	// Генерим ID
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	id := string(out)[:len(string(out))-1]
	fmt.Printf("%s\n", id)

	// Сохраняем задачу (running), отправляем ответ
	tasks[id] = &Task{Url: url, Status: running}

	w.WriteHeader(http.StatusAccepted)
	idStruct := struct {
		Id string `json:"id"`
	}{id}
	answer, err := json.Marshal(idStruct)
	w.Write(answer)
	fmt.Fprintln(w)

	// Скачиваем файл и считаем хэш (в фоне)
	go makeHash(url, id)

	fmt.Println("================END=OF=POST=======================")
}

func CheckRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====================GET==========================")
	r.ParseForm()
	id := r.FormValue("id")
	fmt.Println(id)

	//  Ищем id в задачах
	task, ok := tasks[id]
	fmt.Println(task)
	var answer []byte
	if ok {
		w.WriteHeader(intStatus(task.Status))

		if task.Status == done {
			answer, _ = json.Marshal(tasks[id])
		} else {
			statusStruct := struct {
				Status string `json:"status"`
			}{task.Status}
			answer, _ = json.Marshal(statusStruct)
		}
	} else {
		statusStruct := struct {
			Status string `json:"status"`
		}{"not exist"}
		answer, _ = json.Marshal(statusStruct)
	}
	w.Write(answer)
	fmt.Fprintln(w)
	fmt.Println("================END=OF=GET=======================")
}
func main() {

	http.HandleFunc("/submit", SubmitRouterHandler)
	http.HandleFunc("/check", CheckRouterHandler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
