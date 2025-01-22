package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// 启动服务器 main

func main() {
	// 使用Gin来启动一个简单的服务器
	server := gin.Default()
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Println("Server is running on port 8080")
	server.Run()
}
