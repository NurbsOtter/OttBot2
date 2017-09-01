package main

import (
	"OttBot2/metrics"
	"OttBot2/models"
	"OttBot2/settings"
	"OttBot2/telegram"
)

func main() {
	settings.LoadSettings()
	models.MakeDB(settings.GetDBAddr())
	metrics.StartUp()
	//Refactor all of this regex to `` notation and remove double escaping
	telegram.Register(`\/ping`, settings.GetChannelID(), telegram.TestCmd)
	telegram.Register(`.*furaffinity\.net\/(?:view|full)\/(\d*)`, settings.GetChannelID(), telegram.GetFARating)
	telegram.Register(`.*furrynetwork\.com/.*\?viewId=(\d*)`, settings.GetChannelID(), telegram.GetFNRating)
	telegram.Register(`.*e621\.net/post/show/(\d*)`, settings.GetChannelID(), telegram.GetE621IDRating)
	telegram.Register(`.*static1\.e621\.net/data/../../(.+?)\.`, settings.GetChannelID(), telegram.GetE621MD5Rating)
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
	telegram.InitBot(settings.GetBotToken())
}
