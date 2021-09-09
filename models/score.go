package models

import (
	"go-cpulimiter/pkg/config"
)

type Score struct {
	VPSName string
	Score   int
}

func (score *Score) GetScoreData(vpsName string) *Score {
	var nowData Score
	DB.Where(Score{VPSName: vpsName}).Attrs(Score{Score: config.CPUScoreConfig.MaxScore}).FirstOrCreate(&nowData)

	return &nowData
}

func (score *Score) ChangeScore(vpsName string, plusor string) {

	scoreData := score.GetScoreData(vpsName)

	if plusor == "-" {
		scoreData.Score--
		//log.Printf("已更新VMID %d，扣分，目前：%d", vmid, scoreData.Score)
	} else {
		scoreData.Score++
		//log.Printf("已更新VMID %d，加分，目前：%d", vmid, scoreData.Score)
	}

	DB.Model(score).Where("vps_name = ?", vpsName).Update("score", scoreData.Score)
}
