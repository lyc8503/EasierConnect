# EasierConnect

> 🚫 **[Disclaimer]**
> 本程序**按原样提供**, 作者**不对程序的正确性或可靠性提供保证**, 请使用者自行判断具体场景是否适合使用该程序, **使用该程序造成的问题或后果由使用者自行承担**.
> 
> 本程序为 EasyConnect 客户端的开源实现, 旨在帮助高校学生在校外访问校内资源, 没有绕过相关流量审计或安全审查的功能. EasyConnect 的一切权利属深信服所有, 若相关人员对该程序有异议, 请邮箱联系我. (admin@lyc8503.site)

---

**===仍在施工中  Still working in progress.===**

目前正常情况下运行基本稳定, 但速度很慢且错误处理不够鲁棒, 待优化.

Release 中的 Windows 预编译二进制在 VirusTotal 上被少数杀软报毒, Edge / Windows defender 也可能误报, 具体原因不明.

所有二进制文件均使用 Github Actions 编译并自动上传, 相关配置已经公开, 若介意也可自行拉取代码审查后在本地编译.


Complete: 
- [x] Protocol reverse engineering
- [x] Web login reverse engineering
- [x] ~~**now Works on Linux TUN**~~
- [x] **now Works on All platforms(maybe) with socks5**

To-do: 
- [x] Get Session ID & IP
- [x] Socks to L3
- [x] Better code formatting & better logging
- [x] Performace improvement
- [ ] Support 7.6.3 versions
- [x] Better error handling
- [ ] Pack & Release
