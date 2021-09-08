package main

import (
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"log"
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

func DisconnectLibVirtKVM() {
	err := LConnect.Disconnect()
	if err != nil {
		log.Fatalf("无法断开和本机Libvirt的连接: %v", err)
	}
}

func ChangeLibVirtKVMCPULimit() {

}
