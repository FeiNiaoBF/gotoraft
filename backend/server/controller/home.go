package controller

import (
	"github.com/gin-gonic/gin"
)

// HomeController 处理主页请求
func HomeController(c *gin.Context) {
	c.String(200, "Welcome to the Raft Visualization Home Page!")
}
 