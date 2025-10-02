package rnet

import (
	"Rinx/riface"
	"Rinx/utils"
	"fmt"
	"net"
)

// IServer接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器通信的ip版本(ipv4 or ipv6)
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// 当前Server的消息管理模块，用来绑定msgId和对应的业务逻辑
	MsgHandler riface.IMsgHandler
}

// 开启
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Rinx] Version: %s, MaxConnection: %d,  MaxPacketSize: %d\n",
		utils.Config.Version,
		utils.Config.MaxConnection,
		utils.Config.MaxPacketSize)
	go func() {
		// 一、获取服务器套接字(ip:port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("获取套接字失败(addr获取失败)", err)
			return
		}
		// 二、监听服务器套接字，看是否有服务端连接进来
		tcpListener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("监听套接字失败(tcpListener获取失败)", err)
			return
		}
		// 三、阻塞等待客户端连接，执行任务
		var connId uint32
		connId = 0
		for {
			// 接收客户端连接
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("连接客户端失败", err)
				continue
			}
			// 将监听到的TCP连接封装到自己构建的连接模块中，便于调用不同的业务方法
			dealconn := NewConnection(conn, connId, s.MsgHandler)
			connId++
			// 开启协程进行业务处理
			go dealconn.Start()
		}
	}()

}

// 停止
func (s *Server) Stop() {
	// TODO 停止服务，做一些资源的释放

}

// 运行
func (s *Server) Serve() {
	// 启动服务
	s.Start()
	// TODO 做一些额外的功能，比如将服务器注册到zookeeper里面
	// 阻塞，不然Serve调用结束就停止了，Start也会跟着停止
	select {}

}

func (s *Server) AddRouter(msgId uint32, router riface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

// 新建Server实现
func NewServer(name string) riface.IServer {
	newServer := &Server{
		Name:       utils.Config.Name,
		IPVersion:  "tcp4",
		IP:         utils.Config.Host,
		Port:       utils.Config.Port,
		MsgHandler: NewMsgHandler(),
	}
	return newServer
}
