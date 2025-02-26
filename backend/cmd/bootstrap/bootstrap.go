// 应用程序引导程序包，负责初始化各个组件
package bootstrap

import (
	"fmt"
	"gotoraft/config"
	"gotoraft/internal/kvstore/store"
	"gotoraft/internal/observer"
	"gotoraft/internal/router"
	"gotoraft/internal/websocket"
	"gotoraft/pkg/logger"
	"time"
)

// App 表示应用程序实例
type App struct {
	config    *config.Config
	router    *router.Router
	wsManager *websocket.Manager
	store     *store.Store                // kv存储
	observer  *observer.RaftStateObserver // Raft状态观察器
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

	// 3. 初始化存储
	if err := app.initStore(); err != nil {
		return fmt.Errorf("failed to initialize store: %v", err)
	}

	// 4. 初始化WebSocket管理器
	app.initWebSocket()

	// 5. 初始化HTTP路由
	app.initRouter()

	// 6. 初始化状态观察器
	app.initStateObserver()

	return nil
}

// initConfig 初始化配置
func (app *App) initConfig() error {
	err := config.Init()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	app.config = config.GetConfig()
	logger.Info("配置已加载...!")
	return nil
}

// initLogger 初始化日志系统
func (app *App) initLogger() error {
	if err := logger.InitLogger(); err != nil {
		return err
	}

	logger.Infof("应用程序正在初始化...!")

	return nil
}

// initStore 初始化存储
func (app *App) initStore() error {
	peers := []string{"node1", "node2", "node3"} // 示例节点
	app.store = store.NewStore(peers, "node1")
	return nil
}

// initWebSocket 初始化WebSocket管理器
func (app *App) initWebSocket() {
	// 初始化WebSocket管理器
	app.wsManager = websocket.NewManager(websocket.Config{
		MaxConnections:   100,              // 设置最大连接数
		HeartbeatTimeout: 30 * time.Second, // 设置心跳超时时间
	})
	logger.Info("WebSocket管理器已初始化...!")
}

// initStateObserver 初始化Raft状态观察器
func (app *App) initStateObserver() error {
	app.observer = observer.NewRaftStateObserver(
		app.store,
		app.wsManager,
	)
	// 启动状态观察
	go app.observer.Start()
	logger.Info("Raft状态观察器已启动...!")
	return nil
}

// initRouter 初始化HTTP路由
func (app *App) initRouter() {
	// 创建路由需要传递所有依赖组件
	app.router = router.NewRouter(
		app.wsManager,
		app.store,
		app.observer,
	)

	// 注册路由
	app.router.RegisterRoutes()

	logger.Info("HTTP路由已初始化")
}

// Run 运行应用程序
func (app *App) Run() error {
	addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port)
	logger.Infof("应用程序启动成功，监听地址: %s", addr)
	return app.router.Run(addr)
}

func (app *App) Shutdown() {
	// 关闭顺序与初始化顺序相反
	app.observer.Stop()
	app.store.Shutdown()
	app.wsManager.Shutdown()
}

func (app *App) initObserver() {
	app.observer = observer.NewObserver(app.store)
	go app.observer.Start()
}
