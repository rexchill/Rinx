package riface

// IRouter 路由接口，根据不同的请求 路由到不同的业务处理逻辑
type IRouter interface {
	// PreHandler 业务处理前的操作
	PreHandler(req IRequest)
	// Handler 进入业务
	Handler(req IRequest)
	// PostHandler 业务处理后的操作
	PostHandler(req IRequest)
}
