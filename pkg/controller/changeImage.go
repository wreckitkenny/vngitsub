package controller

import (
	"encoding/json"
	"os"
	"strings"

	"vngitSub/model"
	"vngitSub/pkg/utils"
)

// ChangeImage changes image tag in Gitlab projects
func ChangeImage(msg []byte) {
	logger := utils.ConfigZap()

	var msgStruct model.Message
	unmarshalErr := json.Unmarshal(msg, &msgStruct)
	if unmarshalErr != nil {
		logger.Errorf("Converting message to object...FAILED: %s", unmarshalErr)
	}

	var env map[string]int
	clusterenv := os.Getenv("CLUSTERENV")
	json.Unmarshal([]byte(clusterenv), &env)

	botName := os.Getenv("BOTNAME")
	rootPath := os.Getenv("ROOTPATH")
	cluster := msgStruct.Cluster
	environment := strings.Split(cluster, "-")[len(strings.Split(cluster, "-"))-1]
	projectID := env[cluster]
	imageName := strings.Split(msgStruct.Image, ":")[0]
	newTag := strings.Split(msgStruct.Image, ":")[1]
	transID := newTag
	blobList := locateBlob(projectID, newTag, imageName)
	oldTag, err := getOldTag(projectID, newTag, blobList, imageName)
	if err != "" {
		logger.Warnf("[%s] %s", transID, err)
	} else {
		logger.Infof("[%s] Getting old tag [%s]...OK", transID, oldTag)
	}

	if oldTag != "" && oldTag != newTag {
		logger.Infof("[%s] Comparing old tag [%s] and new tag [%s]", newTag, oldTag, newTag)
		changeState := changeTagImage(projectID, transID, environment, imageName, oldTag, newTag, blobList, botName, rootPath, msgStruct.Image, msgStruct.Cluster)
		if changeState != "SUCCESSFUL" {
			logger.Errorf("[%s] Changing old tag [%s] to new tag [%s]...%s", transID, oldTag, newTag, changeState)
		} else {
			logger.Infof("[%s] Changing old tag [%s] to new tag [%s]...%s", transID, oldTag, newTag, changeState)
			notifyTelegram(cluster, environment, imageName, oldTag, newTag)
		}
	} else {
		logger.Warnf("[%s] Old tag is not found OR nothing to change", newTag)
	}
}