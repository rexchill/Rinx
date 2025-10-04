package rnet

import (
	"Rinx/riface"
	"Rinx/utils"
	"fmt"
	"net"
)

// Server IServer接口实现，定义一个Server的服务器模块
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
	// 连接管理模块
	ConnManager riface.IConnManager
	// 钩子函数，连接建立时自动实施的方法
	OnConnStart func(conn riface.IConnection)
	// 钩子函数，连接关闭时自动实施的方法
	OnConnStop func(conn riface.IConnection)
}

// Start 开启Server
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Rinx] Version: %s, MaxConnection: %d,  MaxPacketSize: %d\n",
		utils.Config.Version,
		utils.Config.MaxConnection,
		utils.Config.MaxPacketSize)
	go func() {
		// 0、服务启动前开启线程池，只允许启动一个线程池，即只运行一次
		s.MsgHandler.StartWorkerPool()

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

			// 判断当前连接个数是否超过了最大连接数，超过了就不允许此次连接了
			if s.ConnManager.Len() >= utils.Config.MaxConnection {
				// TODO 关闭连接前可以给用户一些反馈信息
				fmt.Println("======> 当前连接数过多，请稍后在试...连接个数为： ", s.ConnManager.Len())
				if err := conn.Close(); err != nil {
				}
				continue

			}
			// 将监听到的TCP连接封装到自己构建的连接模块中，便于调用不同的业务方法
			dealConn := NewConnection(s, conn, connId, s.MsgHandler)
			connId++
			// 开启协程进行业务处理
			go dealConn.Start()
		}
	}()

}

// Stop 停止Server
func (s *Server) Stop() {
	fmt.Println("[Rinx服务器正在关闭...], 服务器名称为： ", s.Name)
	s.ConnManager.ClearAll()
}

// Serve 运行Server
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

// NewServer 新建Server实现
func NewServer(name string) riface.IServer {
	newServer := &Server{
		Name:        utils.Config.Name,
		IPVersion:   "tcp4",
		IP:          utils.Config.Host,
		Port:        utils.Config.Port,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return newServer
}

func (s *Server) GetConnManager() riface.IConnManager {
	return s.ConnManager
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(conn riface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(conn riface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn riface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("连接建立时进行的钩子方法...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn riface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("连接建立时进行的钩子方法...")
		s.OnConnStop(conn)
	}
}
