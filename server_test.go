package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestRunning(t *testing.T) {
	// Посылаем задачу (submit, для простоты используем GET-запрос)
	url := "https://github.com/AleksMa"
	resp, err := http.Get("http://localhost:8000/submit?url=" + url)
	if err != nil {
		t.Errorf("Submit request failed")
		return
	}
	defer resp.Body.Close()

	// Парсим ответ сервера, выделяем id задачи
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Reading answer failed")
	}

	idStruct := struct { Id string `json:"id"`}{}
	err = json.Unmarshal(body, &idStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
	}
	id := idStruct.Id

	// Посылаем check на отправленную ранее задачу
	resp2, err := http.Get("http://localhost:8000/check?id=" + id)
	if err != nil {
		t.Errorf("Check request failed")
		return
	}
	defer resp2.Body.Close()

	// Читаем ответ
	body, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Errorf("Reading answer failed")
	}
	taskStruct := Task{}

	err = json.Unmarshal(body, &taskStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
	}
	// Так как запрос статуса задачи происходит практически мновенно, задача не должна быть обработана
	if taskStruct.Status != "running" {
		t.Errorf("Running test failed")
	}
}

func TestDone(t *testing.T) {
	// Посылаем задачу (submit, для простоты используем GET-запрос)
	url := "https://github.com/AleksMa"
	resp, err := http.Get("http://localhost:8000/submit?url=" + url)
	if err != nil {
		t.Errorf("Submit request failed")
		return
	}
	defer resp.Body.Close()

	// Парсим ответ сервера, выделяем id задачи
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Reading answer failed")
	}

	idStruct := struct { Id string `json:"id"`}{}
	err = json.Unmarshal(body, &idStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
	}
	id := idStruct.Id

	// Ждем, чтобы задача обработалась
	time.Sleep(time.Second * 2)

	// Посылаем check на отправленную ранее задачу
	resp2, err := http.Get("http://localhost:8000/check?id=" + id)
	if err != nil {
		t.Errorf("Check request failed")
		return
	}
	defer resp2.Body.Close()

	// Читаем ответ
	body, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		t.Errorf("Reading answer failed")
	}
	taskStruct := Task{}

	err = json.Unmarshal(body, &taskStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
	}

	// За 2 секунды задача должна быть обработана (причем верно)
	if taskStruct.Url != url || taskStruct.Status != "done" {
		t.Errorf("Done test failed")
	}
}

func TestNotExist(t *testing.T) {
	// Посылаем check с заведомо невозможным id (не являющимся UUID)
	resp2, err := http.Get("http://localhost:8000/check?id=SomeIdThatCantBeInBase")
	if err != nil {
		t.Errorf("Check request failed")
		return
	}
	defer resp2.Body.Close()

	// Читаем ответ
	body, _ := ioutil.ReadAll(resp2.Body)
	taskStruct := Task{}

	err = json.Unmarshal(body, &taskStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
		return
	}
	// Сервер должен вернуть Not Found
	if taskStruct.Status != "not exist" {
		t.Errorf("Not exist test failed")
	}
}

func TestError(t *testing.T) {
	// Посылаем задачу с заведомо неправильным URL
	url := "SomethingButNotURL"
	resp, err := http.Get("http://localhost:8000/submit?url=" + url)
	if err != nil {
		t.Errorf("Submit request failed")
		return
	}
	defer resp.Body.Close()

	// Парсим ответ сервера, выделяем id задачи
	body, _ := ioutil.ReadAll(resp.Body)
	idStruct := struct { Id string `json:"id"`}{}
	json.Unmarshal(body, &idStruct)
	id := idStruct.Id

	// Ждем, чтобы задача обработалась
	time.Sleep(time.Second * 2)

	// Посылаем check на отправленную ранее задачу (невалидную)
	resp2, err := http.Get("http://localhost:8000/check?id=" + id)
	if err != nil {
		t.Errorf("Check request failed")
		return
	}
	defer resp2.Body.Close()

	// Читаем ответ
	body, _ = ioutil.ReadAll(resp2.Body)
	taskStruct := Task{}

	err = json.Unmarshal(body, &taskStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
		return
	}

	// Обработка задачи должна была закончиться ошибкой
	if taskStruct.Status != "error" {
		t.Errorf("Error test failed")
	}
}