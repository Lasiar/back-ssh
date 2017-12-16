package web

import (
	"back-telega/model"
	"back-telega/var"
	"fmt"
	"net/http"
	"log"
	"encoding/json"

	"strings"
	"strconv"
)

func CountPoint(w http.ResponseWriter, r *http.Request) {
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
	jsonBlob, _:= json.Marshal(pointCount)
	fmt.Fprint(w, string(jsonBlob))
}

func InfoPoint(w http.ResponseWriter, r *http.Request) {
	var  infoPoint lib.InfoPointJs
	point :=  r.FormValue("point")
	keys, err := lib.RedisDB.Keys(point+"_[iu][ps]*").Result()
	if err != nil {
		fmt.Fprintf(w, "err get keys redis: %v", err)
	}
	vals, err := lib.RedisDB.MGet(keys...).Result()
	if err != nil {
		infoPoint.Success = false
		jsonBlob, _ :=  json.Marshal(infoPoint)
		fmt.Fprint(w, string(jsonBlob))
		return
	}
	for i, val := range vals {
		switch i {
		case 0:
			infoPoint.Ip = val.(string)
		case 1:
			infoPoint.UserAgent = val.(string)
		}
		infoPoint.Success = true
	}
	jsonBlob, _ :=  json.Marshal(infoPoint)
	fmt.Fprint(w, string(jsonBlob))
}

func ListPoint(w http.ResponseWriter, r *http.Request) {
	var pointList lib.PointList
	keysIp , err := lib.RedisDB.Keys("[0-9]*_ip").Result()
	if err != nil {
		fmt.Fprintf(w, "err get keys; %v", err)
	}
	for _, keyIp := range keysIp {
		point, _ := strconv.Atoi(strings.Trim(keyIp, "_ip"))
		pointList.Point = append(pointList.Point, point)
		}
	jsonBlob, _:= json.Marshal(pointList)
	fmt.Fprint(w, string(jsonBlob))
}