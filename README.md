## Github leaked patrol

Github leaked patrol为一款github泄露巡航工具：

1. 提供了WEB管理端，后台数据库支持SQLITE3、MYSQL和POSTGRES
1. 双引擎搜索，github code接口搜索全局github以及本地搜索例行监控的repos
1. 支持规则管理（github搜索规则及本地repos搜索规则）
1. 支持github token管理和用户管理
1. 支持在WEB中对扫描结果审核

## 使用方法
- 命令行参数如下：
    1. web指令表示启动web管理端
    1. scan指令表示只启动github搜索
    1. scan -m local，表示只启动本地代码搜索功能
    1. scan -m all，表示同时启动github代码搜索与本地Repos搜索功能

![](http://docs.xsec.io/images/github/paper/usage.png)

- 配置好conf/app.ini中的参数后使用WEB参数后启动WEB服务器。默认会监听到本地的8000端口，默认的管理员账户和密码分别为：`xsec`和`x@xsec.io`。
![](http://docs.xsec.io/images/github/web.png)

- 登录WEB管理端，录入github token、规则。
![](http://docs.xsec.io/images/github/rules.png)

- 启动搜索功能：

![](http://docs.xsec.io/images/github/search.png)

- 审核结果
    1. github code搜索结果审核：
![](http://docs.xsec.io/images/github/report1.png)
    1. 本地repos详细搜索结果审核：
![](http://docs.xsec.io/images/github/report2.png)

## 更新记录
### 2018/5/25，修改了本地扫描的逻辑
1. 只扫描在后台添加的仓库了
1. 同时支持git远程地址和本地地址，格式如下：
![](http://docs.xsec.io/images/github/local_check.jpg)

