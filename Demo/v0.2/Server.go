package main

import "Rinx/rnet"

func main() {
	// 创建新的服务器
	server := rnet.NewServer("[Rinx v0.2]")
	// 运行服务器
	server.Serve()
}
