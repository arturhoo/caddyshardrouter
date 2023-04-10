package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		customer := c.Request.Header.Get("X-Customer")
		c.Header("X-Shard", os.Getenv("SHARD"))
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello %s!", customer),
		})
	})

	r.POST("/", func(c *gin.Context) {
		var json struct {
			Customer string `json:"customer"`
		}
		if c.Bind(&json) == nil {
			c.Header("X-Shard", os.Getenv("SHARD"))
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Hello %s!", json.Customer),
			})
		}
	})

	r.Run()
}
