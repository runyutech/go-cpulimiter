package models

import (
	"log"
)

type Score struct {
	VPSID uint `gorm:"primaryKey;unique_index"`
	Score int
}

func (score *Score) GetScoreData(vmid uint) *Score {
	var nowData Score
	err := DB.First(&nowData, vmid).Error
	if err != nil {
		log.Panicln("无法读取VM的CPU积分数据")
	}

	return &nowData
}

func (score *Score) ChangeScore(vmid uint, plusor string) {

	scoreData := score.GetScoreData(vmid)

	if plusor == "-" {
		scoreData.Score--
		log.Printf("已更新VMID %d，扣分，目前：%d", vmid, score.Score)
	} else {
		scoreData.Score++
		log.Printf("已更新VMID %d，加分，目前：%d", vmid, score.Score)
	}

	DB.Save(&scoreData)
}
