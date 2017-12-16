package main

import (
	"back-telega/model"
	"back-telega/var"
	"back-telega/web"
	"log"
	"net/http"
)

func init() {
	global.ClickDB = model.NewClick("tcp://192.168.0.145:9000")
}

func main() {
	http.HandleFunc("/gateway/telegram/count-point", web.CountPoint)
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
