package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
}

type Opensea_json struct {
	Stats struct {
		FloorPrice float64 `json:"floor_price"`
	} `json:"stats"`
}

func floor_price(token string) float64 {

	// 发送 HTTP 请求
	url := fmt.Sprintf("https://api.opensea.io/api/v1/collection/%s/stats?format=json", token)
	resp, err := http.Get(url)
	checkError(err)
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	var opensea_json Opensea_json
	err = json.Unmarshal(body, &opensea_json)
	checkError(err)
	return opensea_json.Stats.FloorPrice
}
