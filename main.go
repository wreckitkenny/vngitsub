package main

import (
	"vngitSub/pkg/controller"
	"vngitSub/pkg/utils"

	"github.com/gin-gonic/gin"
)

const version string = "1.0.0"

func main(){
	//Load Default Config
	// defaultConfig, err := utils.LoadConfig("config")
	logger := utils.ConfigZap()
	// if err != nil {
	// 	logger.Errorf("Loading configuration...failed: %s", err)
	// } else {
	// 	logger.Debug("Loading configuration...ok")
	// }

	//Configure GIN
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetTrustedProxies(nil)
	logger.Infof("Loading version...%s", version)
	controller.MsgHandler()
}