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

	// 分数处理器，每十五分钟运行一次
	_, err = cron2.AddFunc("@every 15m", CPUScoreHandler)
	if err != nil {
		return
	}
	log.Println("已设定CPU积分处理器任务，十五分钟后开始处理积分")

	// 启动任务
	cron2.Start()
}

// CPUDataCollector CPU使用搜集器，搜集CPU数据到数据库
func CPUDataCollector() {

	log.Printf("正在搜集所有运行中的VPS的CPU使用率数据..")
	//列出所有Domain
	domains, _, err := LConnect.ConnectListAllDomains(1, 0)
	if err != nil {
		log.Fatalf("无法获取VPS信息: %v", err)
	}

	var totalcount uint

	for _, domain := range domains {
		totalcount++
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

				cputimepercent = 100 * (cputime2 - cputime1) / (1000000000 * uint64(cpucount))

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
				usage.AddRecord(uint(d.ID), cputimepercent, cpucount)

				//打印数据
				//log.Printf("正在读取：ID:%d, Name:%s, VCPU:%d, VCPU Usage:%d%% \n", d.ID, d.Name, cpucount, cputimepercent)

			}
		}(domain)
	}
	log.Printf("本次CPU使用率数据搜集完毕，共读取了%d个VPS的数据", totalcount)
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
		scoreDataBefore := score.GetScoreData(vpsID)

		log.Printf("正在统计VPSID：%d，CPU平均使用率：%.2f%%", vpsID, cpuAVG)

		if cpuAVG >= float64(config.CPUUsageConfig.Check) && scoreDataBefore.Score > config.CPUScoreConfig.MinScore {
			//过去十五分钟平均CPU使用率超过基准值，扣分
			score.ChangeScore(vpsID, "-")
			log.Printf("VPSID：%d，过去十五分钟平均CPU使用率：%.2f%%，扣分，当前积分：%d", vpsID, cpuAVG, scoreDataBefore.Score-1)
		} else if cpuAVG < float64(config.CPUUsageConfig.Check) && scoreDataBefore.Score < config.CPUScoreConfig.MaxScore {
			//过去十五分钟平均CPU使用率低于基准值，加分
			score.ChangeScore(vpsID, "+")
			log.Printf("VPSID：%d，过去十五分钟平均CPU使用率：%.2f%%，加分，当前积分：%d", vpsID, cpuAVG, scoreDataBefore.Score+1)
		}

		//极值处理
		scoreDataNow := score.GetScoreData(vpsID)
		if scoreDataNow.Score <= 0 && scoreDataBefore.Score > 0 {
			//如果被扣到0分了，调用接口限制其CPU使用率
			ChangeLibVirtKVMCPULimit(vpsID, config.CPUUsageConfig.Limited)
			log.Printf("VPSID：%d，积分已归零，正在限制其CPU限制至：%d%%", vpsID, config.CPUUsageConfig.Limited)
		} else if scoreDataNow.Score > 0 && scoreDataBefore.Score <= 0 {
			//如果分数回到0以上，调用接口恢复其CPU使用率
			ChangeLibVirtKVMCPULimit(vpsID, config.CPUUsageConfig.Normal)
			log.Printf("VPSID：%d，积分已恢复，正在恢复其CPU限制至：%d%%", vpsID, config.CPUUsageConfig.Normal)
		}

	}
}
