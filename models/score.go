package models

import (
	"go-cpulimiter/pkg/config"
	"log"
)

type Score struct {
	VPSID uint
	Score int
}

func (score *Score) GetScoreData(vmid uint) *Score {
	var nowData Score
	DB.Where(Score{VPSID: vmid}).Attrs(Score{Score: config.CPUScoreConfig.MaxScore}).FirstOrCreate(&nowData)

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
