package main

import (
    "github.com/AntonioLangiu/pdf_merge_bot/bot"
    "github.com/AntonioLangiu/pdf_merge_bot/common"
)

func main() {
    configuration := common.LoadConfiguration()
    bot.LoadBot(configuration)
}
