# EasierConnect

> 🚫 **[Disclaimer]**
> 本程序**按原样提供**, 作者**不对程序的正确性或可靠性提供保证**, 请使用者自行判断具体场景是否适合使用该程序, **使用该程序造成的问题或后果由使用者自行承担**.
> 
> 本程序为 EasyConnect 客户端的开源实现, 旨在帮助高校学生在校外访问校内资源, 没有绕过相关流量审计或安全审查的功能. EasyConnect 的一切权利属深信服所有, 若相关人员对该程序有异议, 请邮箱联系我. (admin@lyc8503.site)

---

### 背景

国内众多高校都在使用深信服(Sangfor)公司的 EasyConnect 作为学校 VPN 的解决方案.

但官方客户端使用上有诸多问题:

- L3 全局代理, 路由规则配置不当, 减慢网络速度, 劣化 NAT 类型

  (比如想要 SSH 连接学校服务器的同时访问 Google 搜索资料, 代理也会经过 VPN, 速度可能大幅下降甚至不可用)
- 与其他代理软件产生冲突, 与其他软件(如游戏反作弊)容易产生冲突
- 安装驱动&后台常驻, 占用资源, 安装该公司自行签发的 CA, 劫持系统 DLL, 降低系统安全性
- 软件设计对系统侵入性强, 卸载后仍有残留, 容易造成未知的问题

现有如 @Hagb 实现的 [docker-easyconnect](https://github.com/Hagb/docker-easyconnect) 方案使用 Docker 容器封印 EasyConnect, 也可以使用虚拟机封印 EasyConnect.

不过我们在一些如手机/笔记本电脑/路由器之类的"边缘设备"上部署 Docker/虚拟机还是开销相对较大, 使用 vps 中转 EasyConnect 的裸 socks 代理也有一定的安全风险和不便(国内 vps 带宽受限).

于是有了这个应该是目前最优雅简洁, 易于使用的 EasyConnect 客户端的替代品.

### 使用方法

本软件为 EasyConnect 的开源实现, 可以直接提供 Socks5 代理供其他应用使用.

在 Release 中下载或自行拉去代码编译对应平台/架构的独立二进制文件在 cli 下执行即可.

目前正常情况下运行基本稳定, 性能对比官方客户端无明显下降, 但仍需在不同使用场景下的更多测试.


Complete: 
- [x] Web login & binary protocol reimplement
- [x] Get Session ID & IP
- [x] Socks to L3
- [x] Support 7.6.x versions
- [x] Better code formatting & better logging
- [x] Performace improvement
- [x] ~~**now Works on Linux TUN**~~
- [x] **now Works on All platforms(maybe) with socks5**
- [x] Pack & Release

To-do: 
- [ ] More tests
