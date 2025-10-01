package main

import (
	"fmt"
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
		// 连接成功，传输数v
		_, err := conn.Write([]byte("Hello Rinx v0.4..."))
		if err != nil {
			fmt.Println("客户端发送数据失败...", err)
			return
		}
		// 接收数据
		recv := make([]byte, 512)
		count, err := conn.Read(recv)
		if err != nil {
			fmt.Println("接收数据失败...", err)
			return
		}
		// 接受数据成功，打印验证
		fmt.Printf("接收数据为： %s  ,接收字节数为 %d \n", recv[:count], count)
		time.Sleep(5 * time.Second)
	}
}
