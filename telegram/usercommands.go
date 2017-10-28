package telegram

import (
	"OttBot2/models"
	"OttBot2/settings"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"strconv"
	"strings"
)

//Handles non-command messages to record user information/changes
func HandleUsers(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	foundUser := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	models.UpdateAliases(upd.Message.From.FirstName, upd.Message.From.LastName, foundUser.ID)
	fmt.Println(foundUser)
}

func HandleMarkov(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		//Don't continue if message is a command or directed at the bot
		if strings.HasPrefix(upd.Message.Text, "/") {
			return
		}
		if strings.HasPrefix(strings.ToLower(upd.Message.Text), "@robosergal_bot") {
			return
		}

		models.LearnMarkov(upd.Message.Text)
		if rand.Intn(int(settings.GetRandomChance())) == 0 {
			fmt.Println("Randomly responding to message '%s'", upd.Message.Text)
			tgbotapi.NewChatAction(settings.GetChannelID(), tgbotapi.ChatTyping)
			response := models.RandomResponse(upd.Message.Text)
			if response == "" {
				return
			}
			newMess := tgbotapi.NewMessage(settings.GetChannelID(), response)
			bot.Send(newMess)
		}
	}
}

//Get user information by telegram ID
func FindUserByUserID(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		userId := upd.Message.Text[5:]
		userId = strings.Trim(userId, " ")
		usrID, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Error parsing userID. Make sure it's an actual number!")
			fmt.Println(userId)
			fmt.Println(err)
			bot.Send(newMsg)
			return
		}
		user := models.ChatUserFromTGIDNoUpd(int(usrID))
		bot.Send(GetUserInfoResponse(user))
	}
}

//Get user information by username
func FindUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		username := upd.Message.Text[7:]
		username = strings.Trim(username, " ")
		username = strings.ToLower(username) //Woo string processing
		user := models.SearchUserByUsername(username)
		if user != nil {
			curAlias := models.GetLatestAliasFromUserID(user.ID)
			fmt.Println(user)
			outmsg := fmt.Sprintf("User ID: %d\nCurrent Name: %s\nActive User: %t\nMod Ping: %t\n", user.TgID, curAlias.Name, user.ActiveUser, user.PingAllowed)
			newMsg := tgbotapi.NewMessage(settings.GetControlID(), outmsg)
			newMsg.ReplyMarkup = MakeUserInfoInlineKeyboard(user.ID)
			bot.Send(newMsg)
		} else {
			newMsg := tgbotapi.NewMessage(settings.GetControlID(), "User not found.")
			bot.Send(newMsg)
			bot.Send(GetUserInfoResponse(user))
		}
	}
}

//Helper method to generate the response object for the info requests
func GetUserInfoResponse(user *models.ChatUser) tgbotapi.MessageConfig {
	if user != nil {
		curAlias := models.GetLatestAliasFromUserID(user.ID)
		fmt.Println(user)

		outMsg := fmt.Sprintf("User ID: %d", user.TgID)
		if user.UserName != "" {
			outMsg += fmt.Sprintf("\nUsername: @%s", user.UserName)
		}
		outMsg += fmt.Sprintf("\nCurrent name: %s\nActive user: %t\nMod ping: %t\n", curAlias.Name, user.ActiveUser, user.PingAllowed)
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), outMsg)
		newMsg.ReplyMarkup = MakeUserInfoInlineKeyboard(user.ID)
		return newMsg
	} else {
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), "User not found.")
		return newMsg
	}
}

//Helper method to generate the buttons for an initial info request
func MakeUserInfoInlineKeyboard(userId int64) tgbotapi.InlineKeyboardMarkup {
	warnCmd := fmt.Sprintf("/getwarnings %d", userId)
	warnButt := tgbotapi.NewInlineKeyboardButtonData("View warnings", warnCmd)
	aliasCmd := fmt.Sprintf("/getaliases %d", userId)
	aliasButt := tgbotapi.NewInlineKeyboardButtonData("View aliases", aliasCmd)
	modsCmd := fmt.Sprintf("/togglemods %d", userId)
	modsButt := tgbotapi.NewInlineKeyboardButtonData("Toggle /mods", modsCmd)
	banCmd := fmt.Sprintf("/ban %d", userId)
	banButt := tgbotapi.NewInlineKeyboardButtonData("Ban user", banCmd)
	keyboardRow1 := tgbotapi.NewInlineKeyboardRow(warnButt, aliasButt)
	keyboardRow2 := tgbotapi.NewInlineKeyboardRow(modsButt, banButt)
	return tgbotapi.NewInlineKeyboardMarkup(keyboardRow1, keyboardRow2)
}

