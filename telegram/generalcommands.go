package telegram

import (
	"OttBot2/metrics"
	"OttBot2/models"
	"OttBot2/settings"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
	"strconv"
)

//Help response for main channel
func MainChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetChannelID(), "/mods - Call the chat moderators\n/ask <word> - Ask the bot about a word")
	bot.Send(newMsg)
}

//Help response for control channel
func ControlChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetControlID(), "/info <@username OR TelegramID> - Get information about a username\n/warn <@username OR TelegramID> <warning message> - Record a warning for a user\n/find <display name> - Find a user by their display name\n/status - Get bot status information\n/count - Get the number of chains in the markov database")
	bot.Send(newMsg)
}

//Handles the /mods command
func SummonMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	if user.PingAllowed {
		newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
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
		newMess := tgbotapi.NewMessage(settings.GetControlID(), "Time since startup: "+metrics.TimeSinceStart().String())
		bot.Send(newMess)
	}
}

//Handles furaffinity links
func GetFARating(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		regex := regexp.MustCompile(`.*furaffinity\.net\/(?:view|full)\/(\d*)`)
		procString := regex.FindStringSubmatch(upd.Message.Text)
		if procString != nil {
			rating := FALookup(procString[1])
			if rating != "" && rating != "General" {
				mainMess := tgbotapi.NewMessage(settings.GetChannelID(), "\u2757 The linked image above is rated NSFW \u2757")
				bot.Send(mainMess)
				fromId := strconv.Itoa(upd.Message.From.ID)
				modMess := tgbotapi.NewMessage(settings.GetControlID(), "User ID "+fromId+" just linked a NSFW image")
				bot.Send(modMess)
				newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
				bot.Send(newFwd)
			}
		}
	}
}

//Handles furry network links
func GetFNRating(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		regex := regexp.MustCompile(`.*furrynetwork\.com/.*\?viewId=(\d*)`)
		procString := regex.FindStringSubmatch(upd.Message.Text)
		if procString != nil {
			rating := FNLookup(procString[1])
			if rating != 0 {
				mainMess := tgbotapi.NewMessage(settings.GetChannelID(), "\u2757 The linked image above is rated NSFW \u2757")
				bot.Send(mainMess)
				fromId := strconv.Itoa(upd.Message.From.ID)
				modMess := tgbotapi.NewMessage(settings.GetControlID(), "User ID "+fromId+" just linked a NSFW image")
				bot.Send(modMess)
				newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
				bot.Send(newFwd)
			}
		}
	}
}

//Handles normal E621 links
func GetE621IDRating(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		regex := regexp.MustCompile(`.*e621\.net/post/show/(\d*)`)
		procString := regex.FindStringSubmatch(upd.Message.Text)
		if procString != nil {
			rating := E621IDLookup(procString[1])
			if rating != "" && rating != "s" {
				mainMess := tgbotapi.NewMessage(settings.GetChannelID(), "\u2757 The linked image above is rated NSFW \u2757")
				bot.Send(mainMess)
				fromId := strconv.Itoa(upd.Message.From.ID)
				modMess := tgbotapi.NewMessage(settings.GetControlID(), "User ID "+fromId+" just linked a NSFW image")
				bot.Send(modMess)
				newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
				bot.Send(newFwd)

			}
		}
	}
}

//Handles direct image E621 links
func GetE621MD5Rating(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		regex := regexp.MustCompile(`.*static1\.e621\.net/data/../../(.+?)\.`)
		procString := regex.FindStringSubmatch(upd.Message.Text)
		if procString != nil {
			rating := E621MD5Lookup(procString[1])
			if rating != "" && rating != "s" {
				mainMess := tgbotapi.NewMessage(settings.GetChannelID(), "\u2757 The linked image above is rated NSFW \u2757")
				bot.Send(mainMess)
				fromId := strconv.Itoa(upd.Message.From.ID)
				modMess := tgbotapi.NewMessage(settings.GetControlID(), "User ID "+fromId+" just linked a NSFW image")
				bot.Send(modMess)
				newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
				bot.Send(newFwd)
			}
		}
	}
}
