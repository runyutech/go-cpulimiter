package models

import (
	"time"
)

type Usage struct {
	VPSID       uint
	CPUNow      uint64
	CreatedTime time.Time
}

type UsageAvg struct {
	Usage
	CPUAVG float64
}

func (usage *Usage) GetAllRecords() ([]UsageAvg, error) {
	var avgUsage []UsageAvg
	result := DB.Model(usage).Select("*, AVG(`cpu_now`) as `CPUAVG`").Group("vps_id").Find(&avgUsage)
	return avgUsage, result.Error
}

func (usage *Usage) AddRecord(vmid uint, cpupercent uint64) {

	DB.Create(&Usage{VPSID: vmid, CPUNow: cpupercent, CreatedTime: time.Now()})
	//log.Printf("已创建记录，VMID：%d，CPU：%d%%", vmid, cpupercent)
}

func (usage *Usage) ClearOldRecord() {
	timeNow := time.Now()
	hourAgo := timeNow.Add(-1 * time.Hour)

	DB.Model(usage).Where("created_time < ?", hourAgo).Delete(&Usage{})
}
