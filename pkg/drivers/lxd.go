package drivers

import (
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"go-cpulimiter/models"
	"log"
	"os/exec"
	"strconv"
	"time"
)

var LXDConnect lxd.InstanceServer

type LXDDriver struct{}

// ConnectLXD 连接到LXD类型的主机
func (thisDriver *LXDDriver) ConnectLXD() {

	c, err := lxd.ConnectLXDUnix("/var/snap/lxd/common/lxd/unix.socket", nil)
	if err != nil {
		log.Fatalf("无法连接到本机LXD：%v", err)
	}
	log.Println("已连接本机LXD.")
	LXDConnect = c
}

// DisconnectLXD 断开连接到LXD类型的主机
func (thisDriver *LXDDriver) DisconnectLXD() {
	LXDConnect.Disconnect()
}

// ChangeLXDLimit 调整LXD类型的主机的限制
func (thisDriver *LXDDriver) ChangeLXDLimit(vpsName string, percent uint) {

	quota := strconv.Itoa(int(percent))

	cmd := exec.Command("/bin/sh", "-c", "lxc config set "+vpsName+" limits.cpu.allowance="+quota+"ms/100ms")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("在调整VPS的CPU限制时出现错误：%s", err)
	}
}

// CollectLXDCPUData CPU使用搜集器，搜集CPU数据到数据库
func (thisDriver *LXDDriver) CollectLXDCPUData() {

	log.Printf("正在搜集所有运行中的VPS的CPU使用率数据..")

	full, err := LXDConnect.GetContainers()
	if err != nil {
		log.Fatalf("无法获取VPS信息: %v", err)
	}

	var totalcount uint

	for _, container := range full {
		go func(ct api.Container) {
			//获取基本信息
			var cputimepercent int64

			if ct.StatusCode == 103 {
				totalcount++

				state, _, err := LXDConnect.GetContainerState(ct.Name)
				if err != nil {
					log.Fatalf("无法获取VPS信息: %v", err)
				}
				cputime1 := state.CPU.Usage
				time.Sleep(1 * time.Second)

				state, _, err = LXDConnect.GetContainerState(ct.Name)
				if err != nil {
					log.Fatalf("无法获取VPS信息: %v", err)
				}
				cputime2 := state.CPU.Usage

				//通过Config获取CPU核心数
				cpucount, _ := strconv.ParseInt(ct.Config["limits.cpu"], 10, 64)
				cputimepercent = 100 * (cputime2 - cputime1) / (1000000000 * cpucount)

				//写入数据库
				usage := models.Usage{}
				usage.AddRecord(ct.Name, uint64(cputimepercent), uint16(cpucount))

				//log.Printf("正在读取：Name:%s, VCPU:%d, VCPU Usage:%d%% \n", ct.Name, cpucount, cputimepercent)

			}
		}(container)
	}
	log.Printf("本次CPU使用率数据搜集完毕，共读取了%d个VPS的数据", totalcount)
}
