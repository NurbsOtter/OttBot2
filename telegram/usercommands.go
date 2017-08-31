package telegram

import (
	"OttBot2/models"
	"OttBot2/settings"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
)

var BotTarget *models.ChatUser

func HandleUsers(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	foundUser := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	models.UpdateAliases(upd.Message.From.FirstName, upd.Message.From.LastName, foundUser.ID)
	fmt.Println(foundUser)
}
func FindUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		username := upd.Message.Text[7:]
		username = strings.Trim(username, " ")
		username = strings.ToLower(username) //Woo string processing
		user := models.SearchUserByUsername(username)
		if user != nil {
			curAlias := models.GetLatestAliasFromUserID(user.ID)
			fmt.Println(user)
			outmsg := fmt.Sprintf("UserID: %d\nCurrent Name: %s\n", user.TgID, curAlias.Name)
			for _, warn := range models.GetUsersWarnings(user) {
				outmsg += warn.WarningText + "\n"
			}
			newMsg := tgbotapi.NewMessage(settings.GetControlID(), outmsg)
			bot.Send(newMsg)
		} else {
			newMsg := tgbotapi.NewMessage(settings.GetControlID(), "User not found.")
			bot.Send(newMsg)
		}
	}

}
func FindUserByUserID(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userName string
	if upd.Message == nil {
		userName = upd.CallbackQuery.Data[5:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		userName = upd.Message.Text[5:]
	}
	fmt.Println(userName)
	userName = strings.Trim(userName, " ")
	usrID, err := strconv.ParseInt(userName, 10, 64)
	if err != nil {
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Error parsing userID. Make sure it's an actual number!")
		fmt.Println(userName)
		fmt.Println(err)
		bot.Send(newMsg)
		return
	}
	user := models.ChatUserFromTGIDNoUpd(int(usrID))
	if user != nil {
		curAlias := models.GetLatestAliasFromUserID(user.ID)
		fmt.Println(user)
		outmsg := fmt.Sprintf("UserID: %d\nCurrent Name:%s\n", user.TgID, curAlias.Name)
		for _, warn := range models.GetUsersWarnings(user) {
			outmsg += warn.WarningText + "\n"
		}
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), outmsg)
		bot.Send(newMsg)
	} else {
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), "User not found.")
		bot.Send(newMsg)
	}
}
func MakeAliasInlineKeyboard(aliases []models.ChatUser) tgbotapi.InlineKeyboardMarkup {
	var aliasButtons []tgbotapi.InlineKeyboardButton
	for _, alias := range aliases {
		latestID := models.GetLatestAliasFromUserID(alias.ID)
		btnCmd := fmt.Sprintf("/info %d", alias.TgID)
		newButt := tgbotapi.NewInlineKeyboardButtonData(latestID.Name, btnCmd)
		aliasButtons = append(aliasButtons, newButt)
	}
	return tgbotapi.NewInlineKeyboardMarkup(aliasButtons)
}
func SummonMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	user := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	if user.PingAllowed {
		newFwd := tgbotapi.NewForward(settings.GetControlID(), upd.Message.Chat.ID, upd.Message.MessageID)
		bot.Send(newFwd)
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Summoning mods!")
		bot.Send(newMsg)
	} else {
		newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Sorry you were banned from /mods")
		bot.Send(newMsg)
	}
}
func WarnUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[7:] //Remove the /warn @
		procString = strings.TrimLeft(procString, " ")
		userName := strings.Split(procString, " ")[0]
		userName = strings.ToLower(userName)
		message := procString[len(userName)+1:]
		models.AddWarningToUsername(userName, message)
		newMess := tgbotapi.NewMessage(settings.GetControlID(), "Warned "+userName)
		bot.Send(newMess)
	}
}
func WarnUserByID(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[6:] //Remove the /warn
		procString = strings.TrimLeft(procString, " ")
		userID, err := strconv.ParseInt(strings.Split(procString, " ")[0], 10, 64)
		if err != nil {
			fmt.Println("Failed to parse a tgid from /warn")
			return
		}
		chatUser := models.ChatUserFromTGIDNoUpd(int(userID))
		var outMsg string
		if chatUser == nil {
			outMsg = fmt.Sprintf("Could not find TGID %d", userID)
		} else {
			models.AddWarningToID(chatUser.ID, procString[len(strings.Split(procString, " ")[0])+1:])
			outMsg = fmt.Sprintf("Warned %d", userID)
		}
		newMess := tgbotapi.NewMessage(settings.GetControlID(), outMsg)
		bot.Send(newMess)
	}
}
func LookupAlias(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[6:]
		procString = strings.Trim(procString, " ")
		foundAliases := models.LookupAlias(strings.ToLower(procString))
		if len(foundAliases) == 0 {
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), "No Aliases found")
			bot.Send(outMsg)
		} else {
			//outString := "Search Results (Capped at 20!):\n"
			/*for _, user := range foundAliases {
				latestAlias := models.GetLatestAliasFromUserID(user.ID)
				if latestAlias == nil {
					outString += fmt.Sprintf("UserName: @%s UserID: %d\n", user.UserName, user.ID)
				} else {
					outString += fmt.Sprintf("UserName: @%s UserID: %d Latest Alias: %s\n", user.UserName, user.ID, latestAlias.Name)
				}

			}
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), outString)*/
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), "Found Users:")
			outMsg.ReplyMarkup = MakeAliasInlineKeyboard(foundAliases)
			bot.Send(outMsg)
		}
	}
}
func SetBanTarget(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[5:]
		procString = strings.Trim(procString, " ")
		userID, err := strconv.ParseInt(strings.Split(procString, " ")[0], 10, 64)
		if err != nil {
			fmt.Println("Failed to parse a tgid from /warn")
			return
		}
		foundUser := models.ChatUserFromTGIDNoUpd(int(userID))
		if foundUser != nil {
			BotTarget = foundUser
			alias := models.GetLatestAliasFromUserID(foundUser.ID)
			var outMsg string
			if alias != nil {
				outMsg = fmt.Sprintf("Banning user %s\n/yes to confirm", alias.Name)
			} else {
				outMsg = fmt.Sprintf("Banning user %s\n/yes to confirm", foundUser.ID)
			}
			msg := tgbotapi.NewMessage(settings.GetControlID(), outMsg)
			bot.Send(msg)
			fmt.Println(BotTarget)
		} else {
			msg := tgbotapi.NewMessage(settings.GetControlID(), "User not found!")
			bot.Send(msg)
		}
	}
}
func ApplyBannination(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() && BotTarget != nil {
		banConfig := tgbotapi.ChatMemberConfig{}
		banConfig.ChatID = settings.GetChannelID()
		banConfig.UserID = int(BotTarget.TgID)
		BotTarget = nil
		bot.KickChatMember(banConfig)
	}
}
func ClearBotTarget(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() && BotTarget != nil {
		BotTarget = nil
		outMsg := tgbotapi.NewMessage(settings.GetControlID(), "Cleared the ban target.\n\n :(")
		bot.Send(outMsg)
	}
}