//Helper method to generate the buttons for an info request after view warnings button is pressed
func MakeUserInfoInlineKeyboardRefreshWarnButton(userId int64) tgbotapi.InlineKeyboardMarkup {
	warnCmd := fmt.Sprintf("/getwarnings %d", userId)
	warnButt := tgbotapi.NewInlineKeyboardButtonData("Refresh warnings", warnCmd)
	aliasCmd := fmt.Sprintf("/getaliases %d", userId)
	aliasButt := tgbotapi.NewInlineKeyboardButtonData("View aliases", aliasCmd)
	modsCmd := fmt.Sprintf("/togglemods %d", userId)
	modsButt := tgbotapi.NewInlineKeyboardButtonData("Toggle /mods", modsCmd)
	banCmd := fmt.Sprintf("/ban %d", userId)
	banButt := tgbotapi.NewInlineKeyboardButtonData("Ban user", banCmd)
	keyboardRow1 := tgbotapi.NewInlineKeyboardRow(warnButt, aliasButt)
	keyboardRow2 := tgbotapi.NewInlineKeyboardRow(modsButt, banButt)
	return tgbotapi.NewInlineKeyboardMarkup(keyboardRow1, keyboardRow2)
}

//Helper method to generate the buttons for an info request after view aliases button is pressed
func MakeUserInfoInlineKeyboardRefreshAliasButton(userId int64) tgbotapi.InlineKeyboardMarkup {
	warnCmd := fmt.Sprintf("/getwarnings %d", userId)
	warnButt := tgbotapi.NewInlineKeyboardButtonData("View warnings", warnCmd)
	aliasCmd := fmt.Sprintf("/getaliases %d", userId)
	aliasButt := tgbotapi.NewInlineKeyboardButtonData("Refresh aliases", aliasCmd)
	modsCmd := fmt.Sprintf("/togglemods %d", userId)
	modsButt := tgbotapi.NewInlineKeyboardButtonData("Toggle /mods", modsCmd)
	banCmd := fmt.Sprintf("/ban %d", userId)
	banButt := tgbotapi.NewInlineKeyboardButtonData("Ban user", banCmd)
	keyboardRow1 := tgbotapi.NewInlineKeyboardRow(warnButt, aliasButt)
	keyboardRow2 := tgbotapi.NewInlineKeyboardRow(modsButt, banButt)
	return tgbotapi.NewInlineKeyboardMarkup(keyboardRow1, keyboardRow2)
}

