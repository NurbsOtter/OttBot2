package telegram

import (
	"OttBot2/models"
	"OttBot2/settings"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
)

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

}
func WarnUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[7:] //Remove the /warn @
		procString = strings.TrimLeft(procString, "cutset")
		userName := strings.Split(procString, " ")[0]
		userName = strings.ToLower(userName)
		message := procString[len(userName)+1:]
		models.AddWarningToUsername(userName, message)
		newMess := tgbotapi.NewMessage(settings.GetControlID(), "Warned "+userName)
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
			outString := "Search Results (Capped at 20!):\n"
			for _, user := range foundAliases {
				outString += fmt.Sprintf("UserName: @%s TelegramID: %d\n", user.UserName, user.TgID)
			}
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), outString)
			bot.Send(outMsg)
		}
	}
}
