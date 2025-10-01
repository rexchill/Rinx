package rnet

import (
	"Rinx/riface"
	"errors"
	"fmt"
	"io"
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
		// 进行处理的拆包对象
		pkg := NewDataPackage()

		// 读取头部字段的二进制流
		header := make([]byte, pkg.GetHeadLen())
		_, err := io.ReadFull(conn.GetTCPConnection(), header)
		if err != nil {
			fmt.Println("读取pkg头部数据失败...", err)
			return
		}

		// 将头部字段数据填入Message，即进行拆包处理
		msg, err := pkg.UnPack(header)
		if err != nil {
			fmt.Println("message头部解析失败...", err)
		}

		// 根据头部字段中的长度读取对应的数据
		msgData := make([]byte, msg.GetDataLen())
		_, err = io.ReadFull(conn.GetTCPConnection(), msgData)
		if err != nil {
			fmt.Println("读取数据msgData失败...", err)
			return
		}

		// 将msgData数据填入Message
		msg.SetMsgData(msgData)

		// 将客户端的请求消息和对应的连接封装在一起
		req := Request{
			conn:    conn,
			message: msg,
		}
		// 执行对应连接注册的路由方法(程序员定义的执行方法)
		go func(req riface.IRequest) {
			conn.router.PreHandler(req)
			conn.router.Handler(req)
			conn.router.PostHandler(req)
		}(&req)
	}
}

func (conn *Connection) SendMsg(msgId uint32, msgData []byte) error {
	// 连接关闭，给出关闭信息并返回
	if conn.isClosed {
		return errors.New("连接已关闭...")
	}

	// 封包
	binaryMsg, err := NewDataPackage().Pack(NewMessage(msgId, msgData))
	if err != nil {
		fmt.Println("封包失败...", err)
		return err
	}

	// 发送给客户端
	_, err = conn.GetTCPConnection().Write(binaryMsg)
	if err != nil {
		fmt.Println("发送二进制数据失败...", err)
	}

	return nil
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
