package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"vngitSub/pkg/utils"
)

func notifyTelegram(cluster string, env string, imageName string, oldTag string, newTag string) {
	logger := utils.ConfigZap()
	telegramToken := os.Getenv("TELEGRAMTOKEN")
	telegramChannel := os.Getenv("TELEGRAMCHANNEL")
	service := strings.Split(imageName, "/")[len(strings.Split(imageName, "/"))-1]

	var channel map[string]string
	json.Unmarshal([]byte(telegramChannel), &channel)

	telegramMessage := "<b>VNGITBOT-V2 has changed version tag for deployment.</b>%0A<b>Service</b>: <code>" + service + "</code>%0A<b>Cluster</b>: <code>" + cluster + "</code>%0A<b>Old tag</b>: <code>" + oldTag + "</code>  ==>  <b>New tag</b>: <code>" + newTag + "</code>"

	telegramURL := "https://api.telegram.org/bot" + telegramToken + "/sendMessage?chat_id=" + channel[env] + "&parse_mode=HTML&text=" + telegramMessage

	res, err := http.Get(telegramURL)
	if err != nil {
		logger.Errorf("[%s] Sending alert notification to Telegram...%s: %s", newTag, res.Status, err)
	}

	logger.Infof("[%s] Sending alert notification to Telegram...%s", newTag, res.Status)
}