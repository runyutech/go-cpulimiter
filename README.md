# 虚拟化CPU积分制自动限制器

这个程序使用积分制自动限制本机运行的KVM虚拟机的CPU使用率，当超过基准比例时自动扣分，低于基准时加分，当扣分超过设定值则强制限制CPU性能上限到一个值，目前支持LXD和Libvirt两种后端驱动。


## 使用方法
首次运行会自动创建一个config.ini配置文件在程序运行的目录，根据需要调整后重新运行即可。

```ini
[app]
#驱动程序: lxd / libvirt-kvm
Driver = "libvirt-kvm"

[usage]
# CPU基准值，平均超过多少CPU使用百分比扣分
Check = 20
# 正常情况的CPU性能百分比
Normal = 100
# 负分情况的CPU性能百分比
Limited = 20

[score]
# 每台VM初始化分数以及最高分
MaxScore = 30
# VM积分分数下限
MinScore = -168
```
