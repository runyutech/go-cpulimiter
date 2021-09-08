package models

import (
	"log"
)

type Score struct {
	VPSID uint
	Score int `gorm:"default:1"`
}

func (score *Score) GetScoreData(vmid uint) *Score {
	var nowData Score
	DB.Where(Score{VPSID: vmid}).FirstOrCreate(&nowData)

	return &nowData
}

func (score *Score) ChangeScore(vmid uint, plusor string) {

	scoreData := score.GetScoreData(vmid)

	if plusor == "-" {
		scoreData.Score--
		log.Printf("已更新VMID %d，扣分，目前：%d", vmid, scoreData.Score)
	} else {
		scoreData.Score++
		log.Printf("已更新VMID %d，加分，目前：%d", vmid, scoreData.Score)
	}

	DB.Model(score).Where("vps_id = ?", vmid).Update("score", scoreData.Score)
}
