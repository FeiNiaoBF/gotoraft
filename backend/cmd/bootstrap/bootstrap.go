// 应用程序引导程序包，负责初始化各个组件
package bootstrap

import (
	"fmt"
	"gotoraft/config"
	"gotoraft/internal/router"
	"gotoraft/pkg/logger"
	"gotoraft/pkg/store"
	"gotoraft/pkg/websocket"
)

// App 表示应用程序实例
type App struct {
	config    *config.Config
	router    *router.Router
	wsManager *websocket.Manager
	store     store.Store
}

// NewApp 创建一个新的 App 实例
func NewApp() *App {
	return &App{}
}

// Init 初始化后端系统
func (app *App) Init() error {
	// 1. 加载配置
	if err := app.initConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %v", err)
	}

	// 2. 初始化日志系统
	if err := app.initLogger(); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// 3. 初始化WebSocket管理器
	app.initWebSocket()

	// 4. 初始化HTTP路由
	app.initRouter()

	return nil
}

// initConfig 初始化配置
func (app *App) initConfig() error {
	err := config.Init()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	app.config = config.GetConfig()
	return nil
}

// initLogger 初始化日志系统
func (app *App) initLogger() error {
	if err := logger.InitLogger(); err != nil {
		return err
	}

	logger.Infof("应用程序正在初始化... ")

	return nil
}

// initWebSocket 初始化WebSocket管理器
func (app *App) initWebSocket() {
	// 初始化WebSocket管理器
	app.wsManager = websocket.NewManager()
	logger.Info("WebSocket管理器已初始化")
}

// initRouter 初始化HTTP路由
func (app *App) initRouter() {
	// 创建路由管理器，并传入WebSocket管理器
	app.router = router.NewRouter(app.wsManager)
	// 注册所有路由
	app.router.RegisterRoutes()
	logger.Info("HTTP路由和WebSocket已初始化")
}

// Run 运行应用程序
func (app *App) Run() error {
	addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port)
	logger.Infof("应用程序启动成功，监听地址: %s", addr)
	return app.router.Engine().Run(addr)
}
