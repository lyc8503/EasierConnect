# EasierConnect
Still working in progress.

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
- [ ] Support 7.6.3 versions
- [ ] Better error handling
- [ ] Better code formatting & better logging
- [ ] Pack & Release
