package telegram

import (
	"OttBot2/models"
	"OttBot2/settings"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

//Handles non-command messages to record user information/changes
func HandleUsers(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	foundUser, err := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
	if err != nil {
		LogCommand(upd, err)
	}

	err = models.UpdateAliases(upd.Message.From.FirstName, upd.Message.From.LastName, foundUser.TgID)
	// This is disabled to avoid creating a duplicate log on every message, since this commmand essentially runs on every message
	//LogCommand(upd, err)
}

//Get user information by telegram ID
func FindUserByUserID(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		userId := upd.Message.Text[5:]
		userId = strings.Trim(userId, " ")
		usrID, err := strconv.ParseInt(userId, 10, 64)

		// Invalid ID
		if err != nil {
			newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Unable to parse userID. Make sure it's an actual number!")
			bot.Send(newMsg)
			LogCommand(upd, err)
			return
		}

		// Valid ID
		user, err := models.ChatUserFromTGIDNoUpd(int(usrID))
		bot.Send(GetUserInfoResponse(user))
		LogCommand(upd, err)
	}
}

//Get user information by username
func FindUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		username := upd.Message.Text[7:]
		username = strings.Trim(username, " ")
		username = strings.ToLower(username)
		user, err := models.SearchUserByUsername(username)

		// Problem finding user
		if err != nil {
			bot.Send(GetUserInfoResponse(user))
			LogCommand(upd, err)
			return
		}

		curAlias, err := models.GetLatestAliasFromUserID(user.TgID)
		// Problem retrieving alias
		if err != nil {
			newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Unexpected error while retrieving user's aliases.")
			bot.Send(newMsg)
			LogCommand(upd, err)
			return
		}

		outMsg := fmt.Sprintf("User ID: %d", user.TgID)
		if user.UserName != "" {
			outMsg += fmt.Sprintf("\nUsername: @%s", user.UserName)
		}
		outMsg += fmt.Sprintf("\nCurrent name: %s\nMod ping: %t\n", curAlias.Name, user.PingAllowed)
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), outMsg)
		newMsg.ReplyMarkup = MakeUserInfoInlineKeyboard(user.TgID)
		bot.Send(newMsg)
		LogCommand(upd, err)
	}
}

//Helper method to generate the response object for the info requests
func GetUserInfoResponse(user *models.ChatUser) tgbotapi.MessageConfig {
	// User found
	if user != nil {
		curAlias, err := models.GetLatestAliasFromUserID(user.TgID)
		// Problem retrieving alias
		if err != nil {
			newMsg := tgbotapi.NewMessage(settings.GetControlID(), "Unexpected error while retrieving user's aliases.")
			return newMsg
		}

		// Alias found
		outMsg := fmt.Sprintf("User ID: %d", user.TgID)
		if user.UserName != "" {
			outMsg += fmt.Sprintf("\nUsername: @%s", user.UserName)
		}
		outMsg += fmt.Sprintf("\nCurrent name: %s\nMod ping: %t\n", curAlias.Name, user.PingAllowed)
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), outMsg)
		newMsg.ReplyMarkup = MakeUserInfoInlineKeyboard(user.TgID)
		return newMsg

	} else {
		// User not found
		newMsg := tgbotapi.NewMessage(settings.GetControlID(), "User not found.")
		return newMsg
	}
}

