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
)

// App 表示应用程序实例
type App struct {
	config    *config.Config
	router    *router.Router
	wsManager *websocket.Manager
	store     *store.Store
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

// initStore 初始化Raft存储
func (app *App) initStore() error {
	raftConfig := app.config.Raft
	if raftConfig == nil {
		return fmt.Errorf("raft configuration is missing")
	}

	// 创建存储实例
	s := store.New(false) // 使用持久化存储
	
	// 设置Raft存储路径和监听地址
	s.RaftDir = raftConfig.DataDir
	s.RaftBind = fmt.Sprintf("%s:%d", raftConfig.Host, raftConfig.Port)

	// 打开存储，如果是bootstrap节点则启用单节点模式
	if err := s.Open(raftConfig.Bootstrap, app.config.NodeID); err != nil {
		return fmt.Errorf("failed to open store: %v", err)
	}

	// 如果配置了join地址，则尝试加入集群
	if !raftConfig.Bootstrap && raftConfig.JoinAddr != "" {
		if err := s.Join(app.config.NodeID, raftConfig.JoinAddr); err != nil {
			return fmt.Errorf("failed to join cluster: %v", err)
		}
		logger.Info("成功加入Raft集群", 
			"nodeId", app.config.NodeID,
			"joinAddr", raftConfig.JoinAddr,
		)
	}

	app.store = s
	logger.Info("Raft存储初始化成功",
		"dataDir", s.RaftDir,
		"bindAddr", s.RaftBind,
		"bootstrap", raftConfig.Bootstrap,
	)
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

// initStateObserver 初始化Raft状态观察器
func (app *App) initStateObserver() {
	if app.store == nil || app.wsManager == nil {
		logger.Error("初始化状态观察器失败：store或wsManager未初始化")
		return
	}

	observer := observer.NewRaftStateObserver(app.store, app.wsManager, app.config.NodeID)
	observer.Start()

	logger.Info("Raft状态观察器初始化成功")
}

// Run 运行应用程序
func (app *App) Run() error {
	addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port)
	logger.Infof("应用程序启动成功，监听地址: %s", addr)
	return app.router.Engine().Run(addr)
}
