package riface

// TODO IMsgHandler和IRouter需要互换，应该是根据Router找MsgHandler
// 多路由策略，根据不同的消息类型路由到不同的业务处理方法
type IMsgHandler interface {
	// 执行业务逻辑
	DoMsgHandler(request IRequest)
	// 添加处理逻辑
	AddRouter(msgId uint32, router IRouter)
}
