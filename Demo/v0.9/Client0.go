package main

import (
	"Rinx/rnet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("客户端开始连接...")
	time.Sleep(time.Second)
	// 连接到服务器
	conn, err := net.Dial("tcp4", "127.0.0.1:19991")
	if err != nil {
		fmt.Println("连接失败...", err)
		return
	}
	for {
		// 连接成功
		// 封包
		pkg := rnet.NewDataPackage()
		msgSend := rnet.NewMessage(0, []byte("Client0..."))
		binaryMsg, err := pkg.Pack(msgSend)
		if err != nil {
			fmt.Println("客户端封装数据失败...", err)
			continue
		}

		// 发送
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("发送数据失败...", err)
			continue
		}

		// 读取数据头部
		header := make([]byte, pkg.GetHeadLen())
		if _, err := io.ReadFull(conn, header); err != nil {
			fmt.Println("读取数据头部失败...", err)
			continue
		}
		// 解析数据头部
		msgRecv, err := pkg.UnPack(header)
		if err != nil {
			fmt.Println("解包数据头部失败...", err)
			continue
		}

		// 读取数据内容
		body := make([]byte, msgRecv.GetDataLen())
		if _, err := io.ReadFull(conn, body); err != nil {
			fmt.Println("读取数据内容失败...", err)
			continue
		}
		// 解析数据内容
		msgRecv.SetMsgData(body)

		// 接受数据成功，打印验证
		fmt.Println("==> Recv Msg: ID=", msgRecv.GetMsgId(), ", len=", msgRecv.GetDataLen(), ", data=", string(msgRecv.GetMsgData()))
		time.Sleep(5 * time.Second)
	}
}
