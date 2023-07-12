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
	environment := msgStruct.Environment
	projectID := env[cluster + "-" + environment]
	imageName := strings.Split(msgStruct.Image, ":")[0]
	newTag := strings.Split(msgStruct.Image, ":")[1]
	transID := newTag
	blobList := locateBlob(projectID, newTag, imageName)
	oldTag := getOldTag(projectID, newTag, blobList, imageName)

	if oldTag != "" && oldTag != newTag {
		logger.Infof("[%s] Comparing old tag [%s] and new tag [%s]...", newTag, oldTag, newTag)
		for blob := 0; blob < len(blobList); blob++ {
			if changeTagImage(projectID, transID, environment, imageName, oldTag, newTag, blobList[blob], botName, rootPath) {
				logger.Infof("[%s][%d] Changing old tag [%s] to new tag [%s] successfully", transID, blob, oldTag, newTag)
			}
		}
		notifyTelegram(cluster, environment, imageName, oldTag, newTag)
	} else {
		logger.Warnf("[%s] Old tag is not found OR nothing to change", newTag)
	}
}