//Callback handler to update a get user info response to add warnings for the user
func DisplayWarnings(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[12:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	chatUser := models.ChatUserFromID(userId)
	infoMessage := GetUserInfoResponse(chatUser)
	outmsg := infoMessage.Text
	warnings := models.GetUsersWarnings(chatUser)
	if len(warnings) > 0 {
		outmsg += "\nWarnings:"
		for _, warn := range warnings {
			outmsg += "\n" + warn.WarningText
		}
	} else {
		outmsg += "\n No warnings found"
	}
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outmsg)
	inlineKeyboard := MakeUserInfoInlineKeyboardRefreshWarnButton(chatUser.ID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

//Callback handler to update a find by alias request after a user button is clicked
func CallbackInfo(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[13:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	chatUser := models.ChatUserFromID(userId)
	infoMessage := GetUserInfoResponse(chatUser)
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, infoMessage.Text)
	inlineKeyboard := MakeUserInfoInlineKeyboard(chatUser.ID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

//Callback handler to update a get user info response to add all known user aliases
func DisplayAliases(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[11:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	chatUser := models.ChatUserFromID(userId)
	infoMessage := GetUserInfoResponse(chatUser)
	outmsg := infoMessage.Text
	aliases := models.GetAliases(chatUser)
	if len(aliases) > 0 {
		outmsg += "\nKnown aliases:"
		for _, alias := range aliases {
			outmsg += "\n" + alias.Name
		}
	} else {
		outmsg += "\n No Aliases found"
	}
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outmsg)
	inlineKeyboard := MakeUserInfoInlineKeyboardRefreshAliasButton(chatUser.ID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

//Warn a user by username
func WarnUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[7:] //Remove the /warn @
		procString = strings.TrimLeft(procString, " ")
		userName := strings.Split(procString, " ")[0]
		userName = strings.ToLower(userName)
		message := procString[len(userName)+1:]
		user := models.SearchUserByUsername(userName)
		if user != nil {
			models.AddWarningToID(user.ID, message)
			newMess := tgbotapi.NewMessage(settings.GetControlID(), "Warned "+userName)
			bot.Send(newMess)
		} else {
			newMess := tgbotapi.NewMessage(settings.GetControlID(), "Could not find user")
			bot.Send(newMess)
		}
	}
}

//Warn a user by telegram ID
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

//Look up users by their alias
func LookupAlias(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[6:]
		procString = strings.Trim(procString, " ")
		foundAliases := models.LookupAlias(strings.ToLower(procString))
		if len(foundAliases) == 0 {
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), "No Aliases found")
			bot.Send(outMsg)
		} else {
			//No idea what kind of cap on responses this will run into in the future, seems to be max 8 buttons
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), "Found Users:")
			if len(foundAliases) > 8 {
				foundAliases = foundAliases[:8]
			}
			outMsg.ReplyMarkup = MakeAliasInlineKeyboard(foundAliases)
			bot.Send(outMsg)
		}
	}
}

//Helper method to generate the buttons for the lookup by alias command
func MakeAliasInlineKeyboard(aliases []models.ChatUser) tgbotapi.InlineKeyboardMarkup {
	var aliasButtons []tgbotapi.InlineKeyboardButton
	for _, alias := range aliases {
		latestID := models.GetLatestAliasFromUserID(alias.ID)
		btnCmd := fmt.Sprintf("/callbackinfo %d", alias.ID)
		newButt := tgbotapi.NewInlineKeyboardButtonData(latestID.Name, btnCmd)
		aliasButtons = append(aliasButtons, newButt)
	}

	if len(aliasButtons) <= 4 {
		return tgbotapi.NewInlineKeyboardMarkup(aliasButtons)
	} else {
		return tgbotapi.NewInlineKeyboardMarkup(aliasButtons[:4], aliasButtons[4:])
	}
}

