# MD5 Calculator  

## Установка и запуск  
С установленным и настроенным Go.  
В терминале:  
* ```git clone https://github.com/AleksMa/MD5_Calculator.git```  
* ```cd MD5_Calculator```  
* ```go run server.go```    
---

Поскольку сторонние библиотеки не были использованы, 
установка и сборка с помощью *glide.sh* не требуется.

Сервер слушает на порту 8000.  
Пожалуйста, проверьте, свободен ли этот порт.

## Тесты  
На каждый возможный статус ответа сервера подготовлены go-тесты ( *server_test.go* )  
Для их запуска необходимо в директории проекта ( *somepath/MD5_Calculator* ) запустить в терминале
``` go test ```.   
Перед этим, разумеется, необходимо запустить сам сервер.

## API
Веб-сервис, считающий MD5-хеш некоторого ресурса по его URL.
* **POST**-запрос на **/submit** с параметром **url**. Создает задачу с идентификатором в формате UUID, 
в качестве ответа отправляет данный идентификатор. В фоновом режиме происходит загрузка и хеширование ресурса с переданным **url**.   
Пример работы:  
*Запрос:* ```curl -X POST -d "url=https://golang.org" http://localhost:8000/submit```  
*Ответ:* ```{"id":"ac1cf466-eef6-47d6-9c0a-1adefe3f36a9"}```  
Можно также отправлять и **GET**-запрос:  
*Запрос:* ```curl -X GET http://localhost:8000/submit?url=https://golang.org```  
*Ответ:* ```{"id":"1b97926c-8c87-427d-b46d-8ef8419bc174"}```  
И даже **HEAD**:  
*Запрос:* ```curl -I http://localhost:8000/submit?url=https://golang.org```  
*Ответ:* ```HTTP/1.1 202 Accepted...```    
* **GET**-запрос на **/check** с параметром **id**. 
Возвращает статус задачи с данным **id** и соответствующий код ответа:
 {"not exist", 404; "running", 202; "done", 200; "error", 500}. 
 Если статус задачи "done", в ответе возвращает также хеш и url.
 Кроме того, запросы могут быть отправлены и как **POST**, и как **HEAD**.  
 **GET:**  
*Запрос*: ```curl -X GET http://localhost:8000/check?id=ac1cf466-eef6-47d6-9c0a-1adefe3f36a9```  
*Ответ:* ```{"md5":"8acf67d3ff985922cc057a52b3348e83","status":"done","url":"https://golang.org"}```   
**HEAD:**   
*Запрос*: ```curl -I http://localhost:8000/check?id=NotAnUUID```  
*Ответ:* ```HTTP/1.1 404 Not Found...```   

<p align="center">
  <img src="https://cdn-images-1.medium.com/max/1200/1*yh90bW8jL4f8pOTZTvbzqw.png" alt="PES"/>
</p>
