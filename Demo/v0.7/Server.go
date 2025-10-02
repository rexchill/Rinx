package main

import (
	"Rinx/riface"
	"Rinx/rnet"
	"fmt"
)

type Hello struct {
	rnet.BaseRouter
}

func (r *Hello) Handler(request riface.IRequest) {
	fmt.Println("开始处理业务...")
	fmt.Println("接收到来自客户端的msgId", request.GetMessage().GetMsgId(), "msgData", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("Hello Rinx..."))
	if err != nil {
		fmt.Println("业务处理失败...", err)
	}
}

// PingRouter 程序员写的业务方法
type PingRouter struct {
	rnet.BaseRouter
}

func (r *PingRouter) Handler(request riface.IRequest) {
	fmt.Println("开始处理业务...")
	fmt.Println("接收到来自客户端的msgId", request.GetMessage().GetMsgId(), "msgData", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("ping成功..."))
	if err != nil {
		fmt.Println("业务处理失败...", err)
	}
}

func main() {
	// 创建新的服务器
	server := rnet.NewServer("[Rinx v0.7]")
	// 注册路由方法
	server.AddRouter(0, &Hello{})
	server.AddRouter(1, &PingRouter{})
	// 运行服务器
	server.Serve()
}
