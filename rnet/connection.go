package rnet

import (
	"Rinx/riface"
	"fmt"
	"net"
)

// Connection 连接模块
type Connection struct {

	// 与客户端进行TCP连接的 连接对象
	conn *net.TCPConn

	// 连接的id
	ConnID uint32

	// 当前的连接状态
	isClosed bool

	// 当前连接所要执行的业务方法
	router riface.IRouter

	// 用于通知当前连接已经退出的 channel
	ExitChan chan struct{}
}

// NewConnection 初始化连接模块
func NewConnection(conn *net.TCPConn, connID uint32, router riface.IRouter) *Connection {
	con := &Connection{
		conn:     conn,
		ConnID:   connID,
		isClosed: false,
		router:   router,
		ExitChan: make(chan struct{}),
	}
	return con
}
func (conn *Connection) StartReader() {
	fmt.Println("读业务协程开始工作...")
	defer fmt.Println("读业务结束...连接ID为：", conn.ConnID, "远程客户端地址为： ", conn.RemoteAddr().String())
	defer conn.Stop()

	// 开始处理业务
	for {
		buf := make([]byte, 512)
		_, err := conn.conn.Read(buf)
		if err != nil {
			fmt.Println("读取数据失败...", err)
			continue // 继续尝试读取
		}
		// 将客户端的请求消息和对应的连接封装在一起
		req := Request{
			conn: conn,
			data: buf,
		}
		// 执行对应连接注册的路由方法(程序员定义的执行方法)
		go func(req riface.IRequest) {
			conn.router.PreHandler(req)
			conn.router.Handler(req)
			conn.router.PostHandler(req)
		}(&req)
	}
}

func (conn *Connection) Start() {
	fmt.Println("连接开始...连接ID为： ", conn.ConnID)
	// 启动该连接的读业务
	go conn.StartReader()
	// TODO 写业务
}

func (conn *Connection) Stop() {

}

func (conn *Connection) GetTCPConnection() *net.TCPConn {
	return conn.conn
}

func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}

func (conn *Connection) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

func (conn *Connection) Send(data []byte) error {
	return nil
}
