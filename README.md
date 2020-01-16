<p align="center">
    <img src="https://neu.ee/img/logo.png" width="200" alt="ipgw"/>
</p>

<h1 align="center">IPGW Tool</h1>
<p align="center">
<img src="https://img.shields.io/github/release-date/iMyOwn/ipgw" alt="">
<img src="https://img.shields.io/github/license/imyown/ipgw" alt="">
<img src="https://img.shields.io/github/go-mod/go-version/iMyOwn/ipgw" alt="">
<img src="https://img.shields.io/github/languages/code-size/iMyOwn/ipgw" alt="">
</p>

> 东北大学目前唯一非官方跨平台校园网关客户端 😛

<p align="center">
    <img src="https://neu.ee/img/banner@v1.3.1.png" alt="banner"/>
</p>



**部分功能仅用以测试网关与一网通，请勿用于违法违纪用途，使用者自行承担责任，后果自负**




# 目录

* [简介](#简介)
* [功能](#功能)
* [下载](#下载)
* [快速使用](#快速使用)
  * [登陆](#登陆)
  * [登出](#登出)
  * [强制下线](#强制下线)
  * [查询](#查询)
  * [工具](#工具)
* [更新](#更新)
* [命令说明](#命令说明)
  * [Login](#login)
  * [Logout](#logout)
  * [Kick](#kick)
  * [List](#list)
  * [Test](#test)
  * [Tool](#Tool)
    * [Get](#Tool-Get)
    * [List](#Tool-List)
    * [Remove](#Tool-Remove)
    * [Update](#Tool-Update)
  * [Update](#update)
  * [Fix](#fix)
  * [Version](#version)
* [常见问题](#常见问题)
* [二次开发](#二次开发)
  * [关于文本](#关于文本)
  * [关于扩展](#关于扩展)
* [开源协议](#开源协议)




# 简介

1. 每次连接校园网之后都要打开网页进行登录，页面的渲染与密码的输入重复枯燥且无聊。

2. 对于校内部分Linux服务器，没有图形化界面，无法访问网页，网关操作十分繁琐。

为了解决这些问题，`ipgw` 诞生了。

不断加入新功能的`ipgw`，如今已经完全覆盖了除【更换套餐】【更改密码】以外的**所有网关操作**。

如果有新的功能建议，欢迎在本仓库[新建Issue](https://github.com/iMyOwn/ipgw/issues/new)



# 功能

- 使用账号密码登陆
- 使用Cookie登陆
- 登陆时伪装设备
- **无参数快速登陆**
- 使用账号密码登出
- 使用Cookie登出
- **无参数快速登出**
- **强制指定设备下线**
- 检查网络与登陆情况
- **自动更新**
- 修复配置文件
- **信息查询**
  - 查看本地信息
  - 查询账号信息
  - 查询已登陆设备
  - 查询当前套餐
  - 查询扣款记录
  - 查询充值记录
  - 查询使用日志
- **支持工具扩展**



# 下载

本工具为x64架构的linux、osx、windows系统提供了预编译程序。

预编译程序可以通过以下方式下载。

>  也可以在[本仓库的Release页面](https://github.com/iMyOwn/ipgw/releases)或[NEU.ee的Release目录](https://neu.ee/release)进行手动下载

## Linux or OSX
```shell script
# linux
wget https://neu.ee/release/v1.3.1/linux/ipgw && chmod +x ipgw && mv ipgw /usr/local/bin && ipgw version

# osx
# 使用terminal
wget https://neu.ee/release/v1.3.1/osx/ipgw && chmod +x ipgw && mv ipgw /usr/local/bin && ipgw version
```

若遇到问题请参阅[常见问题](#常见问题)，或[寻找帮助](https://github.com/iMyOwn/ipgw/issues/new)
## Win
1. 下载 [ipgw.exe](https://neu.ee/release/v1.3.1/win/ipgw.exe)
2. 下载 [配置脚本](https://neu.ee/release/v1.3.1/win/install.bat)
3. 将配置脚本与`ipgw.exe`放置于同一目录下，右键使用**管理员权限**打开配置脚本，会自动配置并弹出`系统属性`设置窗口
4. 点击`环境变量`打开设置窗口，在**系统环境变量**中找到`Path`，选中后点击`编辑`，在弹出的窗口点击`新建`，输入`%ipgw%`并保存，点击`确认`关闭设置窗口
5. 打开`cmd`(可通过win+r并输入cmd打开)，输入`ipgw version`，若无报错，即配置成功
> 配置成功后下载的`ipgw.exe`与`install.bat`可以删除

## Other

其他架构或系统，可以通过下载源代码自主编译的方式使用

```shell
git clone https://github.com/iMyOwn/ipgw.git 
cd ipgw 
go build -ldflags "-w -s -X ipgw/base.Version=v1.3.1" -o ipgw 
```

# 快速使用

利用`ipgw`，能够大大简化对网关的操作.

`ipgw`为命令行工具，osx与linux系统请在`terminal`中使用，windows系统请在`cmd`中使用

## 登陆

在没有保存过账号的情况下登陆

  ```shell script
  ipgw login -u 学号 -p 密码
  ```

可以在登陆时保存账号

  ```shell script
  ipgw login -u 学号 -p 密码 -s
  ```

**在保存了账号后，可以直接登陆**

  ```shell script
  ipgw
  ```

> 默认配置文件保存在用户目录下，名称为`.ipgw`，暂不支持自定义路径，暂不支持保存多个用户

## 登出

在没有保存过账号的情况下登出

  ```shell script
  ipgw logout -u 学号 -p 密码
  ```

在保存了账号后，可以直接登出

  ```shell script
  ipgw logout
  ```

**如果该次登陆使用的是本工具，则无论是否保存账号，都可直接登出**

  ```shell script
  ipgw logout
  ```


## 强制下线

强制指定SID的设备断开校园网

  ```shell script
  ipgw kick SID1 SID2 SID3 ...
  ```

>【注意】该操作可以强制任何人的设备断开校园网，由于滥用该操作造成的一切后果由使用者自负。

强制当前账号在某些设备下线
  ```shell script
  ipgw list -d
  # 会输出类似信息，该命令的具体操作请参阅下文 list 部分
  # No.0 01-01 xx:xx:xx  xxx.xxx.xxx.xxx   xxxxxxxx
  # No.1 01-01 yy:yy:yy  yyy.yyy.yyy.yyy   yyyyyyyy
  # No.2 01-01 zz:zz:zz  zzz.zzz.zzz.zzz   zzzzzzzz

  # 根据最后的八位数字来强制下线
  ipgw kick xxxxxxxx zzzzzzzz
  ```

## 查询

列出本地保存的信息

  ```shell script
  ipgw list -l
  ```

列出**当前登陆账号所有信息**

  ```shell script
  ipgw list -a
  ```

列出**当前登陆账号的已登录设备、账号信息与第一页使用日志**

  ```shell script
  ipgw list -d -i -h 1
  ```

列出**当前登陆账号的已登录设备、账号信息与第一页使用日志**，缩写形式

  ```shell script
  ipgw list -hid
  ```
> 默认获取第一页日志，因此1可以省略

可以通过`-s`、`-u -p`、`-c`查询 已保存账号 / 指定账号 / 指定Cookie 的信息，例如
  - 列出已保存账号的第三页使用日志
    ```shell script
    ipgw list -sh 3
    ```
  - 列出指定账号的的第三页使用日志
    ```shell script
    ipgw list -u 学号 -p 密码 -h 3
    ```
  - 列出指定Cookie对应账号的第三页使用日志
    ```shell script
    ipgw list -c Cookie -h 3
    ```


## 工具

`ipgw`自`v1.3.1`开始支持工具扩展。

列出可用工具

  ```shell script
  ipgw tool list
  ```
列出本地工具

  ```shell script
  ipgw tool list -l
  ```

下载指定工具

  ```shell script
  ipgw tool get [tool name]
  ```

删除指定工具

  ```shell script
  ipgw tool remove [tool name]
  ```

更新指定工具

  ```shell script
  ipgw tool update [tool name]
  ```





# 更新
获取最新版本信息并自动下载更新
```shell script
ipgw update
```

强制自动更新，无论当前是否已是最新版本
```shell script
ipgw update -f
```

【注意】大版本更新后可能出现旧配置文件无法解析，使用`ipgw fix`修复配置文件即可




#  命令说明
### 用法

```
ipgw <command> [arguments]
```
### 命令列表

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

每个命令都已经给出了使用示例，可以使用`ipgw help <command>`查看

以下内容和`ipgw help <command>`的输出相同

## Login

### 用法

```shell script
ipgw login [-u username] [-p password] [-s save] [-c cookie] [-d device] [-i info] [-v view all] 
```
### 参数列表

```
  -u    登陆账号
  -p    登陆密码
  -s    保存该账号
  -c    使用cookie登陆
  -d    使用指定设备信息
  -i    登陆后输出账号信息
  -v    输出所有中间信息
```

### 使用示例

```shell script
  ipgw
    # 效果等同于 ipgw login -i
    # [推荐] 在已经使用-s保存了账号信息的情况下，直接执行ipgw即可完成登陆

  ipgw login -u 学号 -p 密码
    # 使用指定账号登陆网关

  ipgw login -u 学号 -p 密码 -s
    # 本次登陆的账号信息将被保存在用户目录下的.ipgw文件中

  ipgw login
    # 在已经使用-s保存了账号信息的情况下，可以直接使用已经保存的账号登录

  ipgw login -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    # 使用指定cookie登陆

  ipgw login -d win
    # 使用指定设备信息登陆，可选的有win linux osx，默认使用匿名设备信息

  ipgw login -i
    # 登陆成功后输出账号信息，包括账号余额、已使用时长、已使用流量等

  ipgw login [arguments] -v
    # 打印登陆过程中的每一步信息
```



## Logout

### 用法

```shell script
ipgw logout [-u username] [-p password] [-c cookie] [-v view all]
```
### 参数列表

```
  -u    登出账号
  -p    登出密码
  -c    使用cookie登出
  -v    输出所有中间信息
```

### 使用示例

```shell script
  ipgw logout
    # 若本次登陆是通过本工具，则直接登出
    # 若直接登出失败，且有未失效的Cookie，将使用Cookie登出
    # 若Cookie登出失败，且已使用-s保存了账号信息，将使用该账号登出

  ipgw logout -u 学号 -p 密码
    # 使用指定账号登出网关

  ipgw logout -c "ST-XXXXXX-XXXXXXXXXXXXXXXXXXXX-tpass"
    # 使用指定cookie登出

  ipgw logout [arguments] -v
    # 打印登出过程中的每一步详细信息
```

## Kick
### 用法

```shell script
ipgw kick [-v view all] sid1 sid2 sid3 ...
```
### 参数列表

```
  -v    输出所有中间信息
```

### 使用示例

```shell script
  ipgw kick XXXXXXX YYYYYYYY
    # 使指定SID的设备下线

  ipgw kick -v XXXXXXX
    # 使指定SID的设备下线并输出详细的中间信息
```



## List

### 用法

```shell script
ipgw list [-f full] [-v view all] [-s saved] [-u username] [-p password] [-c cookie] [-a all] [-l local info] [-d devices] [-i net info] [-r recharge] [-b bill] [-h history] page
```
### 参数列表

```
  -s    使用保存的账号查询
  -c    使用cookie查询
  -u    使用指定账号查询，需配合 -p
  -p    使用指定账号查询
  -a    列出所有信息
  -l    列出本地保存的账号及网络信息
  -i    列出校园网套餐信息
  -r    列出充值记录
  -d    列出登陆设备
  -b    列出历史账单
  -h    列出校园网使用日志
  -f    输出所有查询结果的详细信息
  -v    输出所有中间信息
```

### 使用示例

```shell script
  ipgw list
    # 效果等同于 ipgw list -l

  ipgw list -l
    # 列出本地保存的账号及会话信息
    # 包括 已保存账号 Cookie CAS

  ipgw list -a
    # 效果等同于 ipgw list -birdh 1
    # 列出当前登陆账号所有信息，必须是使用本工具登陆

  ipgw list -i
    # 查看当前登陆账号的校园网套餐信息
    # 包括 套餐 使用流量 使用时长 余额 使用次数
    # 可使用 -u -p 或 -s 或 -c 查询指定的账号

  ipgw list -r
    # 列出当前登陆账号的充值记录
    # 可使用 -u -p 或 -s 或 -c 查询指定的账号

  ipgw list -d
    # 列出当前登陆账号的已登录设备
    # 可使用 -u -p 或 -s 或 -c 查询指定的账号

  ipgw list -b
    # 列出当前登陆账号的历史付费记录
    # 可使用 -u -p 或 -s 或 -c 查询指定的账号

  ipgw list -h 1
    # 列出当前登陆账号的使用记录的第一页，每页20条
    # 可使用 -u -p 或 -s 或 -c 查询指定的账号

  ipgw list -af
    # 列出所有信息的具体查询结果

  ipgw list -av
    # 列出中间信息
```



## Test

### 用法

```shell script
ipgw test [-v view all]
```
### 参数列表

```
  -v    输出所有中间信息
```

### 使用示例

```shell script
  ipgw test
    # 测试校园网连接与登陆情况

  ipgw test -v
    # 测试校园网连接与登陆情况并输出详细中间信息
```



## Tool

### 用法

```shell script
ipgw tool <command> [arguments]
```



### Tool Get

#### 用法

```shell script
ipgw tool get tool1 tool2 ...
```

#### 参数列表

无

#### 使用示例

```shell script
  ipgw tool get teemo
    # 下载teemo
```


### Tool List

#### 用法

```shell script
ipgw tool list [-a all] [-l local]
```

#### 参数列表

```
  -a    查看所有工具
  -l    查看本地工具
```

#### 使用示例

```shell script
  ipgw tool list
    # 查看可用工具(API兼容)
    
  ipgw tool list -a
    # 查看所有工具
    
  ipgw tool list -l
    # 查看本地已有的工具
```


### Tool Remove

#### 用法

```shell script
ipgw tool remove tool1 tool2 ...
```

#### 参数列表

无

#### 使用示例

```shell script
  ipgw tool remove teemo
    # 从本机删除teemo
```


### Tool Update

#### 用法

```shell script
ipgw tool update [-f force] tool1 tool2 ...
```

#### 参数列表

```
  -f    强制更新
```

#### 使用示例

```shell script
  ipgw tool update teemo
    # 检查更新并升级teemo
    
  ipgw tool update -f teemo
    # 强制更新teemo
```





## Update

### 用法

```shell script
ipgw update [-f force] [-v view all]
```
### 参数列表

```
  -f    强制更新
  -v    输出中间信息与详细报错信息
```

### 使用示例

```shell script
  ipgw update
    # 检查更新并自动更新

  ipgw update -f
    # 强制下载最新版更新
```



## Fix

### 用法

```shell script
ipgw fix
```
### 参数列表

无

### 使用示例

```shell script
  ipgw fix
    # 修复配置文件
```



## Version

### 用法

```shell script
ipgw version [-l list]
```
### 参数列表

```
  -l    输出完整版本功能
```

### 使用示例

```shell script
  ipgw version
    # 查看版本

  ipgw version -l
    # 查看当前版本完整功能
```



# 常见问题

> Permission denied

程序没有执行权限，使用`chmod +x ipgw`赋予可执行权限即可。

<br/>

> ipgw: command not found

这是*nix系统下的报错，没有正确将`ipgw`程序放置于环境变量Path所列出的目录中，推荐将程序放置于`/usr/local/bin`目录下

<br/>

> 'ipgw' is not recognized as an internal or external command,
> operable program or batch file.

这是win系统下的报错，没有正确将`ipgw`程序放置于环境变量Path所列出的目录中，推荐将程序放置于一个不会经常变动的路径下，然后将该路径加入环境变量Path、

<br/>

> wget: command not found

linux用户请使用自己系统对应的包管理工具安装wget，如ubuntu可以使用`sudo apt-get install wget`

mac用户可以使用homebrew安装wget，`brew install wget`；但有可能系统里甚至没有homebrew，建议手动下载然后`chmod +x ipgw`并`mv ipgw /usr/local/bin`

<br/>

> 更新失败或者提示网络失败、获取失败怎么办？

可以加上`-v`重新执行一遍命令，查看输出的具体信息与具体报错，一般来讲具体信息与具体报错已经能给出引发错误的原因。

如果错误确实由程序的bug引发，欢迎[提交Bug](https://github.com/iMyOwn/ipgw/issues/new)

<br/>

# 二次开发

```shell script
# Clone
git clone https://github.com/iMyOwn/ipgw.git
cd ipgw
```

在大部分情况下，`Makefile`无需作修改，修改程序后可以直接构建与加壳

```
# To build
make all VERSION=v1.3.1

# To release
make release VERSION=v1.3.1
```
> 加壳需要预先安装UPX



## 关于文本

基本上所有的输出文本都独立在了各个包中的`text.go`中，方便定制化输出

少部分输出文本在包的`impl.go`中

上下文`Ctx`的输出编写在`ctx`包中



## 关于扩展
To be done.



# 开源协议

MIT license.

