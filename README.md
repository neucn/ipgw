# IPGW Tool
![](https://img.shields.io/github/release-date/iMyOwn/ipgw)
![](https://img.shields.io/github/license/imyown/ipgw)
![](https://img.shields.io/github/go-mod/go-version/iMyOwn/ipgw)
![](https://img.shields.io/github/languages/code-size/iMyOwn/ipgw)

官网正在建设中 [NEU.ee](https://neu.ee)

所有的发布版本请见本仓库Release或 [NEU.ee/release](https://neu.ee/release)

**部分功能仅用以测试网关与一网通，请勿用于违法违纪用途，使用者自行承担责任，后果自负**

## 功能
利用`ipgw`，能够大大简化对网关的操作，以登陆为例:

- 在没有保存过账号的情况下，需要使用`ipgw login -u 学号 -p 密码`登陆

- 使用`ipgw login -u 学号 -p 密码 -s`来保存本次登陆所用的账号密码

- 在保存了账号后，可以使用`ipgw`直接使用账号登陆

- 除了简化基础操作以外，`ipgw login`还提供了其他功能，例如:
  1. 使用Cookie登陆: `ipgw login -c Cookie`
  2. 伪装设备登陆: `ipgw login -d android`

除了登陆，本工具还有以下功能:
```
># ipgw version -l

  Login
    [=]  基础登陆
    [=]  保存账号
    [=]  Cookie登陆
    [=]  伪装设备

  Logout
    [=]  基础登出
    [=]  Cookie登出

  Kick
    [=]  使指定设备下线

  List
    [=]  查看本地信息
    [=]  查询账号信息
    [=]  查询已登陆设备
    [=]  查询当前套餐
    [=]  查询扣款记录
    [=]  查询充值记录
    [=]  查询使用日志

  Test
    [=]  网络检查

  Update
    [=]  获取最新版本更新日志
    [=]  更新到最新版本

  Version
    [=]  查看当前版本信息

  Help
    [=]  查看各命令信息
```

## 下载
当前支持x64架构的linux、osx、windows系统
### Linux or OSX
```shell script
# linux
wget https://neu.ee/release/v1.1.0/linux/ipgw && chmod +x ipgw && mv ipgw /usr/local/bin
# osx
wget https://neu.ee/release/v1.1.0/osx/ipgw && chmod +x ipgw && mv ipgw /usr/local/bin
```
### Win
1. 下载 [https://neu.ee/release/v1.1.0/win/ipgw.exe](https://neu.ee/release/v1.1.0/linux/ipgw)
2. 将`ipgw.exe`放置于加入了Path环境变量的路径下


## 使用
用法:
```
ipgw <command> [arguments]
```
命令:
```
version     版本查询
login       基础登陆
logout      基础登出
list        获取各类信息
kick        使指定设备下线
test        校园网测试
fix         修复配置文件
update      更新版本
```

每个命令都已经给出了使用示例，请使用`ipgw help <command>`查看

如`ipgw help login`

```
用法: ipgw login [-u username] [-p password] [-s save] [-c cookie] [-d device] [-i info] [-v full view]

提供登陆校园网关功能
  -u    登陆账号
  -p    登陆密码
  -s    保存该账号
  -c    使用cookie登陆
  -d    使用指定设备信息
  -i    登陆后输出账号信息
  -v    输出所有中间信息

  ipgw
    效果等同于 ipgw login -i
    [推荐] 在已经使用-s保存了账号信息的情况下，直接执行ipgw即可完成登陆
  ipgw login -u 学号 -p 密码
    使用指定账号登陆网关
  ipgw login -u 学号 -p 密码 -s
    本次登陆的账号信息将被保存在用户目录下的.ipgw文件中
  ipgw login
    在已经使用-s保存了账号信息的情况下，可以直接使用已经保存的账号登录
  ipgw login -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    使用指定cookie登陆
  ipgw login -d android
    使用指定设备信息登陆，可选的有win linux osx，默认使用匿名设备信息
  ipgw login -i
    登陆成功后输出账号信息，包括账号余额、已使用时长、已使用流量等
  ipgw login [arguments] -v
    打印登陆过程中的每一步信息
```

默认配置文件保存在用户目录下，名称为`.ipgw`，暂不支持自定义路径，暂不支持保存多个用户

## 参与开发或定制化

```shell script
# Clone
git clone https://github.com/iMyOwn/ipgw.git
cd ipgw

# To build
make all VERSION=v1.1.0

# To release
make release VERSION=v1.1.0
```
### 关于文本
基本上所有的输出文本都独立在了各个包中的`text.go`中，方便定制化输出

少部分输出文本在各个包的`impl.go`中

上下文`Ctx`的输出编写在`base/ctx/ctx`中

help命令比较特殊，它的文本在项目目录下的`text`包中，方便`main.go`使用

### 关于扩展
添加新功能请新建一个包
1. 使用`ctx.GetCtx()`获取到全局的上下文
2. 使用`ctx.GetClient()`获取到Cookie可复用的全局http客户端
3. 网关的Cookie保存于`Ctx.User.Cookie`中
4. 一网通的Cookie保存于`Ctx.User.CAS`中
5. 若该功能需要定制化flag解析，请模仿`list`包的写法，`Cmd`对象中的`CustomFlags`设为`true`,并自行编写一个解析函数于命令开始时解析
6. 基本的登录函数，通用的参数提取函数都在`share`包下