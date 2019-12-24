package version

var (
	tipQuery         = "获取最新版本信息中..."
	tipAlreadyLatest = "当前已是最新版本"
	tipLatest        = "最新版本: %s\n"

	errNet = "请检查网络连接"

	failQuery = "获取失败"

	detail = `当前版本进展:
  Login
    [=]  基础登陆
    [=]  保存账号
    [=]  Cookie登陆
    [=]  伪装设备
    [=]  状态持久化
  
  Logout
    [=]  基础登出
    [=]  Cookie登出
  
  List
    [=]  本地信息
    [=]  账号信息
    [=]  登陆设备
    [=]  套餐信息
    [=]  扣款记录
    [=]  充值记录
    [=]  使用日志

  Toggle
    [ ]  多用户模式
    [ ]  基础切换

  Kick
    [=]  指定设备下线

  Test
    [=]  校园网连通性测试
    [ ]  校园网测速

  Version
    [=]  检查版本

    [=]  有功能想法请发送邮件至 i@shangyes.net`
)
