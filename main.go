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

	//Any channel commands
	telegram.Register(`^\/ping`, 0, telegram.TestCmd)
	telegram.Register(`^\/mods`, 0, telegram.SummonMods)

	//Main channel commands
	telegram.Register(`^\/help`, settings.GetChannelID(), telegram.MainChannelHelp)
	telegram.Register(`.*furaffinity\.net\/(?:view|full)\/(\d*)`, settings.GetChannelID(), telegram.GetFARating)
	telegram.Register(`.*furrynetwork\.com/.*\?viewId=(\d*)`, settings.GetChannelID(), telegram.GetFNRating)
	telegram.Register(`.*e621\.net/post/show/(\d*)`, settings.GetChannelID(), telegram.GetE621IDRating)
	telegram.Register(`.*static1\.e621\.net/data/../../(.+?)\.`, settings.GetChannelID(), telegram.GetE621MD5Rating)
	telegram.Register(`.*`, settings.GetChannelID(), telegram.HandleUsers)

	//Control channel commands
	telegram.Register(`^\/help`, settings.GetControlID(), telegram.ControlChannelHelp)
	telegram.Register(`^\/info @.+`, settings.GetControlID(), telegram.FindUserByUsername)
	telegram.Register(`^\/info \d+`, settings.GetControlID(), telegram.FindUserByUserID)
	telegram.Register(`^\/warn @.+ .+`, settings.GetControlID(), telegram.WarnUserByUsername)
	telegram.Register(`^\/warn \d+ .+`, settings.GetControlID(), telegram.WarnUserByID)
	telegram.Register(`^\/find .+`, settings.GetControlID(), telegram.LookupAlias)
	telegram.Register(`^\/ban \d+`, settings.GetControlID(), telegram.PreBan)
	telegram.Register(`^\/status`, settings.GetControlID(), telegram.GetBotStatus)

	//Callbacks
	telegram.RegisterCallback(`^\/togglemods \d+`, telegram.ToggleMods)
	telegram.RegisterCallback(`^\/getwarnings \d+`, telegram.DisplayWarnings)
	telegram.RegisterCallback(`^\/info \d+`, telegram.FindUserByUserID)
	telegram.RegisterCallback(`^\/preconfirmban \d+`, telegram.PreConfirmBan)
	telegram.RegisterCallback(`^\/cancelban`, telegram.CancelBan)
	telegram.RegisterCallback(`^\/confirmban \d+`, telegram.ConfirmBan)

	telegram.RegisterNewMember(settings.GetChannelID(), telegram.HandleNewMember)
	telegram.RegisterLeftMember(settings.GetChannelID(), telegram.HandleLeftMember)
	telegram.InitBot(settings.GetBotToken())
}
