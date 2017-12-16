package web

import (
	"back-telega/model"
	"back-telega/var"
	"fmt"
	"net/http"
	"log"
	"encoding/json"
)

func CountPoint(w http.ResponseWriter, r *http.Request) {
	var pointCount global.PointCount
	err := model.PingClick(global.ClickDB)
	if err != nil {
		fmt.Fprintf(w, "ошибка подключения %v", err)
	}
	rows, err := global.ClickDB.Query("select count(distinct point_id) from stat.statistics where created = today()")
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