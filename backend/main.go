package main

import (
	"gotoraft/cmd/bootstrap"
	"gotoraft/pkg/logger"
)

func main() {
	// 初始化各个组件
	app := bootstrap.NewApp()
	if err := app.Init(); err != nil {
		logger.Fatal("Failed to initialize application:", err)
	}
	// 关闭应用程序
	defer app.Shutdown()
	// 启动后端服务器
	if err := app.Run(); err != nil {
		logger.Fatal("Error running application:", err)
	}

}
