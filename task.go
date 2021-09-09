package main

import (
	"github.com/robfig/cron/v3"
	"go-cpulimiter/models"
	"go-cpulimiter/pkg/config"
	"go-cpulimiter/pkg/drivers"
	"log"
)

// InitTasks 初始化定时任务执行
func InitTasks() {
	cron2 := cron.New() //创建一个cron实例

	// CPU使用搜集器，每分钟运行一次

	_, err := cron2.AddFunc("@every 1m", drivers.Driver{}.CPUDataCollector)
	if err != nil {
		return
	}
	log.Println("已设定CPU使用搜集器任务，一分钟后开始统计")

	// 分数处理器，每十五分钟运行一次
	_, err = cron2.AddFunc("@every 15m", CPUScoreHandler)
	if err != nil {
		return
	}
	log.Println("已设定CPU积分处理器任务，十五分钟后开始处理积分")

	// 启动任务
	cron2.Start()
}

// CPUScoreHandler CPU积分处理器，用来操作CPU分数变化时对VM的操作
func CPUScoreHandler() {
	log.Println("正在运行CPU积分处理器.")

	usage := models.Usage{}
	score := models.Score{}
	driver := drivers.Driver{}
	//清理太老的数据
	usage.ClearOldRecord()

	//获取全部积分
	usageData, _ := usage.GetAllRecords()
	for _, v := range usageData {
		vpsName := v.VPSName
		cpuAVG := v.CPUAVG
		scoreDataBefore := score.GetScoreData(vpsName)

		log.Printf("正在统计VPS：%s，过去十五分钟CPU平均使用率：%.2f%%", vpsName, cpuAVG)

		if cpuAVG >= float64(config.CPUUsageConfig.Check) && scoreDataBefore.Score > config.CPUScoreConfig.MinScore {
			//过去十五分钟平均CPU使用率超过基准值，扣分
			score.ChangeScore(vpsName, "-")
			log.Printf("VPS：%s，扣分，当前积分：%d", vpsName, scoreDataBefore.Score-1)
		} else if cpuAVG < float64(config.CPUUsageConfig.Check) && scoreDataBefore.Score < config.CPUScoreConfig.MaxScore {
			//过去十五分钟平均CPU使用率低于基准值，加分
			score.ChangeScore(vpsName, "+")
			log.Printf("VPSID：%s，加分，当前积分：%d", vpsName, scoreDataBefore.Score+1)
		}

		//极值处理
		scoreDataNow := score.GetScoreData(vpsName)
		if scoreDataNow.Score <= 0 && scoreDataBefore.Score > 0 {
			//如果被扣到0分了，调用接口限制其CPU使用率
			driver.ChangeLimit(vpsName, config.CPUUsageConfig.Limited)
			log.Printf("VPSID：%s，积分已归零，正在限制其CPU限制至：%d%%", vpsName, config.CPUUsageConfig.Limited)
		} else if scoreDataNow.Score > 0 && scoreDataBefore.Score <= 0 {
			//如果分数回到0以上，调用接口恢复其CPU使用率
			driver.ChangeLimit(vpsName, config.CPUUsageConfig.Normal)
			log.Printf("VPSID：%s，积分已恢复，正在恢复其CPU限制至：%d%%", vpsName, config.CPUUsageConfig.Normal)
		}

	}
}