//Helper method to generate the buttons for an initial info request
func MakeUserInfoInlineKeyboard(userId int64) tgbotapi.InlineKeyboardMarkup {
	warnCmd := fmt.Sprintf("/getwarnings %d", userId)
	warnButt := tgbotapi.NewInlineKeyboardButtonData("View warnings", warnCmd)
	aliasCmd := fmt.Sprintf("/getaliases %d 0", userId)
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
	aliasCmd := fmt.Sprintf("/getaliases %d 0", userId)
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
func MakeUserInfoInlineKeyboardRefreshAliasButton(userId int64, curAliasPage int64, aliasPagesTotal int64) tgbotapi.InlineKeyboardMarkup {
	warnCmd := fmt.Sprintf("/getwarnings %d", userId)
	warnButt := tgbotapi.NewInlineKeyboardButtonData("View warnings", warnCmd)
	var aliasButt tgbotapi.InlineKeyboardButton
	if curAliasPage == aliasPagesTotal-1 || aliasPagesTotal == 0 {
		aliasCmd := fmt.Sprintf("/getaliases %d 0", userId)
		aliasButt = tgbotapi.NewInlineKeyboardButtonData("Refresh aliases", aliasCmd)
	} else {
		aliasCmd := fmt.Sprintf("/getaliases %d %d", userId, curAliasPage+1)
		aliasButt = tgbotapi.NewInlineKeyboardButtonData("Next Page", aliasCmd)
	}
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

	// Check to see if this is an appropriate CallbackQuery
	if upd.Message == nil {
		userIdStr = strings.Fields(upd.CallbackQuery.Data)[1]
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	// Parse the user ID from the InlineKeyboard's CallbackQuery data
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	chatUser, err := models.ChatUserFromID(userId)
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	infoMessage := GetUserInfoResponse(chatUser)
	outmsg := infoMessage.Text
	warnings, err := models.GetUsersWarnings(chatUser)
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	// Append all found warnings
	if len(warnings) > 0 {
		outmsg += "\nWarnings:"
		for _, warn := range warnings {
			outmsg += "\n" + warn.WarningText
		}
	} else {
		outmsg += "\n No warnings found."
	}
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outmsg)
	inlineKeyboard := MakeUserInfoInlineKeyboardRefreshWarnButton(chatUser.TgID)
	editMsg.ReplyMarkup = &inlineKeyboard

	bot.Send(editMsg)
	AnswerCallback(upd, bot, "")
}

//Callback handler to update a find by alias request after a user button is clicked
func CallbackInfo(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string

	// Check to see if this is an appropriate CallbackQuery
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[13:]
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// Unable to parse user id
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	chatUser, err := models.ChatUserFromID(userId)
	// Unable to retrieve chat user
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	infoMessage := GetUserInfoResponse(chatUser)
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, infoMessage.Text)
	inlineKeyboard := MakeUserInfoInlineKeyboard(chatUser.TgID)
	editMsg.ReplyMarkup = &inlineKeyboard

	bot.Send(editMsg)
	AnswerCallback(upd, bot, "")
}

//Callback handler to update a get user info response to add all known user aliases
func DisplayAliases(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	var curPage int64

	// Check to see if this is an appropriate CallbackQuery
	if upd.Message == nil {
		splitString := strings.Split(upd.CallbackQuery.Data, " ")
		userIdStr = splitString[1]
		curPage, _ = strconv.ParseInt(splitString[2], 10, 32)
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// Unable to parse user id
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	chatUser, err := models.ChatUserFromID(userId)
	// Unable to retrieve chat user
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	infoMessage := GetUserInfoResponse(chatUser)
	outmsg := infoMessage.Text
	aliases, err := models.GetAliases(chatUser)
	// Unable to retrieve aliases
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	// Append all found aliases
	if len(aliases) > 0 {
		outmsg += "\nKnown aliases:"
		if len(aliases) <= 10 {
			for _, alias := range aliases {
				outmsg += "\n" + alias.Name
			}
		} else {
			page := aliases[curPage*5:]
			for _, alias := range page[:5] {
				outmsg += "\n" + alias.Name
			}
		}
	} else {
		outmsg += "\n No Aliases found"
	}

	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outmsg)
	inlineKeyboard := MakeUserInfoInlineKeyboardRefreshAliasButton(chatUser.TgID, curPage, int64(len(aliases)/5))
	editMsg.ReplyMarkup = &inlineKeyboard

	bot.Send(editMsg)
	AnswerCallback(upd, bot, "")
}

//Warn a user by username
func WarnUserByUsername(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[7:] //Remove the /warn @
		procString = strings.TrimLeft(procString, " ")
		userName := strings.Split(procString, " ")[0]
		userName = strings.ToLower(userName)
		message := procString[len(userName)+1:]

		user, err := models.SearchUserByUsername(userName)
		// Unable to find user
		if err != nil {
			newMess := tgbotapi.NewMessage(settings.GetControlID(), "Could not find user")
			bot.Send(newMess)
			LogCommand(upd, err)
			return
		}

		// Found user
		models.AddWarningToID(user.TgID, message)
		newMess := tgbotapi.NewMessage(settings.GetControlID(), "Warned "+userName)
		bot.Send(newMess)
		LogCommand(upd, err)
	}
}

//Warn a user by telegram ID
func WarnUserByID(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[6:] //Remove the /warn
		procString = strings.TrimLeft(procString, " ")
		userID, err := strconv.ParseInt(strings.Split(procString, " ")[0], 10, 64)

		// Unable to parse user id
		if err != nil {
			newMess := tgbotapi.NewMessage(settings.GetControlID(), "Invalid user id")
			bot.Send(newMess)
			LogCommand(upd, err)
			return
		}

		// Problem finding user
		chatUser, err := models.ChatUserFromTGIDNoUpd(int(userID))
		// TODO: Rework ChatUserFromTGIDNoUpd function then clean up this function
		if err != nil || chatUser == nil {
			newMess := tgbotapi.NewMessage(settings.GetControlID(), "Could not find user")
			bot.Send(newMess)
			LogCommand(upd, err)
			return
		}

		var outMsg string
		err = models.AddWarningToID(chatUser.TgID, procString[len(strings.Split(procString, " ")[0])+1:])
		// Problem adding warning
		if err != nil {
			newMess := tgbotapi.NewMessage(settings.GetControlID(), "Encountered problem while trying to warn user")
			bot.Send(newMess)
			LogCommand(upd, err)
			return
		}

		outMsg = fmt.Sprintf("Warned %d", userID)
		newMess := tgbotapi.NewMessage(settings.GetControlID(), outMsg)
		bot.Send(newMess)
		LogCommand(upd, err)
	}
}

