package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sheldygg/nats-proxy/internal/pkg/config"
	"net/http"
)

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": config.GetConfig()})
}
