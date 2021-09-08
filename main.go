package main

import (
	"go-cpulimiter/models"
	"go-cpulimiter/pkg/config"
)

func init() {

	// 链接、构建SQLite数据库
	models.Init()

	// 初始化配置文件
	config.Init()

	// 初始化任务的执行
	InitTasks()

}

func main() {

	// 连接到KVM
	ConnectLibVirtKVM()

	//退出时销毁Libvirt连接
	defer DisconnectLibVirtKVM()
	select {}
}
