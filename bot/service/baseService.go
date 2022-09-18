package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type BaseService struct {
}

// AtALL @所有人服务
func (b BaseService) AtALL(msg string, groupID int64) {
	b.SendGroupMessage("[CQ:at,qq=all]"+msg, groupID, false)
}

// SendGroupMessage 发送群消息服务
func (b BaseService) SendGroupMessage(msg string, groupID int64, randomFlag bool) (messageID int64) {
	if randomFlag == true {
		t := rand.Int63n(1000 * 3)
		log.Printf("time: %d\n", t)
		time.Sleep(time.Millisecond * time.Duration(t))
	}
	url := BaseURL + "/send_group_msg"
	data, err := json.Marshal(map[string]interface{}{
		"group_id": groupID,
		"message":  msg,
	})
	if err != nil {
		log.Println("序列化数据失败! " + err.Error())
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Println("发送消息失败! " + err.Error())
	}

	resDataByte, _ := ioutil.ReadAll(response.Body)
	var resData map[string]interface{}
	json.Unmarshal(resDataByte, &resData)
	status := resData["status"].(string)
	if status != "ok" {
		log.Println(resData)
		return 0
	}
	log.Println(resData)
	return int64(resData["data"].(map[string]interface{})["message_id"].(float64))
}

// DeleteMessage 撤回消息
func (b BaseService) DeleteMessage(msgID int64) {
	url := BaseURL + "/delete_msg"
	data, err := json.Marshal(map[string]interface{}{
		"message_id": msgID,
	})
	if err != nil {
		log.Println("序列化数据失败! " + err.Error())
	}
	_, err = http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Println("撤销消息失败! " + err.Error())
	}
}
