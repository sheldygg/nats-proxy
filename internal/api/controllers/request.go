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
	Timeout int                 `json:"timeout" binding:"required"`
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

	response, err := ctx.Nats.RequestMsg(msg, time.Duration(request.Timeout)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "response": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"ok": true, "response": string(response.Data)})
	}
}
