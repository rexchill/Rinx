package riface

import "net"

// IConnection 连接模块的抽象方法
type IConnection interface {

	// Start 开启连接-->让当前的连接开始工作
	Start()

	// Stop 停止连接-->结束当前连接的工作
	Stop()

	// GetTCPConnection 获取与客户端进行TCP连接的 连接对象
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接id
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的套接字(ip:port)
	RemoteAddr() net.Addr

	// SendMsg 发送数据给客户端
	SendMsg(uint32, []byte) error

	// SetProperty 设置连接属性
	SetProperty(string, any)

	// GetProperty 获取连接属性
	GetProperty(string) (any, error)

	// Remove 移除连接属性
	Remove(string)
}

// HandleFunc 定义一个处理业务的方法
// conn ==> 与客户端进行TCP连接的 连接对象
// data ==> 处理业务的数据
// count ==> 数据的长度
type HandleFunc func(conn *net.TCPConn, data []byte, count int) error
