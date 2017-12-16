package main

import (
	"back-telega/model"
	"back-telega/var"
	"back-telega/web"
	"log"
	"net/http"
)

func init() {
	lib.ClickDB = model.NewClick("tcp://192.168.0.145:9000")
	lib.RedisDB = model.NewRedis()
}

func main() {
	http.HandleFunc("/gateway/telegram/count-point", web.CountPoint)
	http.HandleFunc("/gateway/telegram/info-point", web.InfoPoint)
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
