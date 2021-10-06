package drivers

import "go-cpulimiter/pkg/config"

type Driver struct{}

func (d *Driver) Connect() {
	driverConfig := config.AppConfig.Driver
	if driverConfig == "libvirt-kvm" {
		drivernow := LibvirtKVMDriver{}
		drivernow.ConnectLibvirtKVM()
	} else if driverConfig == "lxd" {
		drivernow := LXDDriver{}
		drivernow.ConnectLXD()
	}
}

func (d *Driver) Disconnect() {
	driverConfig := config.AppConfig.Driver
	if driverConfig == "libvirt-kvm" {
		drivernow := LibvirtKVMDriver{}
		drivernow.DisconnectLibvirtKVM()
	} else if driverConfig == "lxd" {
		drivernow := LXDDriver{}
		drivernow.DisconnectLXD()
	}
}

func (d *Driver) CPUDataCollector() {
	driverConfig := config.AppConfig.Driver
	if driverConfig == "libvirt-kvm" {
		drivernow := LibvirtKVMDriver{}
		drivernow.CollectLibvierKVMCPUData()
	} else if driverConfig == "lxd" {
		drivernow := LXDDriver{}
		drivernow.CollectLXDCPUData()
	}
}

func (d *Driver) ChangeLimit(vpsName string, percent uint, cpuCount int) {
	driverConfig := config.AppConfig.Driver
	if driverConfig == "libvirt-kvm" {
		drivernow := LibvirtKVMDriver{}
		drivernow.ChangeLibvierKVMLimit(vpsName, percent, cpuCount)
	} else if driverConfig == "lxd" {
		drivernow := LXDDriver{}
		drivernow.ChangeLXDLimit(vpsName, percent)
	}
}
