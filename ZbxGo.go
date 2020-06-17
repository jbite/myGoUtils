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
	Auth        string
}

func NewZabbix(url string, contentType string) *Zabbix {
	z := &Zabbix{}
	z.url = url
	z.contentType = contentType
	return z
}

func (z *Zabbix) Login(username string, password string) (result map[string]interface{}, err error) {
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
	result, err = z.Post()
	if err != nil {
		return nil, err
	}
	z.Auth = result["result"].(string)
	return result, nil
}

func (z *Zabbix) Post() (result map[string]interface{}, err error) {
	// fmt.Printf("%s\n", z.requestBody)
	resp, err := http.Post(z.url, z.contentType, bytes.NewBuffer(z.requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	// fmt.Println(result)
	return result, nil
}

func (z *Zabbix) Request(method string, params map[string]interface{}) (result map[string]interface{}, err error) {
	if z.Auth != "" {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  method,
			"params":  params,
			"auth":    z.Auth,
			"id":      2,
		})
		z.requestBody = requestBody
		result, err = z.Post()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	err = errors.New("not login")
	return nil, err
}
