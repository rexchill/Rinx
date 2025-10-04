package riface

// IConnManager 连接模块的接口
type IConnManager interface {
	// Add 增加一个连接
	Add(conn IConnection)
	// Remove 从连接管理模块删除一个连接，但没有关闭该连接
	Remove(conn IConnection)
	// Get 根据连接ID获取某个连接
	Get(Id uint32) (IConnection, error)
	// Len 查询总的连接个数
	Len() int
	// ClearAll 从连接管理模块删除所有的连接，并将连接关闭(conn.Stop()方法)
	ClearAll()
}
