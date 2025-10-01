package rnet

// 消息模块通过长度字段读取数据，解决粘包问题 (另：还可以设置结束标志，标志消息的结束)
type Message struct {
	// 消息id
	MsgId uint32
	// 消息长度 解决粘包问题
	DateLen uint32
	// 消息内容
	MsgData []byte
}

func NewMessage(msgId uint32, msgData []byte) *Message {
	return &Message{
		MsgId:   msgId,
		DateLen: uint32(len(msgData)),
		MsgData: msgData,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

func (m *Message) GetDataLen() uint32 {
	return m.DateLen
}

func (m *Message) GetMsgData() []byte {
	return m.MsgData
}

func (m *Message) SetMsgId(id uint32) {
	m.MsgId = id
}

func (m *Message) SetDataLen(len uint32) {
	m.DateLen = len
}

func (m *Message) SetMsgData(data []byte) {
	m.MsgData = data
}
