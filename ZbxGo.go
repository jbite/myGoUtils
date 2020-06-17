package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Zabbix struct {
	url         string
	contentType string
	requestBody []byte
}

func NewZabbix(url string, contentType string) *Zabbix {
	z := &Zabbix{}
	z.url = url
	z.contentType = contentType
	return z
}

func (z *Zabbix) login(username string, password string) (result map[string]interface{}, err error) {
	z.requestBody, err = json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "user.login",
		"params": map[string]interface{}{
			"user":     username,
			"password": password,
		},
		"id":   1,
		"auth": nil,
	})

	if err != nil {
		return nil, err
	}
	resp, err := http.Post(z.url, z.contentType, bytes.NewBuffer(z.requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	return result, nil
}
func main() {
	url := "http://10.39.0.108/zabbix/api_jsonrpc.php"
	contentType := "application/json-rpc"
	z := NewZabbix(url, contentType)
	r, err := z.login("Admin", "password")

	if err != nil {
		fmt.Println("err", err)
	}

	log.Println(r)
}
