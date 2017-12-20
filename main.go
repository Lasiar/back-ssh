package main

import (
	"back-telega/model"
	"back-telega/var"
	"back-telega/web"
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
	lib.ClickDB = model.NewClick("tcp://192.168.0.145:9000")
	lib.RedisDB = model.NewRedis()
}

func main() {
	keyChannel := make(chan int)
	listen := web.Listen(keyChannel)
	ticker := time.NewTicker(time.Second * 1)
	go listenParser(ticker, keyChannel)
	http.HandleFunc("/gateway/telegram/count-point", web.CountPoint)
	http.HandleFunc("/gateway/telegram/info-point", web.InfoPoint)
	http.HandleFunc("/gateway/telegram/list-point", web.ListAllPoint)
	http.HandleFunc("/gateway/telegram/listen", listen)
	http.HandleFunc("/gateway/telegram/list-point-today", web.ListPointToday)
	err := http.ListenAndServe(":8181", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func listenParser(ticker *time.Ticker, keyChannel chan int) {
	delete := time.NewTicker(time.Second * 30)
	var keys []int
	keyDelete := make(chan int)
	oldLenKeys := -1
	for {
		select {
		case k := <-keyChannel:
			keys = append(keys, k)
			go func(delete *time.Ticker, key chan int) {
				select {
				case <-delete.C:
					key <- k
				}
			}(delete, keyDelete)
		case <-ticker.C:
			if len(keys) != oldLenKeys {
				fmt.Println(keys)
				oldLenKeys = len(keys)
			}
		case k := <-keyDelete:
			keys = deleteFromValue(k, keys)

		}
	}
}

func deleteFromValue(value int, array []int) ([]int) {
	fmt.Println("Удаляю", value)
	var arrayCopy []int
	for _, a := range array{
		if a != value {
			arrayCopy = append(arrayCopy, a)
		}
	}
	fmt.Println("Удалил")
	return arrayCopy
}
