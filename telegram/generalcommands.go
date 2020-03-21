package telegram

import (
	"OttBot2/metrics"
	"OttBot2/models"
	"OttBot2/settings"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

//Help response for main channel
func MainChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetChannelID(), "/mods - Call the chat moderators\n/ping - Check to see if the bot is online")
	bot.Send(newMsg)
	LogCommand(upd, nil)
}

//Help response for control channel
func ControlChannelHelp(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	newMsg := tgbotapi.NewMessage(settings.GetControlID(), "/info <@username OR TelegramID> - Get information about a username\n/warn <@username OR TelegramID> <warning message> - Record a warning for a user\n/find <display name> - Find a user by their display name\n/status - See how long the bot has been up\n/ping - Check to see if the bot is online")
	bot.Send(newMsg)
	LogCommand(upd, nil)
}

//Handles the /mods command
func SummonMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user, err := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	if err != nil {
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "User attempting to summon mods is not a known user. Aborting!")
		bot.Send(newMsg)
		LogCommand(upd, err)
		return
	}

	if user.PingAllowed {
		// Prefixes an "@" symbol if a username is available, and skips it if not
		var alertSenderPrefix string
		if upd.Message.From.UserName != "" {
			alertSenderPrefix = "@"
		} else {
			alertSenderPrefix = ""
		}
		// Informs the main channel that moderators have been summoned
		newAlertMainMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Summoning mods!")
		sentAlertMainMsg, _ := bot.Send(newAlertMainMsg)

		// Informs the control channel that moderators have been summoned
		modsCmdMsgSlice := strings.SplitN(upd.Message.Text, " ", 2)
		var modsCmdMsgText string
		if len(modsCmdMsgSlice) == 2 {
			modsCmdMsgText = modsCmdMsgSlice[1]
		} else {
			modsCmdMsgText = ""
		}
		newAlertControlMsg := tgbotapi.NewMessage(settings.GetAnnounceChannel(), fmt.Sprintf("%s%s is requesting moderator assistance!\n%s", alertSenderPrefix, upd.Message.From.String(), modsCmdMsgText))

		// Adds buttons to view the alert message in the main channel and also to resolve the alert (deletes all messages associated with the alert)
		viewAlertButt := tgbotapi.NewInlineKeyboardButtonURL("View Alert", fmt.Sprintf("https://t.me/%s/%d", upd.Message.Chat.UserName, upd.Message.MessageID))
		resolveAlertButt := tgbotapi.NewInlineKeyboardButtonData("Resolve Alert", fmt.Sprintf("/resolvealert %d %d", upd.Message.MessageID, sentAlertMainMsg.MessageID))
		newAlertControlMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(viewAlertButt, resolveAlertButt))
		bot.Send(newAlertControlMsg)
	} else {
		// User is not allowed to use /mods command
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Sorry, you are banned from /mods")
		bot.Send(newMsg)
	}
	LogCommand(upd, err)
}

// ResolveAlert handler to delete messages which were created during the /mods command
func ResolveAlert(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// Check to see if this is an appropriate CallbackQuery

	if upd.Message == nil {
		modsCommandMsgID, err := strconv.Atoi(strings.Fields(upd.CallbackQuery.Data)[1])
		if err != nil {
			AnswerCallback(upd, bot, err.Error())
		}

		alertMainMsgID, err := strconv.Atoi(strings.Fields(upd.CallbackQuery.Data)[2])
		if err != nil {
			AnswerCallback(upd, bot, err.Error())
		}

		alertControlMsgID := upd.CallbackQuery.Message.MessageID

		// Deletes the /mods command if it is in the main channel
		deleteModsCommandMsg := tgbotapi.NewDeleteMessage(settings.GetChannelID(), modsCommandMsgID)
		_, err = bot.DeleteMessage(deleteModsCommandMsg)
		if err != nil {
			// No callback here because PM's would trigger this
			AnswerCallback(upd, bot, "")
		}

		// Deletes the alert message if it is in the main channel
		deleteAlertMainMsg := tgbotapi.NewDeleteMessage(settings.GetChannelID(), alertMainMsgID)
		_, err = bot.DeleteMessage(deleteAlertMainMsg)
		if err != nil {
			// No callback here because PM's would trigger this
			AnswerCallback(upd, bot, "")
		}

		// Deletes the alert message in the announce channel
		deleteControlAlertMsg := tgbotapi.NewDeleteMessage(settings.GetAnnounceChannel(), alertControlMsgID)
		_, err = bot.DeleteMessage(deleteControlAlertMsg)
		if err != nil {
			AnswerCallback(upd, bot, err.Error())
		}
	}
	AnswerCallback(upd, bot, "")
}

//Returns uptime of the bot
func GetBotStatus(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), fmt.Sprintf("Time since startup: %s", metrics.TimeSinceStart().Round(time.Second).String()))
		bot.Send(newMsg)
	}
	LogCommand(upd, nil)
}
