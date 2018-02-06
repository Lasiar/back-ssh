package system

import (
	"back-telega/lib"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
)

func WorkerSendGoodStatistic(jsonStatistic chan lib.GoodJson, jsonRequest chan lib.RequestGoodStatistic) {
	var arrPoint []int
	var arrJsOld []lib.GoodJson
	for {
		select {
		case js := <-jsonRequest:
			arrPoint = append(arrPoint, js.Point...)
		case js := <-jsonStatistic:
			fmt.Println("Пришло сообщение")
			if len(arrJsOld) != 0 {
				for i, jsOld := range arrJsOld {
					switch {
					case jsOld.Point != js.Point:
						break
					case jsOld.Md5 == js.Md5:
						fmt.Println("WARNING")
					default:
						arrJsOld = reWrite(js, arrJsOld, i)
					}
				}
			} else {
				arrJsOld = append(arrJsOld, js)
			}
			for _, v := range arrPoint{
				if js.Point == v {
					sendGoodStatistic(js)
				}
			}
		}
	}
}

func reWrite(goodJson lib.GoodJson, arrGoodJson []lib.GoodJson, element int) []lib.GoodJson{
	var newArray []lib.GoodJson
	for i, js := range arrGoodJson{
		if element == i {
			newArray = append(newArray, goodJson)
		} else {
			newArray = append(newArray, js)
		}
	}
	return newArray
}


func sendGoodStatistic(js lib.GoodJson) {
	url := "http://127.0.0.1:8282/good"
	jsonStr, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}
	req.Header.Set("X-Custom-Header", "json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}