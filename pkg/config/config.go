package config

import (
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type App struct {
	Driver string //驱动类型
}

var AppConfig = &App{}

type CPUUsage struct {
	Check   uint //CPU验证值
	Normal  uint //CPU正常值
	Limited uint //CPU被限制后的值
}

var CPUUsageConfig = &CPUUsage{}

type CPUScore struct {
	MaxScore int
	MinScore int
}

var CPUScoreConfig = &CPUScore{}

var cfg *ini.File

// Init 初始化配置
func Init(configFile []byte) {
	var err error

	//获取运行目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//判断与创建配置文件
	if _, err := os.Stat(dir + "/config.ini"); os.IsNotExist(err) {
		err := ioutil.WriteFile(dir+"/config.ini", configFile, 0755)
		if err != nil {
			log.Fatalf("无法写入配置文件.")
		}

		log.Printf("已为您创建配置文件，请根据需求修改config.ini再重新运行程序")
		os.Exit(0)
	}

	//读取配置文件
	cfg, err = ini.Load(dir + "/config.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	//映射配置到结构体
	mapTo("usage", CPUUsageConfig)
	mapTo("score", CPUScoreConfig)
	mapTo("app", AppConfig)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
