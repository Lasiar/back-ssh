package configure

import (
	"back-telega/lib"
	"io/ioutil"
	"fmt"
	"log"
	"encoding/json"
)

func ReadConfig() lib.Configure {
	var config lib.Configure
	var tmp struct {
		A []json.Number `json:"chat-id"`
	}
	file, err := ioutil.ReadFile("config/config")
	if err != nil {
		log.Panic("Can not read configuration file", err)
	}

	err = json.Unmarshal(file, &tmp)
	if err != nil {
		fmt.Println("Unmarshal config", err)
	}
	config.ChatId = make([]int64, len(tmp.A))
	for i := range tmp.A{
		v, err := tmp.A[i].Int64()
		if err != nil {
			log.Println(err)
		}
		config.ChatId[i] = v
	}
	return config
}
