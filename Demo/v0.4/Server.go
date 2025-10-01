package main

import (
	"Rinx/riface"
	"Rinx/rnet"
	"fmt"
)

// PingRouter 程序员写的业务方法
type PingRouter struct {
	rnet.BaseRouter
}

func (r *PingRouter) PreHandler(request riface.IRequest) {
	fmt.Println("方法处理前...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("开始ping..."))
	if err != nil {
		fmt.Println("前置处理逻辑失败...", err)
	}
}

func (r *PingRouter) Handler(request riface.IRequest) {
	fmt.Println("开始处理业务...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping成功..."))
	if err != nil {
		fmt.Println("业务处理失败...", err)
	}
}

func (r *PingRouter) PostHandler(request riface.IRequest) {
	fmt.Println("方法处理后...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("完成ping..."))
	if err != nil {
		fmt.Println("后置处理逻辑失败...", err)
	}

}
func main() {
	// 创建新的服务器
	server := rnet.NewServer("[Rinx v0.4]")
	// 注册路由方法
	server.AddRouter(&PingRouter{})
	// 运行服务器
	server.Serve()
}
