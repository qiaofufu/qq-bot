package main

import (
	"QQ-BOT/bot"
	_ "QQ-BOT/bot/db/mongodb"
)

func main() {
	b := bot.Default()
	panic(bot.Run(b))
}
