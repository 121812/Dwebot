package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/eatmoreapple/openwechat"
)

func main() {

	// 初始化数据库
	Init_db()

	// 初始化日志文件
	logfile, err := os.Create("logfile.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption())

	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Println(err)
		return
	}
	groups, err := self.Groups()
	fmt.Println(groups, err)

	bot.MessageHandler = func(msg *openwechat.Message) {

		fmt.Println(msg.Content)

		// 公众号监测&推送
		if msg.IsArticle() {
			name := ""
			doc := etree.NewDocument()
			doc.ReadFromString(msg.Content)
			// 推送过滤
			if len(doc.FindElements("//appmsg/mmreader/category/name")) != 0 {
				name = doc.FindElements("//appmsg/mmreader/category/name")[0].Text()
				log.Println(name)
			}
			// 公众号名称
			if name == "" {
				items := doc.FindElements("//appmsg/mmreader/category/item")
				for _, item := range items {
					title := item.FindElement("title")
					url := item.FindElement("url")
					// 配置转发群组
					msg.Owner().SendTextToGroup(groups.SearchByNickName(1, "test")[0], name+"公众号更新："+"\n"+title.Text()+"\n"+url.Text())
					log.Println(name+"\n"+title.Text()+"\n"+url.Text(), "\n", url.Text(), msg.Content)
				}
			}
		}

		// 群聊天记录入库
		if msg.IsSendByGroup() && !msg.IsSystem() && !msg.IsArticle() {
			sender, _ := msg.SenderInGroup()
			sendgr, _ := msg.Sender()
			sender_content := msg.Content
			fmt.Println(sender, err, sender_content, sendgr)

			year, month, day := time.Now().Date()
			hour, min, sec := time.Now().Hour(), time.Now().Minute(), time.Now().Second()

			WechatChatLog := Wechat_chat_log{
				Time:         fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec),
				Send_user:    sender.NickName,
				Send_content: sender_content,
				Send_group:   sendgr.NickName,
			}
			fmt.Println(*sendgr)

			if Insert_wechat_chat_log(WechatChatLog) {
				fmt.Println("入库成功")
			}
		}

		// /bot 使用介绍
		if msg.IsText() && len(msg.Content) >= 4 && msg.Content[:4] == "/bot" {
			bot_msg := "/help 获取帮助、 /f 获取地板价、 /mark 标记信息"
			msg.ReplyText(bot_msg)
		}

		// /help 新人帮助
		if msg.IsText() && len(msg.Content) >= 5 && msg.Content[:5] == "/help" {
			help_msg := `
			新人帮助
1. 对待社区成员要有最基本的尊重
2. 不要发表明显冒犯他人的言论
3. 不要发送垃圾信息
4. 不要分享和提供成人相关内容
5. 不要使用冒犯他人的名称与头像
6. 不要泄露任何关于你个人的信息
7. 不要上传任何盗版或非法的内容
8. 避免提供具误导性的错误信息
`
			msg.ReplyText(help_msg)
		}

		// /f 地板价查询
		if msg.IsText() && len(msg.Content) > 2 && msg.Content[:2] == "/f" {
			format_content := strings.Replace(msg.Content[2:], " ", "", -1)
			floor_price_value := floor_price(format_content)
			floor_price_value_res := fmt.Sprintf("Floor price: %.2f", floor_price_value)
			fmt.Println(floor_price_value_res)
			msg.ReplyText(floor_price_value_res)
		}

		// /mark 标记同步discord
		if msg.IsText() && len(msg.Content) > 5 && strings.Contains(msg.Content, "- - - - - - - - - - - - - - -") && strings.Contains(msg.Content, "/mark") {
			mark_content := strings.Split(msg.Content, "- - - - - - - - - - - - - - -")
			discord_send_msg("123456", "from wechat:\n"+mark_content[0])
			msg.ReplyText("标记成功")
		}
	}

	bot.Block()
}
