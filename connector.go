package main

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"log"
	"os/exec"
	"strconv"
)

var LConnect *libvirt.Libvirt

// ConnectLibVirtKVM 连接到Libvirt-KVM类型的主机
func ConnectLibVirtKVM() {

	//连接到libvirt
	l := libvirt.NewWithDialer(dialers.NewLocal())

	err := l.Connect()
	if err != nil {
		log.Fatalf("无法连接到本机Libvirt：%v", err)
	}
	log.Println("已连接到LibVirt.")
	LConnect = l
}

// DisconnectLibVirtKVM 断开连接到Libvirt-KVM类型的主机
func DisconnectLibVirtKVM() {
	err := LConnect.Disconnect()
	if err != nil {
		log.Fatalf("无法断开和本机Libvirt的连接: %v", err)
	}
}

// ChangeLibVirtKVMCPULimit 调整Libvirt-KVM类型的主机的限制
func ChangeLibVirtKVMCPULimit(vmid uint, percent uint) {

	quota := strconv.Itoa(int(percent * 10000))

	cmd := exec.Command("/bin/sh", "-c", "virsh schedinfo --live --set vcpu_quota="+quota+" "+strconv.Itoa(int(vmid))+" --config")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("在调整VPS的CPU限制时出现错误：%s", err)
	}
}
