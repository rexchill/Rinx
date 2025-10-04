package main

import (
	"Rinx/riface"
	"Rinx/rnet"
	"fmt"
	"time"
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

func doConnBuilderPre(c riface.IConnection) {
	fmt.Println("执行连接建立时的钩子方法")
	if err := c.SendMsg(c.GetConnID(), []byte("新连接执行钩子方法")); err != nil {
		fmt.Println(err)
	}
	// 为连接设置属性值
	c.SetProperty("爱好", "羽毛球")
	c.SetProperty("登入时间", time.Now())

}

func doConnBuilderPost(c riface.IConnection) {
	fmt.Println("执行连接结束时的钩子方法")
	// 连接结束时，可以取出属性值
	if hobby, err := c.GetProperty("爱好"); err == nil {
		fmt.Println("爱好为：", hobby)
	}
	if loginTime, err := c.GetProperty("登入时间"); err == nil {
		fmt.Println("登入时间为：", loginTime)
	}

}

func main() {
	// 创建新的服务器
	server := rnet.NewServer("[Rinx v0.7]")
	// 注册路由方法
	server.AddRouter(0, &Hello{})
	server.AddRouter(1, &PingRouter{})
	// 注册钩子方法
	server.SetOnConnStart(doConnBuilderPre)
	server.SetOnConnStop(doConnBuilderPost)
	// 运行服务器
	server.Serve()
}
