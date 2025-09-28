package rnet

import (
	"Rinx/riface"
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
}

// 开启
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at %s:%d\n", s.IP, s.Port)
	go func() {
		// 一、获取服务器套接字(ip:port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("s%:%d", s.IP, s.Port))
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
		for {
			// 接收客户端连接
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("连接客户端失败", err)
				continue
			}
			// 客户端连接服务端成功，服务端开始执行任务并返回给客户端
			// TODO
			go func() {
				for {
					// 暂时执行回显
					readBuf := make([]byte, 512)
					writeBuf := make([]byte, 512)
					count, err := conn.Read(readBuf) // 读取字节放入buf中，并返回读取成功的字节数
					if err != nil {
						fmt.Println("读缓冲区失败(从客户端连接读取字节数据失败)", err)
						continue
					}
					// 成功读取数据后执行业务逻辑
					if _, err := conn.Write(writeBuf[:count]); err != nil {
						fmt.Println("写缓冲区失败(发送给客户端数据失败)", err)
						continue
					}

				}

			}()
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

// 新建Server实现
func NewServer(name string) riface.IServer {
	newServer := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      19991,
	}
	return newServer
}
