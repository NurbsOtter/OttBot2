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

//Handles markov learning and random responses
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

//Handles toggling of a user's ability to use /ask
func ToggleMarkov(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
	chatUser := models.ChatUserFromID(userId)
	models.SetMarkovUse(chatUser.ID, !chatUser.MarkovAskAllowed)
	chatUser.MarkovAskAllowed = !chatUser.MarkovAskAllowed
	outMsg := GetUserInfoResponse(chatUser)
	editMsg := tgbotapi.NewEditMessageText(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, outMsg.Text)
	inlineKeyboard := MakeUserInfoInlineKeyboard(chatUser.ID)
	editMsg.ReplyMarkup = &inlineKeyboard
	bot.Send(editMsg)
}

//Handles /ask commands
func MarkovTalkAbout(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetChannelID() {
		foundUser := models.ChatUserFromTGID(upd.Message.From.ID, upd.Message.From.UserName)
		if foundUser.MarkovAskAllowed == true {
			if len(upd.Message.Text) < 6 {
				message := getUsernameOrFirstName(upd) + ": make sure you put what you want to ask about after the /ask command!"
				newMess := tgbotapi.NewMessage(settings.GetChannelID(), message)
				bot.Send(newMess)
			}
			tgbotapi.NewChatAction(settings.GetChannelID(), tgbotapi.ChatTyping)
			message := upd.Message.Text[5:]
			message = strings.Trim(message, " ")
			words := strings.Split(message, " ")
			generated := models.GenerateMarkovResponse(words[0])
			var response string
			if generated == "" {
				response = "I haven't learned that word yet."
			} else {
				response = getUsernameOrFirstName(upd) + ": " + generated
			}
			newMess := tgbotapi.NewMessage(settings.GetChannelID(), response)
			bot.Send(newMess)
		} else {
			newMsg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Sorry, you are banned from /ask")
			bot.Send(newMsg)
		}
	}
}

//Helper function to determine if a user should be addressed by username or first name
func getUsernameOrFirstName(upd tgbotapi.Update) string {
	if len(upd.Message.From.UserName) < 1 {
		return upd.Message.From.FirstName
	} else {
		return "@" + upd.Message.From.UserName
	}
}

//Returns the number of rows in markov database
func MarkovCount(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if upd.Message.Chat.ID == settings.GetControlID() {
		message := "Current number of chains in database: " + strconv.Itoa(models.MarkovCount())
		newMess := tgbotapi.NewMessage(settings.GetControlID(), message)
		bot.Send(newMess)
	}
}
