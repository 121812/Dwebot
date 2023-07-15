# Dwebot 适用于Dao的微信机器人
Dwebot = Dao + WeChat + Bot， 这是一个使用 Go 语言编写的机器人，适合Dao群体在微信协作中使用


## 功能特性

- 查询NFT地板价
- 微信标记同步discord
- 公众号更新推送群组
- 群聊消息入库

## 安装

1. 下载依赖包：
    ```shell
    go mod download
    ```
2. 配置文件 
    ```go
    // 打开 config.json 配置 Discord_bot_auth，Mysql 连接
    {
    "Discord_bot_auth": ""
    "Mysql_host":"",
    "Mysql_port":"",
    "Mysql_db":"",
    "Mysql_user":"",
    "Mysql_password":""
    }
    ```
3. 配置 discord 频道
    ```go
    // 打开 main.go 找到该行
    discord_send_msg("频道ID", "from wechat:\n"+mark_content[0])
    ```
4. 配置 公众号检测转发群组
    ```go
    // 打开 main.go 找到以下关键字配置

    // 公众号名称
	if name == "公众号名称" {

    // 配置 群组
    msg.Owner().SendTextToGroup(groups.SearchByNickName(1, "群组名称")[0], name+"公众号更新："+"\n"+title.Text()+"\n"+url.Text())
    ```
5. 构建项目：
    ```shell
    go build
    ```
6. 运行项目:
    ```
    chmod +x ./Dwebot
    ./Dwebot
    ```
7. 项目日志：
    ```shell
    logfile.log
    ```

## 项目依赖
- openwechat 框架：https://github.com/eatmoreapple/openwechat
- discordgo 框架：https://github.com/bwmarrin/discordgo
