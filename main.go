package main

import (
	"back-telega/lib"
	"back-telega/model"
	"back-telega/system"
	"back-telega/web"
	"log"
	"net/http"
)

func init() {
	lib.ClickDB = model.NewClick("tcp://192.168.0.145:9000")
	lib.RedisDB = model.NewRedis()
}

func main() {
	logs := make(chan lib.Json)
	comeGoodStatisticChan := make(chan lib.GoodJson)
	comeInitialSendGoodJSChan := make(chan lib.RequestGoodStatistic)

	handleRecivedBadStatistic := web.RecivedBadStatistic(logs)
	handleWorkerComeGoodStatistic := web.ComeGoodStatistic(comeGoodStatisticChan)
	handleInitialGoodPoint := web.InitialGoodPoint(comeInitialSendGoodJSChan)

	go system.WorkerSendGoodStatistic(comeGoodStatisticChan, comeInitialSendGoodJSChan)

	http.HandleFunc("/gateway/telegram/count-point", web.CountPoint)
	http.HandleFunc("/gateway/telegram/info-point", web.InfoPoint)
	http.HandleFunc("/gateway/telegram/list-point", web.ListAllPoint)
	http.HandleFunc("/gateway/telegram/list-point-today", web.ListPointToday)
	http.HandleFunc("/gateway/telegram/initial-good-point", handleInitialGoodPoint)
	http.HandleFunc("/gateway/telegram/create/bad", handleRecivedBadStatistic)
	http.HandleFunc("/gateway/telegram/create/good", handleWorkerComeGoodStatistic)

	err := http.ListenAndServe(":8181", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
