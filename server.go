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

// Map, хранящий все полученные задачи ( id -> {hash, status, url} )
var (
	tasks = make(map[string]*Task)
)

// Отображение текстовых статусов задач в коды состояния HTTP
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

// Изменение статуса задачи на "упавшую с ошибкой в процессе выполнения"
func hashError(id string) {
	tasks[id].Status = internalError
}

// Горутина, скачивающая файл по id, считающая его хеш и обнавляющая состояние задачи по id
func makeHash(url string, id string) {
	log.Println("==Downloading and encoding resource==")
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Request error. URL: %v; ID: %v. Cannot download resource", url, id)
		hashError(id)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Request error. URL: %v; ID: %v. Cannot read resource", url, id)
		hashError(id)
		return
	}

	hash := md5.Sum(body)
	log.Printf("Hash: %x", hash)

	task := tasks[id]
	task.Status = done
	task.Hash = hex.EncodeToString(hash[:16])
}

func SubmitRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("==SUBMIT request==")
	// Получаем URL
	r.ParseForm()
	url := r.FormValue("url")	// Если параметра url нет или url задан некорректно, задача примет статус "error"
	log.Printf("URL: %v", url)

	// Генерим UUID
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	id := string(out)[:len(string(out))-1]	// Стандартная функция записывает в конец UUID символ '\n'
	log.Printf("UUID: %v", id)

	// Сохраняем задачу (статус "running"), отправляем ответ
	tasks[id] = &Task{Url: url, Status: running}

	w.WriteHeader(http.StatusAccepted)

	// Создаем временную структуру, содержащую статус, для записи в json
	idStruct := struct {
		Id string `json:"id"`
	}{id}

	answer, err := json.Marshal(idStruct)
	w.Write(answer)
	fmt.Fprintln(w)

	// Скачиваем файл и считаем хэш (в горутине)
	go makeHash(url, id)
}

func CheckRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("==CHECK request==")
	r.ParseForm()
	// При отстуствующем/некорректном параметре позднее получим "not exist"
	id := r.FormValue("id")
	log.Printf("ID: %v", id)

	//  Ищем id в задачах
	task, ok := tasks[id]
	var answer []byte
	statusStruct := struct { Status string `json:"status"` }{}
	if ok {
		w.WriteHeader(intStatus(task.Status))

		if task.Status == done {
			answer, _ = json.Marshal(tasks[id])
		} else {
			statusStruct.Status = task.Status
			answer, _ = json.Marshal(statusStruct)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		statusStruct.Status = "not exist"
		answer, _ = json.Marshal(statusStruct)
	}
	w.Write(answer)
	fmt.Fprintln(w)
}
func main() {

	http.HandleFunc("/submit", SubmitRouterHandler)
	http.HandleFunc("/check", CheckRouterHandler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}
