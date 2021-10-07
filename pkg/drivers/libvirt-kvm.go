package drivers

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"go-cpulimiter/models"
	"log"
	"os/exec"
	"strconv"
	"time"
)

var LConnect *libvirt.Libvirt

type LibvirtKVMDriver struct{}

// ConnectLibvirtKVM 连接到Libvirt-KVM类型的主机
func (thisDriver *LibvirtKVMDriver) ConnectLibvirtKVM() {

	//连接到libvirt
	l := libvirt.NewWithDialer(dialers.NewLocal())

	err := l.Connect()
	if err != nil {
		log.Fatalf("无法连接到本机Libvirt：%v", err)
	}
	log.Println("已连接到本机LibVirt.")
	LConnect = l
}

// DisconnectLibvirtKVM 断开连接到Libvirt-KVM类型的主机
func (thisDriver *LibvirtKVMDriver) DisconnectLibvirtKVM() {
	err := LConnect.Disconnect()
	if err != nil {
		log.Fatalf("无法断开和本机Libvirt的连接: %v", err)
	}
}

// ChangeLibvierKVMLimit 调整Libvirt-KVM类型的主机的限制
func (thisDriver *LibvirtKVMDriver) ChangeLibvierKVMLimit(vpsName string, percent uint, cpuCount int) {

	quota := strconv.Itoa(int(percent * 1000))
	quotaGlobal := strconv.Itoa(int(percent*1000) * cpuCount)

	cmd := exec.Command("/bin/sh", "-c", "virsh schedinfo --live --set vcpu_period=100000 "+vpsName+" --config && virsh schedinfo --live --set vcpu_quota="+quota+" "+vpsName+" --config && virsh schedinfo --live --set global_period=100000 "+vpsName+" --config && virsh schedinfo --live --set global_quota="+quotaGlobal+" "+vpsName+" --config")

	err := cmd.Run()
	if err != nil {
		log.Fatalf("在调整VPS的CPU限制时出现错误：%s", err)
	}
}

// CollectLibvierKVMCPUData CPU使用搜集器，搜集CPU数据到数据库
func (thisDriver *LibvirtKVMDriver) CollectLibvierKVMCPUData() {

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
				usage.AddRecord(d.Name, cputimepercent, cpucount)

				//打印数据
				//log.Printf("正在读取：ID:%d, Name:%s, VCPU:%d, VCPU Usage:%d%% \n", d.ID, d.Name, cpucount, cputimepercent)
			}
		}(domain)
	}
	log.Printf("本次CPU使用率数据搜集完毕，共读取了%d个VPS的数据", totalcount)
}
