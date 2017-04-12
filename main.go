package main

import (
    "github.com/AntonioLangiu/pdf_merge_bot/common"
    "github.com/AntonioLangiu/pdf_merge_bot/bot"
)

func main() {
    configuration := common.LoadConfiguration()
    bot.LoadBot(configuration)
}
