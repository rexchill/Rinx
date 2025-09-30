package rnet

import (
	"Rinx/riface"
	"Rinx/utils"
	"bytes"
	"encoding/binary"
	"errors"
)

type DataPackage struct{}

func NewDataPackage() riface.IDataPackage {
	return &DataPackage{}
}
func (pkg *DataPackage) GetHeadLen() uint32 {
	// MsgId:4字节	MsgLen:4字节
	return 8
}

func (pkg *DataPackage) Pack(msg riface.IMessage) ([]byte, error) {
	// 创建一个写入缓冲区
	buf := bytes.NewBuffer([]byte{})

	// 将message的id写入
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将message的长度写入
	if err := binary.Write(buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将message的数据写入
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}

func (pkg *DataPackage) UnPack(binaryData []byte) (riface.IMessage, error) {
	// 创建一个读取缓冲区
	buf := bytes.NewReader(binaryData)

	// 创捷一个message对象，用于解包装填后返回
	msg := &Message{}

	/*
		注意：msg是指向Message的指针，但是msg.MsgId等是通过msg这个指针访问到的Message字段的值，而不是msg.MsgId的地址，
		因此要用&msg.MsgId才能修改其变量

		binary.Read()第三个参数必须传入指针类型，以及必须是固定大小的类型(即结构体的大小必须固定，不能有变长字段)
	*/
	// 将id读取到message
	if err := binary.Read(buf, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}

	// 将长度读取到message
	if err := binary.Read(buf, binary.LittleEndian, &msg.DateLen); err != nil {
		return nil, err
	}

	// 判断dataLen的长度是否超出我们允许的最大包长度
	if utils.Config.MaxPacketSize > 0 && utils.Config.MaxPacketSize < msg.DateLen {
		return nil, errors.New("message包数据长度超过最大长度...")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，在外部再从conn读取一次数据
	return msg, nil

}
