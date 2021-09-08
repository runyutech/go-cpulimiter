package main

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/robfig/cron/v3"
	"go-cpulimiter/models"
	"go-cpulimiter/pkg/config"
	"log"
	"time"
)

// InitTasks 初始化定时任务执行
func InitTasks() {
	cron2 := cron.New() //创建一个cron实例

	// CPU使用搜集器，每分钟运行一次
	_, err := cron2.AddFunc("@every 1m", CPUDataCollector)
	if err != nil {
		return
	}
	log.Println("已设定CPU使用搜集器任务，一分钟后开始统计")

	// 分数处理器，每五分钟运行一次
	_, err = cron2.AddFunc("@hourly", CPUScoreHandler)
	if err != nil {
		return
	}
	log.Println("已设定CPU积分处理器搜集器任务，一小时后开始处理积分")

	// 启动任务
	cron2.Start()
}

// CPUDataCollector CPU使用搜集器，搜集CPU数据到数据库
func CPUDataCollector() {

	log.Printf("正在列出所有Domain")
	//列出所有Domain
	domains, _, err := LConnect.ConnectListAllDomains(1, 0)
	if err != nil {
		log.Fatalf("无法获取VPS信息: %v", err)
	}

	//

	for _, domain := range domains {
		go func(d libvirt.Domain) {
			//获取Domain基本信息
			var cputimepercent uint64

			status, _, _, _, cputime1, err := LConnect.DomainGetInfo(d)
			if err != nil {
				return
			}
			if status == 1 {
				time.Sleep(1 * time.Second)

				_, _, _, cpucount, cputime2, err := LConnect.DomainGetInfo(d)
				if err != nil {
					return
				}

				cputimepercent = 100 * (cputime2 - cputime1) / 1000000000

				//获取调度器参数列表
				//parameters, _ := LConnect.DomainGetSchedulerParameters(d, 5)

				//根据调度器获取CPU限制
				//var cpuquota int64
				//for _, value := range parameters {
				//	if value.Field == "vcpu_quota" {
				//		cpuquota = value.Value.I.(int64) / 10000
				//	}
				//}

				//写入数据库
				usage := models.Usage{}
				usage.AddRecord(uint(d.ID), cputimepercent)

				//打印数据
				log.Printf("正在读取：ID:%d, Name:%s, VCPU:%d, VCPU Usage:%d%% \n", d.ID, d.Name, cpucount, cputimepercent)

			}
		}(domain)
	}
}

// CPUScoreHandler CPU积分处理器，用来操作CPU分数变化时对VM的操作
func CPUScoreHandler() {
	log.Println("正在运行CPU积分处理器.")

	usage := models.Usage{}
	score := models.Score{}
	//清理太老的数据
	usage.ClearOldRecord()

	//获取全部积分
	usageData, _ := usage.GetAllRecords()
	for _, v := range usageData {
		vpsID := v.VPSID
		cpuAVG := v.CPUAVG
		scoreDataNow := score.GetScoreData(vpsID)

		log.Printf("正在统计VPSID：%d，CPU平均使用率：%d", vpsID, cpuAVG)

		if cpuAVG >= config.CPUUsageConfig.Check && scoreDataNow.Score > config.CPUScoreConfig.MinScore {
			//过去一小时平均CPU使用率超过基准值，扣分
			score.ChangeScore(vpsID, "-")
			log.Printf("VPSID：%d，过去一小时平均CPU使用率：%d，扣分，当前积分：%d", vpsID, cpuAVG, scoreDataNow.Score)
		} else if cpuAVG < config.CPUUsageConfig.Check && scoreDataNow.Score < config.CPUScoreConfig.MaxScore {
			//过去一小时平均CPU使用率低于基准值，加分
			score.ChangeScore(vpsID, "+")
			log.Printf("VPSID：%d，过去一小时平均CPU使用率：%d，加分，当前积分：%d", vpsID, cpuAVG, scoreDataNow.Score)
		}

		//todo 极值处理

	}
}