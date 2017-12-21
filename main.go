package main

import (
	"back-telega/model"
	"back-telega/var"
	"back-telega/web"
	"fmt"
	"log"
	"net/http"
	"time"
	//	"encoding/json"
	"bytes"
	//"io/ioutil"
	"io/ioutil"
	"encoding/json"
)

func init() {
	lib.ClickDB = model.NewClick("tcp://192.168.0.145:9000")
	lib.RedisDB = model.NewRedis()
}

func main() {
	keyChannel := make(chan int)
	logs := make(chan lib.Json)
	//	keysChannel := make(chan []int)
	listen := web.ListenAddKey(keyChannel)
	print := web.MakeHello(logs)
	ticker := time.NewTicker(time.Second * 1)
	go listenParser(ticker, keyChannel, logs)
	http.HandleFunc("/gateway/telegram/count-point", web.CountPoint)
	http.HandleFunc("/gateway/telegram/info-point", web.InfoPoint)
	http.HandleFunc("/gateway/telegram/list-point", web.ListAllPoint)
	http.HandleFunc("/gateway/telegram/listen-add-key", listen)
	http.HandleFunc("/gateway/telegram/list-point-today", web.ListPointToday)
	http.HandleFunc("/gateway/telegram/print", print)
	err := http.ListenAndServe(":8181", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func listenParser(ticker *time.Ticker, keyChannel chan int, logs chan lib.Json) {
	delete := time.NewTicker(time.Second * 30)
	//	var keys []int
	keyDelete := make(chan int)
	oldLenKeys := -1
	for {
		select {
		case k := <-keyChannel:
			lib.Keys = append(lib.Keys, k)
			go func(delete *time.Ticker, key chan int) {
				select {
				case <-delete.C:
					key <- k
				}
			}(delete, keyDelete)
		case <-ticker.C:
			if len(lib.Keys) != oldLenKeys {
				fmt.Println(lib.Keys)
				oldLenKeys = len(lib.Keys)
			}
		case k := <-keyDelete:
			lib.Keys = deleteFromValue(k, lib.Keys)
		case js := <-logs:
			fmt.Println("Отправляю на телеграм: ", js)
			SendAllStatistic(js)
		}
	}
}

func deleteFromValue(value int, array []int) []int {
	fmt.Println("Удаляю", value)
	var arrayCopy []int
	for _, a := range array {
		if a != value {
			arrayCopy = append(arrayCopy, a)
		}
	}
	fmt.Println("Удалил")
	return arrayCopy
}

func SendAllStatistic(jsonRaw lib.Json) {
	url := "http://127.0.0.1:8282/listen"
	jsonStr,_ := json.Marshal(jsonRaw)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка отправки статистики на телеграм:", err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
