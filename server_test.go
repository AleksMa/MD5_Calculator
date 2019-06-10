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
	body, _ := ioutil.ReadAll(resp.Body)
	idStruct := struct { Id string `json:"id"`}{}
	json.Unmarshal(body, &idStruct)
	id := idStruct.Id

	// Посылаем check на отправленную ранее задачу
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
	body, _ := ioutil.ReadAll(resp.Body)
	idStruct := struct { Id string `json:"id"`}{}
	json.Unmarshal(body, &idStruct)
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
	body, _ = ioutil.ReadAll(resp2.Body)
	taskStruct := Task{}

	err = json.Unmarshal(body, &taskStruct)
	if err != nil {
		t.Errorf("JSON parse failed")
		return
	}
	if taskStruct.Url != url || taskStruct.Status != "done" {
		t.Errorf("Done test failed")
	}
}

func TestNotExist(t *testing.T) {
	// Посылаем check с заведомо невозможным id
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

	// Посылаем check на отправленную ранее задачу
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
	if taskStruct.Status != "error" {
		t.Errorf("Error test failed")
	}
}