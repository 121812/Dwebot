package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Configuration struct {
	Discord_bot_auth string
}

func discord_send_msg(id string, cotent string) {

	// 打开文件
	config_file, _ := os.Open("config.json")

	// 关闭文件
	defer config_file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(config_file)

	conf := Configuration{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	decoder.Decode(&conf)

	discord, err := discordgo.New("Bot " + conf.Discord_bot_auth)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
		return
	}

	// 建立连接
	err = discord.Open()
	if err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}

	// 发生消息
	discord.ChannelMessageSend(id, cotent)
}