// LookupAlias finds a user even if only part of one of their aliases is provided
func LookupAlias(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		procString := upd.Message.Text[6:]
		procString = strings.Trim(procString, " ")

		foundAliases, err := models.LookupAlias(strings.ToLower(procString))
		// Problem finding aliases
		if err != nil {
			newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Unexpected error while retrieving user's aliases.")
			bot.Send(newMsg)
			LogCommand(upd, err)
			return
		}

		// Aliases not found
		if len(foundAliases) == 0 {
			outMsg := tgbotapi.NewMessage(settings.GetControlID(), "No Aliases found")
			bot.Send(outMsg)
			LogCommand(upd, err)
			return
		}

		// Aliases found
		outMsg := tgbotapi.NewMessage(settings.GetControlID(), "Found Users:")
		if len(foundAliases) > 8 {
			foundAliases = foundAliases[:8]
		}
		outMsg.ReplyMarkup = MakeAliasInlineKeyboard(foundAliases)
		bot.Send(outMsg)
		LogCommand(upd, err)
	}
}

//Helper method to generate the buttons for the lookup by alias command
func MakeAliasInlineKeyboard(aliases []models.ChatUser) tgbotapi.InlineKeyboardMarkup {
	var aliasButtons []tgbotapi.InlineKeyboardButton
	for _, alias := range aliases {
		latestID, err := models.GetLatestAliasFromUserID(alias.TgID)
		// Problem retrieving latest alias
		if err != nil {
			return tgbotapi.NewInlineKeyboardMarkup()
		}

		btnCmd := fmt.Sprintf("/callbackinfo %d", alias.TgID)
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
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// Unable to parse id
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	foundUser, err := models.ChatUserFromID(userId)
	// Unable to retrieve user
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	infoMessage := GetUserInfoResponse(foundUser)
	outMsg := infoMessage.Text
	outMsg += "\nDo you want to ban this user?"
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg)
	inlineKeyboard := MakeBanInlineKeyboard(foundUser.TgID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
	AnswerCallback(upd, bot, "")
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
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// Unable to parse user id
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	foundUser, err := models.ChatUserFromID(userId)
	// Unable to retrieve user
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	infoMessage := GetUserInfoResponse(foundUser)
	outMsg := infoMessage.Text
	outMsg += "\nAre you ABSOLUTELY SURE you want to ban this user?"
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg)
	inlineKeyboard := MakeBanConfirmInlineKeyboard(userId)
	editMsg.ReplyMarkup = &inlineKeyboard

	bot.Send(editMsg)
	AnswerCallback(upd, bot, "")
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
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// Unable to parse user id
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	foundUser, err := models.ChatUserFromID(userId)
	// Unable to retrieve user
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	banConfig := tgbotapi.KickChatMemberConfig{}
	banConfig.ChatID = settings.GetChannelID()
	banConfig.UserID = int(foundUser.TgID)
	_, err = bot.KickChatMember(banConfig)
	// Error kicking chat member
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

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
	AnswerCallback(upd, bot, "")
}

//Handles toggling of a user's ability to use /mods
func ToggleMods(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var userIdStr string
	if upd.Message == nil {
		userIdStr = upd.CallbackQuery.Data[12:]
	} else {
		AnswerCallback(upd, bot, "")
		return
	}

	userIdStr = strings.Trim(userIdStr, " ")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	// Unable to parse user id
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	chatUser, err := models.ChatUserFromID(userId)
	// Unable to retrieve user
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	err = models.SetModPing(chatUser.TgID, !chatUser.PingAllowed)
	// Error setting mod ping
	if err != nil {
		AnswerCallback(upd, bot, err.Error())
		return
	}

	chatUser.PingAllowed = !chatUser.PingAllowed
	outMsg := GetUserInfoResponse(chatUser)
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg.Text)
	inlineKeyboard := MakeUserInfoInlineKeyboard(chatUser.TgID)
	editMsg.ReplyMarkup = &inlineKeyboard

	bot.Send(editMsg)
	AnswerCallback(upd, bot, "")
}
