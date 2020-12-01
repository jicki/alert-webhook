package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	MsgType string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
		Mentioned_list string `json:"mentioned_list"`
		Mentioned_mobile_list string `json:"mentioned_mobile_list"`
	} `json:"text"`

}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:annotations`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}

type Notification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:receiver`
	GroupLabels       map[string]string `json:groupLabels`
	CommonLabels      map[string]string `json:commonLabels`
	CommonAnnotations map[string]string `json:commonAnnotations`
	ExternalURL       string            `json:externalURL`
	Alerts            []Alert           `json:alerts`
}

var defaultRobot = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=a7967dbc-f04c-40f6-a888-02dac71ffe1f"
var defaultRobot1 = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=063d487f-058f-4a7f-9e8f-16f5f6d6196e"

// 告警接收人
func SendMessage(notification Notification,defaultRobot string) {
	var msgres = make(map[string]string)
	msgres["mentioned_mobile_list"] = "18270117154"
	//var buffer bytes.Buffer
	//buffer.WriteString(notification.GroupKey)
	msg,err := json.Marshal(notification.GroupLabels)
	if err != nil {
		log.Println("notification.GroupLabels Marshal failed,", err)
		return
	}
	msg1,err := json.Marshal(notification.CommonAnnotations["summary"])
	if err != nil {
		log.Println("notification.CommonAnnotations Marshal failed,", err)
		return
	}
	//msg2,err := json.Marshal(notification.CommonAnnotations["description"])
	//if err != nil {
	//	log.Println("notification.CommonAnnotations Marshal failed,", err)
	//	return
	//}
	// 告警消息
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("告警: %v\n",string(msg)))
	buffer.WriteString(fmt.Sprintf("Endpoint: %v\n",string(msg1)))
	//buffer.WriteString(fmt.Sprintf("告警描述: %v\n",string(msg2)))
	buffer.WriteString(fmt.Sprintf("告警描述: \"我挂了，快来救我^OO^\"\n"))
	//buffer.WriteString(fmt.Sprintf("mentioned_mobile_list: %v\n",msgres["mentioned_mobile_list"]))
	buffer.WriteString(fmt.Sprintf("Status:%v\n",notification.Status))
	// 恢复消息
	var buffer2 bytes.Buffer
	buffer2.WriteString(fmt.Sprintf("告警: %v\n",string(msg)))
	buffer2.WriteString(fmt.Sprintf("Endpoint: %v\n",string(msg1)))
	buffer2.WriteString(fmt.Sprintf("告警描述: \"哈哈哈，我又回来了！！！\"\n"))
	//buffer2.WriteString(fmt.Sprintf("mentioned_mobile_list: %v\n",msgres["mentioned_mobile_list"]))
	buffer2.WriteString(fmt.Sprintf("Status:%v\n",notification.Status))
	//"mentioned_mobile_list": ["18270117154"]
	var m Message
	m.MsgType = "text"
	m.Text.Mentioned_mobile_list = msgres["mentioned_mobile_list"]
	if notification.Status == "resolved" {
		m.Text.Content = buffer2.String()
	}else if notification.Status == "firing" {
		m.Text.Content = buffer.String()
	}
	jsons, err := json.Marshal(m)
	if err != nil {
		log.Println("SendMessage Marshal failed,", err)
		return
	}
	resp := string(jsons)
	client := &http.Client{}

	req, err := http.NewRequest("POST", defaultRobot, strings.NewReader(resp))
	if err != nil {
		log.Println("SendMessage http NewRequest failed,", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		log.Println("SendMessage client Do failed", err)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("SendMessage ReadAll Body failed", err)
		return
	}
	log.Println("SendMessage success,body:", string(body))
}

func Alter(c *gin.Context)  {
	var notification Notification

	err := c.BindJSON(&notification)
	//fmt.Printf("%#v",notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	SendMessage(notification,defaultRobot)
}

func main()  {
	t := gin.Default()
	t.POST("/Alter",Alter)
	t.Run(":8090")
}
