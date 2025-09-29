package riface

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 路由功能，给当前服务注册一个路由业务方法，供客户端连接使用
	AddRouter(router IRouter)
}
