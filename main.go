package main

import (
	"OttBot2/metrics"
	"OttBot2/models"
	"OttBot2/settings"
	"OttBot2/telegram"
	"math/rand"
	"time"
)

func main() {
	settings.LoadSettings()
	models.MakeDB(settings.GetDBAddr())
	metrics.StartUp()
	rand.Seed(time.Now().UnixNano())

	//Any channel commands
	telegram.Register(`^\/ping`, 0, telegram.TestCmd)
	telegram.Register(`^\/mods`, 0, telegram.SummonMods)

	//Main channel commands
	telegram.Register(`^\/help`, settings.GetChannelID(), telegram.MainChannelHelp)
	telegram.Register(`.*`, settings.GetChannelID(), telegram.HandleUsers)
	//Debug command
	telegram.Register(`.*`, 0, telegram.DebugShowID)

	//Control channel commands
	telegram.Register(`^\/help`, settings.GetControlID(), telegram.ControlChannelHelp)
	telegram.Register(`^\/info @.+`, settings.GetControlID(), telegram.FindUserByUsername)
	telegram.Register(`^\/info \d+`, settings.GetControlID(), telegram.FindUserByUserID)
	telegram.Register(`^\/warn @.+ .+`, settings.GetControlID(), telegram.WarnUserByUsername)
	telegram.Register(`^\/warn \d+ .+`, settings.GetControlID(), telegram.WarnUserByID)
	telegram.Register(`^\/find .+`, settings.GetControlID(), telegram.LookupAlias)
	telegram.Register(`^\/status`, settings.GetControlID(), telegram.GetBotStatus)

	//Callbacks
	telegram.RegisterCallback(`^\/togglemods \d+`, telegram.ToggleMods)
	telegram.RegisterCallback(`^\/getwarnings \d+`, telegram.DisplayWarnings)
	telegram.RegisterCallback(`^\/getaliases \d+`, telegram.DisplayAliases)
	telegram.RegisterCallback(`^\/ban \d+`, telegram.PreBan)
	telegram.RegisterCallback(`^\/callbackinfo \d+`, telegram.CallbackInfo)
	telegram.RegisterCallback(`^\/preconfirmban \d+`, telegram.PreConfirmBan)
	telegram.RegisterCallback(`^\/confirmban \d+`, telegram.ConfirmBan)

	telegram.InitBot(settings.GetBotToken())
}
