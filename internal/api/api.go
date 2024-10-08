package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sheldygg/nats-proxy/internal/api/router"
	"github.com/sheldygg/nats-proxy/internal/pkg/config"
)

func setConfiguration(configPath string) {
	config.Setup(configPath)
	gin.SetMode(config.GetConfig().Server.Mode)
}

func Run(configPath string) {
	if configPath == "" {
		configPath = "data/config.yml"
	}
	setConfiguration(configPath)
	conf := config.GetConfig()
	web := router.Setup()
	fmt.Println("Go API REST Running on port " + conf.Server.Port)
	fmt.Println("==================>")
	_ = web.Run(":" + conf.Server.Port)
}
