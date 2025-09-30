package riface

type IDataPackage interface {
	// 获取包头长度
	GetHeadLen() uint32
	// 将message打包成二进制流
	Pack(msg IMessage) ([]byte, error)
	// 解析二进制流为message
	UnPack(binaryData []byte) (IMessage, error)
}
