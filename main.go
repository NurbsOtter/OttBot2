package main

import (
	"OttBot2/models"
	"OttBot2/settings"
	"OttBot2/telegram"
	"OttBot2/webroutes"
	"OttBot2/metrics"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	settings.LoadSettings()
	models.MakeDB(settings.GetDBAddr())
	metrics.StartUp()
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
	telegram.Register("\\/status", settings.GetControlID(), telegram.GetBotStatus)
	telegram.RegisterCallback("^\\/info \\d+", telegram.FindUserByUserID)
	telegram.RegisterNewMember(settings.GetChannelID(), telegram.HandleNewMember)
	telegram.RegisterLeftMember(settings.GetChannelID(), telegram.HandleLeftMember)
	go telegram.InitBot(settings.GetBotToken())
	router := gin.Default()
	router.Static("/", "./frontend")
	adminUser := router.Group("/user")
	{
		adminUser.POST("/register", webroutes.PostRegister)
	}
	router.Run(":3000")
}
