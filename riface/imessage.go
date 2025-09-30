package riface

// 消息模块通过长度字段读取数据，解决粘包问题
type IMessage interface {
	// 获取消息id
	GetMsgId() uint32
	// 获取消息长度
	GetDataLen() uint32
	// 获取消息内容
	GetMsgData() []byte
	// 设置消息id
	SetMsgId(id uint32)
	// 设置消息长度
	SetDataLen(len uint32)
	// 设置消息内容
	SetMsgData(data []byte)
}
