package controllers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"time"
)

type RequestData struct {
	Subject string              `json:"subject" binding:"required"`
	Header  map[string][]string `json:"headers" binding:"required"`
	Data    string              `json:"data" binding:"required"`
	Timeout float64             `json:"timeout" binding:"required"`
}

type AppContext struct {
	Nats *nats.Conn
}

func (ctx *AppContext) Request(c *gin.Context) {
	var request RequestData
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	var msg = nats.NewMsg(request.Subject)
	msg.Header = request.Header
	decodedData, err := base64.StdEncoding.DecodeString(request.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}
	msg.Data = decodedData

	log.Printf("Received Msg: Subject: %s, Header: %v, Data: %s", msg.Subject, msg.Header, string(msg.Data))

	timeoutDuration := time.Duration(request.Timeout * float64(time.Second))
	response, err := ctx.Nats.RequestMsg(msg, timeoutDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"response": gin.H{
			"subject": response.Subject,
			"reply":   response.Reply,
			"data":    response.Data,
			"header":  response.Header,
		},
	})
}
