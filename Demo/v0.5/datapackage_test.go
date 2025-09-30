package v0_5

import (
	"Rinx/rnet"
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试dataPackage的封包、拆包的单元测试
func TestDataPackage(t *testing.T) {

	// 服务器注册监听套接字
	listen, err := net.Listen("tcp", "127.0.0.1:19991")
	if err != nil {
		fmt.Println("监听创建失败...", err)
		return
	}

	go func(listener net.Listener) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("连接建立失败...等待新连接...", err)
				continue
			}

			// 开启新协程去处理这个新的连接
			go func(conn net.Conn) {
				// 包对象
				pkg := rnet.NewDataPackage()
				for {
					// TODO 待后续了解
					// 突然发现，如果采用这种package的方式，如果一次读取出现问题，那么后续所有基于此连接的数据都会发生错误
					// 比如某一次读取包头header出现错误，再次循环读取的话，就会把 数据体MsgData 读成包头MsgId和DataLen，然后造成一连串的错误
					// 读取发生错误时先直接return，结束此次连接

					// 读取包头
					header := make([]byte, pkg.GetHeadLen())
					_, err := io.ReadFull(conn, header)
					if err != nil {
						fmt.Println("从缓冲区读取流数据失败...")
						// 由于上述所写的原因，所有直接return
						return
					}
					// 解析包头，写入message
					msHeader, err := pkg.UnPack(header)
					if err != nil {
						fmt.Println("解析包头失败...", err)
						// 由于上述所写的原因，所有直接return
						return
					}
					// 类型断言，将接口转为实现类
					msg := msHeader.(*rnet.Message)
					msg.MsgData = make([]byte, msg.GetDataLen())
					_, err = io.ReadFull(conn, msg.MsgData)
					if err != nil {
						fmt.Println("读取数据体失败...", err)
						return
					}

					// 读取成功后的数据处理...
					fmt.Println("==> Recv Msg: MsgId=", msg.MsgId, ", DateLen=", msg.DateLen, ", MsgData=", string(msg.MsgData))
				}

			}(conn)

		}

	}(listen)

	// 客户端
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:19991")
		if err != nil {
			fmt.Println("连接服务器失败...", err)
			return
		}

		pkg := rnet.NewDataPackage()
		msg1 := &rnet.Message{
			MsgId:   1,
			DateLen: 5,
			MsgData: []byte("hello"),
			//MsgData: []byte{'h', 'e', 'l', 'l', 'o'},
		}
		msg2 := &rnet.Message{
			MsgId:   2,
			DateLen: 9,
			MsgData: []byte("loveChina"),
			//MsgData: []byte{'l', 'o', 'v', 'e', 'C', 'h', 'i', 'n', 'a'},
		}

		// 打包
		sendData1, err := pkg.Pack(msg1)
		if err != nil {
			fmt.Println("打包失败...", err)
			return
		}
		sendData2, err := pkg.Pack(msg2)
		if err != nil {
			fmt.Println("打包失败...", err)
			return
		}

		// 两个包的数据粘在一起
		sendData1 = append(sendData1, sendData2...)

		// 发送
		_, err = conn.Write(sendData1)
		if err != nil {
			fmt.Println("发送数据失败...", err)
		}
	}()
	select {}

}
