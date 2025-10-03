package rnet

import (
	"Rinx/riface"
	"Rinx/utils"
	"errors"
	"fmt"
	"io"
	"net"
)

// Connection 连接模块
type Connection struct {

	// 与客户端进行TCP连接的 连接对象
	Conn *net.TCPConn

	// 连接的id
	ConnID uint32

	// 当前的连接状态
	isClosed bool

	// 当前连接消息类型对应的所要执行的业务方法
	MsgHandler riface.IMsgHandler

	// 用于通知当前连接已经退出的 channel；读协程获取到从客户端发来的退出信号，通知写协程退出
	ExitChan chan struct{}

	// 用于读写协程间的通信
	msgChan chan []byte
}

// NewConnection 初始化连接模块
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler riface.IMsgHandler) *Connection {
	con := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan struct{}),
		msgChan:    make(chan []byte),
	}
	return con
}
func (conn *Connection) StartReader() {
	fmt.Println("[读业务协程开始工作...]")
	defer conn.Stop()
	defer fmt.Println("读业务结束...连接ID为：", conn.ConnID, "远程客户端地址为： ", conn.RemoteAddr().String())

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
			return
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

		// 配置开启了线程池，发送到任务队列中进行集中处理
		if utils.Config.WorkerPoolSize > 0 {
			conn.MsgHandler.SendReqToTaskQueue(&req)
		} else { // 否则开启新协程进行处理
			// 执行对应连接注册的路由方法(程序员定义的执行方法)
			// 根据msgId找到对应注册的方法进行执行
			go conn.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (conn *Connection) StartWriter() {
	fmt.Println("[写业务协程开始工作...]")
	defer fmt.Println("[写协程退出...], 客户端地址为：", conn.RemoteAddr().String())
	for {
		select {
		case data := <-conn.msgChan:
			// TODO 应该在写协程进行业务的处理（？）
			// 接收读业务进行业务处理后的数据
			if _, err := conn.GetTCPConnection().Write(data); err != nil {
				fmt.Println("[写协程发送数据失败...]", err)
				return
			}
		case <-conn.ExitChan:
			// 客户端退出连接，结束写业务
			return
		}
	}

}

func (conn *Connection) SendMsg(msgId uint32, msgData []byte) error {
	// 连接关闭，给出关闭信息并返回
	if conn.isClosed {
		return errors.New("[连接已关闭...]")
	}

	// 封包
	binaryMsg, err := NewDataPackage().Pack(NewMessage(msgId, msgData))
	if err != nil {
		fmt.Println("封包失败...", err)
		return err
	}

	// 通过管道发给写业务
	conn.msgChan <- binaryMsg

	return nil
}
func (conn *Connection) Start() {
	fmt.Println("连接开始...连接ID为： ", conn.ConnID)
	// 启动该连接的读业务
	go conn.StartReader()
	// 启动该连接的写业务
	go conn.StartWriter()
}

func (conn *Connection) Stop() {
	fmt.Println("[执行Stop...]")
	// 连接已经关闭
	if conn.isClosed {
		return
	}
	// 否则，置关闭标识
	conn.isClosed = true
	// 服务端关闭连接
	if err := conn.GetTCPConnection().Close(); err != nil {
		fmt.Println("[服务端关闭连接失败...]", err)
	}
	// 告知Writer连接已关闭
	conn.ExitChan <- struct{}{}
	// 回收资源
	close(conn.ExitChan)
	close(conn.msgChan)
}

func (conn *Connection) GetTCPConnection() *net.TCPConn {
	return conn.Conn
}

func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}

func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

func (conn *Connection) Send(data []byte) error {
	return nil
}
