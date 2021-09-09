package models

import (
	"time"
)

type Usage struct {
	VPSName     string
	CPUNow      uint64
	CPUCount    uint16
	CreatedTime time.Time
}

type UsageAvg struct {
	Usage
	CPUAVG float64
}

func (usage *Usage) GetAllRecords() ([]UsageAvg, error) {
	var avgUsage []UsageAvg
	result := DB.Model(usage).Select("*, AVG(`cpu_now`) as `CPUAVG`").Group("vps_name").Find(&avgUsage)
	return avgUsage, result.Error
}

func (usage *Usage) AddRecord(vpsName string, cpupercent uint64, cpucount uint16) {
	DB.Create(&Usage{VPSName: vpsName, CPUNow: cpupercent, CreatedTime: time.Now(), CPUCount: cpucount})
	//log.Printf("已创建记录，VPSNAME：%s，CPU：%d%%", vpsName, cpupercent)
}

func (usage *Usage) ClearOldRecord() {
	timeNow := time.Now()
	hourAgo := timeNow.Add(-15 * time.Minute)

	DB.Model(usage).Where("created_time < ?", hourAgo).Delete(&Usage{})
}
