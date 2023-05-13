package cqhttp

import (
	"Kira-qbot/model"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"
)

var (
	SendMsgClient http.Client
	SendMsgConfig model.Config
)

func init() {

	configFileName := "config/Config.yml"
	configFile, _ := os.Open(configFileName)
	decoder := yaml.NewDecoder(configFile)
	_ = decoder.Decode(&SendMsgConfig)

	SendMsgClient = http.Client{}
}

func SendMsg(event model.Event, message string) {
	payload := strings.NewReader(fmt.Sprintf("user_id=%d&message=%s", event.UserId, message))

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/send_private_msg", SendMsgConfig.CqHttpHost), payload)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, _ = SendMsgClient.Do(request)
}
