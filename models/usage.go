package models

import (
	"log"
	"time"
)

type Usage struct {
	VPSID       uint `gorm:"primaryKey;unique_index"`
	CPUNow      uint64
	CreatedTime time.Time
}

type UsageAvg struct {
	Usage
	CPUAVG uint
}

func (usage *Usage) GetAllRecords() ([]UsageAvg, error) {
	var avgUsage []UsageAvg
	result := DB.Select("AVG(`CPUNow`) as CPUAVG").Group("VPSID").Find(&usage)
	return avgUsage, result.Error
}

func (usage *Usage) AddRecord(vmid uint, cpupercent uint64) {
	DB.Create(&Usage{VPSID: vmid, CPUNow: cpupercent, CreatedTime: time.Now()})
	log.Printf("已创建记录，VMID：%d，CPU：%d%%", vmid, cpupercent)
}

func (usage *Usage) ClearOldRecord() {
	hourAgo, _ := time.ParseDuration("1h")

	DB.Model(usage).Where("CreatedTime < ?", hourAgo).Delete(&Usage{})
}
