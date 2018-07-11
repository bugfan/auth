package com

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bugfan/logrus"
)

func HasValues(data map[string]interface{}, keys ...string) error {
	for _, v := range keys {
		str := HasValue(data, v)
		if str != "" {
			return errors.New(str)
		}
	}
	return nil
}
func HasValue(data map[string]interface{}, key string) string {
	if _, exist := data[key]; exist {
		return ""
	}
	return fmt.Sprintf("缺少参数[%s]", key)
}
func GetCtxMap(r *http.Request) map[string]interface{} {
	return JsonToMap([]byte(GetCtxJson(r)))
}
func GetCtxJson(r *http.Request) string {
	requestbody, _ := ioutil.ReadAll(r.Body)
	return string(requestbody)
}
func JsonToMap(val []byte) map[string]interface{} {
	data := make(map[string]interface{})
	json.Unmarshal(val, &data)
	logrus.Println("[IN] ", data)
	return data
}
func ToJsonString(v interface{}) string {
	return string(ToJsonByte(v))
}
func ToJsonByte(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

type Result struct {
	Status interface{} `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Msg    string      `json:"msg,omitempty"`
}

type List struct {
	Offset interface{} `json:"offset"`
	Total  int64       `json:"total"`
	List   interface{} `json:"list"`
}
