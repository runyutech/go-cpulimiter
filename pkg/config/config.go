package config

import (
	"github.com/go-ini/ini"
	"log"
)

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
func Init() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("usage", CPUUsageConfig)
	mapTo("score", CPUScoreConfig)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
