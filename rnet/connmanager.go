package rnet

import (
	"Rinx/riface"
	"fmt"
	"sync"
)

type ConnManager struct {
	// 所有连接的集合
	connMap map[uint32]riface.IConnection
	// 对连接集合进行操作的锁
	lock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]riface.IConnection),
		/*
			在 Go 中，sync.RWMutex 是一个值类型，并且它会有默认的零值。
			对于 sync.RWMutex，它的零值已经是一个有效的互斥锁实例。
			因此，在没有显式赋值时，sync.RWMutex 会被自动初始化为零值，即一个没有被锁定的锁。
			锁的零值已经是有效的，所有功能都可以正常使用。
		*/
	}
}

// Add 增加一个连接
func (cm *ConnManager) Add(conn riface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.connMap[conn.GetConnID()] = conn
	fmt.Println("连接添加成功...", "连接ID为 ", conn.GetConnID())
}

// Remove 删除一个连接
func (cm *ConnManager) Remove(conn riface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connMap, conn.GetConnID())
	fmt.Println("连接删除成功...", "连接ID为 ", conn.GetConnID(), "剩余连接个数为", cm.Len())

}

// Get 根据连接ID获取某个连接
func (cm *ConnManager) Get(Id uint32) (riface.IConnection, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	// 查询到对应连接
	if c, ok := cm.connMap[Id]; ok {
		return c, nil
	}
	// 未找到对应连接
	fmt.Println("没有该连接,请确认连接ID...")
	return nil, fmt.Errorf("没有该连接,请确认连接ID [%d]...", Id)

}

// Len 查询总的连接个数
func (cm *ConnManager) Len() int {
	return len(cm.connMap)
}

// ClearAll 清除所有的连接
func (cm *ConnManager) ClearAll() {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	for id, c := range cm.connMap {
		// 中止连接并释放资源
		c.Stop()
		// 删除连接
		delete(cm.connMap, id)
	}
	fmt.Println("清除所有连接成功，目前连接数 = ", cm.Len())

}
