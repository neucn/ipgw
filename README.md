<p align="center">
    <img src="https://github.com/neucn/ipgw/raw/master/.doc/logo.png?raw=true" width="200" alt="ipgw"/>
</p>

<h2 align="center">IPGW</h2>
<h3 align="center">东北大学非官方跨平台校园网关客户端</h3>
<p align="center">
<img src="https://img.shields.io/github/v/release/neucn/ipgw" alt="">
<img src="https://img.shields.io/github/issues/neucn/ipgw?color=rgb%2877%20199%20166%29" alt="">
<img src="https://img.shields.io/github/downloads/neucn/ipgw/total?color=ea8f14&label=users" alt="">
<img src="https://img.shields.io/github/license/neucn/ipgw" alt="">
</p>

<p align="center"><a href="#安装">安装</a> | <a href="#快速开始">快速开始</a> | <a href="https://github.com/neucn/ipgw/issues/new">反馈</a></p>

# 安装

## Windows

在 Powershell 中执行以下命令安装

```powershell
iwr https://raw.githubusercontent.com/neucn/ipgw/master/install.ps1 -useb | iex
```

## Linux/FreeBSD/OSX

在 shell 中执行以下命令安装

```shell
curl -fsSL https://raw.githubusercontent.com/neucn/ipgw/master/install.sh | sh
```

## Others

其他系统的同学请 clone 到本地后自行编译

# 快速开始

> 须知：本项目的最初目的仅在于满足作者本人的日常使用，因此工具的输出文本中同时存在中英文。
>
> 欢迎有兴趣的同学提起 [Pull Request](https://github.com/neucn/ipgw/pulls) 将输出文本统一为中文 😀

保存 ipgw 账号 (密码将被加密存储)

```shell
ipgw config account add -u "学号" -p "密码" --default
```

使用默认账号快速登录，需要先保存至少一个 ipgw 账号在本地

```shell
ipgw
```

快速登出

```shell
ipgw logout
```

查看校园网信息，如套餐详情、使用记录、扣费记录等

```shell
ipgw info -a
```

检测校园网连接状况

```shell
ipgw test
```

更新工具

```shell
ipgw update
```

更多命令及其可配置项请使用 `ipgw help` 与 `ipgw help [command name]` 查看
