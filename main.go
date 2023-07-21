package main

import (
	"os"
	"vngitSub/pkg/controller"
	"vngitSub/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main(){
	logger := utils.ConfigZap()

	//Configure GIN
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetTrustedProxies(nil)

	version, err := os.ReadFile("VERSION")
	if err != nil {
		logger.Errorf("Loading version...FAILED: %s", err)
	} else {
		logger.Infof("Loading version...%s", version)
	}
	controller.ValidateMongoConnection()
	controller.MsgHandler()
}