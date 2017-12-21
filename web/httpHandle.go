package web

import (
	"back-telega/model"
	"back-telega/var"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func ListPointToday(w http.ResponseWriter, r *http.Request) {
	var listPointToday lib.PointAllList
	var onePoint int
	err := model.PingClick(lib.ClickDB)
	if err != nil {
		fmt.Fprintf(w, "ошибка подключения %v", err)
	}
	rows, err := lib.ClickDB.Query("select distinct point_id from stat.statistics where created = today()")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		if err := rows.Scan(&onePoint); err != nil {
			log.Fatal(err)
		}
		listPointToday.Point = append(listPointToday.Point, onePoint)
	}
	jsonBlob, _ := json.Marshal(listPointToday)
	fmt.Fprint(w, string(jsonBlob))
}

func CountPoint(w http.ResponseWriter, _ *http.Request) {
	var pointCount lib.PointCount
	err := model.PingClick(lib.ClickDB)
	if err != nil {
		fmt.Fprintf(w, "ошибка подключения %v", err)
	}
	rows, err := lib.ClickDB.Query("select count(distinct point_id) from stat.statistics where created = today()")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		if err := rows.Scan(&pointCount.Count); err != nil {
			log.Fatal(err)
		}
	}
	jsonBlob, _ := json.Marshal(pointCount)
	fmt.Fprint(w, string(jsonBlob))
}

func InfoPoint(w http.ResponseWriter, r *http.Request) {
	var infoPoint lib.InfoPointJs
	match, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	point := r.FormValue("point")
	keys, err := lib.RedisDB.Keys(point + "_[iu][ps]*").Result()
	if err != nil {
		fmt.Fprintf(w, "err get keys redis: %v", err)
	}
	vals, err := lib.RedisDB.MGet(keys...).Result()
	if err != nil {
		infoPoint.Success = false
		jsonBlob, _ := json.Marshal(infoPoint)
		fmt.Fprint(w, string(jsonBlob))
		return
	}
	for _, val := range vals {
		if match.MatchString(val.(string)) {
			infoPoint.Ip = val.(string)
		} else {
			infoPoint.UserAgent = val.(string)

		}
		infoPoint.Success = true
	}
	jsonBlob, _ := json.Marshal(infoPoint)
	fmt.Fprint(w, string(jsonBlob))
}

func ListAllPoint(w http.ResponseWriter, r *http.Request) {
	var pointList lib.PointAllList
	keysIp, err := lib.RedisDB.Keys("[0-9]*_ip").Result()
	if err != nil {
		fmt.Fprintf(w, "err get keys; %v", err)
	}
	for _, keyIp := range keysIp {
		point, _ := strconv.Atoi(strings.Trim(keyIp, "_ip"))
		pointList.Point = append(pointList.Point, point)
	}
	jsonBlob, _ := json.Marshal(pointList)
	fmt.Fprint(w, string(jsonBlob))
}

func ListenAddKey(keyChan chan int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		fmt.Fprint(w, "all ok	")
		keyInt, _ := strconv.Atoi(key)
		keyChan <- keyInt
	}
}

func MakeHello(logger chan lib.Json) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprint(w, key)
		fmt.Fprint(w, "all ok")
		decoder := json.NewDecoder(r.Body)
		var t lib.Json
		err := decoder.Decode(&t)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("keys: ", lib.Keys)
		fmt.Println("id=", t.Point)
		for _, k := range lib.Keys {
			fmt.Println(t.Point, "=", k)
			if t.Point == k {
				fmt.Println("Отправил", t)
				logger <- t
				fmt.Println("Отправил логгер")
			}
		}
	}
}
