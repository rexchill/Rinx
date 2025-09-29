package rnet

import "Rinx/riface"

type Request struct {
	// 建立的连接
	conn riface.IConnection
	// 客户端的请求信息
	data []byte
}

// 接口类型就是指针类型
func (req *Request) GetConnection() riface.IConnection {
	return req.conn
}

func (req Request) GetData() []byte {
	return req.data
}
