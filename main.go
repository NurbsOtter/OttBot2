package main

import (
	"OttBot2/models"
	"OttBot2/settings"
	"OttBot2/telegram"
	"OttBot2/webroutes"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	settings.LoadSettings()
	models.MakeDB(settings.GetDBAddr())
	telegram.Register("\\/ping", settings.GetChannelID(), telegram.TestCmd)
	telegram.Register(".*", settings.GetChannelID(), telegram.HandleUsers)
	telegram.Register("^\\/info @\\D+", settings.GetControlID(), telegram.FindUserByUsername)
	telegram.Register("^\\/warn @.+", settings.GetControlID(), telegram.WarnUserByUsername)
	telegram.Register("^\\/find .+", settings.GetControlID(), telegram.LookupAlias)
	telegram.Register("^\\/mods", 0, telegram.SummonMods)
	telegram.Register("^\\/warn \\d+", settings.GetControlID(), telegram.WarnUserByID)
	telegram.Register("^\\/ban \\d+", settings.GetControlID(), telegram.SetBanTarget)
	telegram.Register("^\\/yes", settings.GetControlID(), telegram.ApplyBannination)
	telegram.Register("^\\/no", settings.GetControlID(), telegram.ClearBotTarget)
	telegram.Register("^\\/info \\d+", settings.GetControlID(), telegram.FindUserByUserID)
	telegram.RegisterCallback("^\\/info \\d+", telegram.FindUserByUserID)
	go telegram.InitBot(settings.GetBotToken())
	router := gin.Default()
	router.Static("/", "./frontend")
	adminUser := router.Group("/user")
	{
		adminUser.POST("/register", webroutes.PostRegister)
	}
	router.Run(":3000")
}
