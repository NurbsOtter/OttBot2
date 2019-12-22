package telegram

import (
	"OttBot2/metrics"
	"OttBot2/models"
	"OttBot2/settings"
	"fmt"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

//Help response for main channel
func MainChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetChannelID(), "/mods - Call the chat moderators\n/ping - Check to see if the bot is online")
	bot.Send(newMsg)
}

//Help response for control channel
func ControlChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetControlID(), "/info <@username OR TelegramID> - Get information about a username\n/warn <@username OR TelegramID> <warning message> - Record a warning for a user\n/find <display name> - Find a user by their display name\n/status - See how long the bot has been up\n/ping - Check to see if the bot is online")
	bot.Send(newMsg)
}

//Handles the /mods command
func SummonMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user, err := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	if err != nil {

	}
	if user.PingAllowed {
		newFwd := tgbotapi.NewForward(settings.GetAnnounceChannel(), upd.Message.Chat.ID, upd.Message.MessageID)
		bot.Send(newFwd)
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Summoning mods!")
		bot.Send(newMsg)
	} else {
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Sorry, you are banned from /mods")
		bot.Send(newMsg)
	}
}

//Returns uptime of the bot
func GetBotStatus(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		newMess := tgbotapi.NewMessage(settings.GetControlID(), fmt.Sprintf("Time since startup: %s", metrics.TimeSinceStart().Round(time.Second).String()))
		bot.Send(newMess)
	}
}
