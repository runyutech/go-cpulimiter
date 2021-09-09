package main

import (
	_ "embed"
	"github.com/common-nighthawk/go-figure"
	"go-cpulimiter/models"
	"go-cpulimiter/pkg/config"
	"go-cpulimiter/pkg/drivers"
	"log"
)

//go:embed conf/config.ini
var ConfigFile []byte

func init() {
	figure.NewFigure("rainyun.com", "slant", true).Print()
	log.Println("=====================================")
	log.Println("虚拟化平台CPU使用积分系统 v0.1")
	log.Println("现已支持LibVirt-KVM / LXD")
	log.Println("Powered by RainYun LLC")
	log.Println("=====================================")

	// 初始化配置文件
	config.Init(ConfigFile)

	// 链接、构建SQLite数据库
	models.Init()

	// 初始化任务的执行
	InitTasks()
}

func main() {

	// 连接到KVM
	driver := drivers.Driver{}

	driver.Connect()

	// 先搜集一次CPU数据
	driver.CPUDataCollector()

	//退出时销毁Libvirt连接
	defer driver.Disconnect()
	select {}
}
