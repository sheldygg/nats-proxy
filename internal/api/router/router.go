package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/sheldygg/nats-proxy/internal/api/controllers"
	"github.com/sheldygg/nats-proxy/internal/pkg/config"
	"io"
	"log"
	"os"
)

func Setup() *gin.Engine {
	app := gin.New()

	// Logging to a file.
	f, _ := os.Create("log/api.log")
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(f)

	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())

	conf := config.GetConfig()
	natsConn, err := nats.Connect(conf.NatsServers[0])
	if err != nil {
		log.Fatalf("error connect to nats: %v", err)
	}

	appCtx := controllers.AppContext{Nats: natsConn}

	app.GET("/getConfig", controllers.GetConfig)
	app.POST("/request", appCtx.Request)

	return app
}
