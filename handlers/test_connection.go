package handlers

import (
	"github.com/gin-gonic/gin"
	"go-mongo-api/config"
	"net/http"
)

func TestConnection(c *gin.Context) {
	err := config.MongoDB.Client().Ping(c, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to connect to MongoDB"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully connected to MongoDB"})
}