//Handles the first step in the ban process by displaying the target to a user
func PreBan(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[4:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	foundUser := models.ChatUserFromID(userId)
	infoMessage := GetUserInfoResponse(foundUser)
	outMsg := infoMessage.Text
	outMsg += "\nDo you want to ban this user?"
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg)
	inlineKeyboard := MakeBanInlineKeyboard(foundUser.ID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

//Helper method to generate the buttons for a pre ban request
func MakeBanInlineKeyboard(userId int64) tgbotapi.InlineKeyboardMarkup {
	var banButtons []tgbotapi.InlineKeyboardButton
	confirmCmd := fmt.Sprintf("/preconfirmban %d", userId)
	confirmButt := tgbotapi.NewInlineKeyboardButtonData("Confirm ban", confirmCmd)
	cancelCmd := fmt.Sprintf("/callbackinfo %d", userId)
	cancelButt := tgbotapi.NewInlineKeyboardButtonData("Cancel ban", cancelCmd)
	banButtons = append(banButtons, confirmButt)
	banButtons = append(banButtons, cancelButt)
	return tgbotapi.NewInlineKeyboardMarkup(banButtons)
}

//Handles the callback when a user presses the first confirm ban button
func PreConfirmBan(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[14:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	foundUser := models.ChatUserFromID(userId)
	infoMessage := GetUserInfoResponse(foundUser)
	outMsg := infoMessage.Text
	outMsg += "\nAre you ABSOLUTELY SURE you want to ban this user?"
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg)
	inlineKeyboard := MakeBanConfirmInlineKeyboard(userId)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

//Helper method to generate the buttons for the final ban request
func MakeBanConfirmInlineKeyboard(userId int64) tgbotapi.InlineKeyboardMarkup {
	var banButtons []tgbotapi.InlineKeyboardButton
	confirmCmd := fmt.Sprintf("/confirmban %d", userId)
	confirmButt := tgbotapi.NewInlineKeyboardButtonData("Yes, I am sure", confirmCmd)
	cancelCmd := fmt.Sprintf("/callbackinfo %d", userId)
	cancelButt := tgbotapi.NewInlineKeyboardButtonData("No, cancel ban", cancelCmd)
	banButtons = append(banButtons, confirmButt)
	banButtons = append(banButtons, cancelButt)
	return tgbotapi.NewInlineKeyboardMarkup(banButtons)
}

//Handles the callback when the user presses the final ban confirmation button
func ConfirmBan(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[11:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	foundUser := models.ChatUserFromID(userId)
	banConfig := tgbotapi.ChatMemberConfig{}
	banConfig.ChatID = settings.GetChannelID()
	banConfig.UserID = int(foundUser.TgID)
	bot.KickChatMember(banConfig)
	models.SetActiveUserState(userId, false)
	foundUser.ActiveUser = false
	infoMessage := GetUserInfoResponse(foundUser)
	outMsg := infoMessage.Text
	outMsg += "\nUser banned by "
	if upd.CallbackQuery.From.UserName != "" {
		outMsg += "@" + upd.CallbackQuery.From.UserName
	} else if upd.CallbackQuery.From.FirstName != "" {
		outMsg += upd.CallbackQuery.From.FirstName + " " + upd.CallbackQuery.From.LastName
	}
	outMsg += " TGID: " + strconv.Itoa(upd.CallbackQuery.From.ID)
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg)
	bot.Send(editMsg)
}

//Handles member join events
func HandleNewMember(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		foundUser := models.ChatUserFromTGID(upd.Message.NewChatMember.ID, upd.Message.NewChatMember.UserName)
		models.UpdateAliases(upd.Message.NewChatMember.FirstName, upd.Message.NewChatMember.LastName, foundUser.ID)
	}
}

//Handles member left events
func HandleLeftMember(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		foundUser := models.ChatUserFromTGID(upd.Message.LeftChatMember.ID, upd.Message.LeftChatMember.UserName)
		models.UpdateAliases(upd.Message.LeftChatMember.FirstName, upd.Message.LeftChatMember.LastName, foundUser.ID)
		models.SetActiveUserState(foundUser.ID, false)
	}
}

//Handles toggling of a user's ability to use /mods
func ToggleMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[12:]
		config := tgbotapi.NewCallback(upd.CallbackQuery.ID, "") //We don't need this so get it outta da way.
		bot.AnswerCallbackQuery(config)
	} else {
		return
	}
	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	chatUser := models.ChatUserFromID(userId)
	models.SetModPing(chatUser.ID, !chatUser.PingAllowed)
	chatUser.PingAllowed = !chatUser.PingAllowed
	outMsg := GetUserInfoResponse(chatUser)
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg.Text)
	inlineKeyboard := MakeUserInfoInlineKeyboard(chatUser.ID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

func MarkovTalkAbout(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		if len(upd.Message.Text) < 6 {
			message := getUsernameOrFirstName(upd) + ": make sure you put what you want to ask about after the /ask command!"
			newMess := tgbotapi.NewMessage(settings.GetChannelID(), message)
			bot.Send(newMess)
		}
		tgbotapi.NewChatAction(settings.GetChannelID(), tgbotapi.ChatTyping)
		keyword := upd.Message.Text[5:]
		keyword = strings.Trim(keyword, " ")
		response := models.GenerateMarkovResponse(keyword)
		var message string
		if response == "" {
			message = "I haven't learned that word yet."
		} else {
			message = getUsernameOrFirstName(upd) + ": " + response
		}
		newMess := tgbotapi.NewMessage(settings.GetChannelID(), message)
		bot.Send(newMess)
	}
}

func getUsernameOrFirstName(upd tgbotapi.Update) string {
	if len(upd.Message.From.UserName) < 1 {
		return upd.Message.From.FirstName
	} else {
		return "@" + upd.Message.From.UserName
	}
}

func MarkovCount(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		message := "Current number of chains in database: " + strconv.Itoa(models.MarkovCount())
		newMess := tgbotapi.NewMessage(settings.GetControlID(), message)
		bot.Send(newMess)
	}
}
