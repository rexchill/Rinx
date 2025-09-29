package riface

// IRequest 将建立的连接和客户端的请求信息封装在一起
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
}
