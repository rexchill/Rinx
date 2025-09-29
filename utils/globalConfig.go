package utils

import (
	"Rinx/riface"
	"encoding/json"
	"os"
)

// 全局配置，从配置文件中读取
type GlobalConfig struct {
	// 服务器名字
	Name string
	// 服务器的全局对象
	TcpServer riface.IServer
	// 服务器主机的IP
	Host string
	// 服务器主机监听的端口号
	Port int
	// 服务器的版本号
	Version string
	// 读写数据包的最大值
	MaxPacketSize uint32
	// 服务器主机允许的最大客户端连接个数
	MaxConnection int
}

// 定义全局对象，用于传递配置参数
var Config *GlobalConfig

func (g *GlobalConfig) Reload() {
	// 从配置文件读取
	file, err := os.ReadFile("conf/rinx.json")
	if err != nil {
		panic(err)
	}
	// 将配置文件数据解析到 全局配置对象中
	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}

/*
1、init 函数是一个特殊的包初始化函数，用于在程序启动时、main 函数执行前，对包进行初始化设置（如变量赋值、注册组件、打开连接、解析配置等）
2、Go 运行时会在程序启动时自动执行。在所有全局变量初始化之后，main 函数执行之前
3、按文件名的字典序依次执行，同一文件内的 init 按书写顺序执行
4、所有导入包的全局变量初始化 -> 所有导入包的 init() 函数（递归深度优先） → 当前包的全局变量 → 当前包的 init() → main()
*/
func init() {
	// 默认配置
	Config = &GlobalConfig{
		Name:          "RinxServer",
		Host:          "0.0.0.0",
		Port:          19991,
		Version:       "V0.4",
		MaxPacketSize: 4096,
		MaxConnection: 1000,
	}

	// 加载配置文件
	Config.Reload()
}
