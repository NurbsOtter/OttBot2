package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
)

var commands []*BotCommand
var callbacks []*BotCommand
var newmembercommands []*MemberCommand
var leftmembercommands []*MemberCommand

type BotCommand struct {
	MatchCmd   *regexp.Regexp
	Chan       int64
	HandleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)
}

type MemberCommand struct {
	Chan       int64
	HandleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)
}

func TestCmd(updateIn tgbotapi.Update, botIn *tgbotapi.BotAPI) {
	c := tgbotapi.NewMessage(updateIn.Message.Chat.ID, "Butts")
	botIn.Send(c)
}
func Register(regexIn string, chanIn int64, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)) {
	if commands == nil {
		commands = []*BotCommand{}
	}
	newCommand := &BotCommand{}
	newCommand.MatchCmd = regexp.MustCompile(regexIn)
	newCommand.Chan = chanIn
	newCommand.HandleFunc = handleFunc
	commands = append(commands, newCommand)
}
func RegisterNewMember(chanIn int64, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)) {
	if newmembercommands == nil {
		newmembercommands = []*MemberCommand{}
	}
	newCommand := &MemberCommand{}
	newCommand.Chan = chanIn
	newCommand.HandleFunc = handleFunc
	newmembercommands = append(newmembercommands, newCommand)
}
func RegisterLeftMember(chanIn int64, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)) {
	if leftmembercommands == nil {
		leftmembercommands = []*MemberCommand{}
	}
	newCommand := &MemberCommand{}
	newCommand.Chan = chanIn
	newCommand.HandleFunc = handleFunc
	leftmembercommands = append(leftmembercommands, newCommand)
}
func RegisterCallback(regexIn string, handleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)) {
	if callbacks == nil {
		callbacks = []*BotCommand{}
	}
	newCallback := &BotCommand{}
	newCallback.MatchCmd = regexp.MustCompile(regexIn)
	newCallback.HandleFunc = handleFunc
	callbacks = append(callbacks, newCallback)
}
func ProcessMessage(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := upd.Message.Text
	if message == "" {
		return
	}
	for _, cmd := range commands {
		if cmd.MatchCmd.Match([]byte(message)) {
			if cmd.Chan == upd.Message.Chat.ID || cmd.Chan == 0 {
				go cmd.HandleFunc(upd, bot)
			}
		}
	}
}
func ProcessNewMember(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := upd.Message.NewChatMember.ID
	if data == 0 {
		return
	}
	for _, cmd := range newmembercommands {
		if cmd.Chan == upd.Message.Chat.ID || cmd.Chan == 0 {
			go cmd.HandleFunc(upd, bot)
		}
	}
}
func ProcessLeftMember(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := upd.Message.LeftChatMember.ID
	if data == 0 {
		return
	}
	for _, cmd := range leftmembercommands {
		if cmd.Chan == upd.Message.Chat.ID || cmd.Chan == 0 {
			go cmd.HandleFunc(upd, bot)
		}
	}
}
func ProcessCallback(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := upd.CallbackQuery.Data
	if data == "" {
		return
	}
	for _, cback := range callbacks {
		if cback.MatchCmd.Match([]byte(data)) {
			go cback.HandleFunc(upd, bot)
		}
	}
}
func InitBot(botToken string) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}
	for update := range updates {
		if update.Message != nil {
			if update.Message.NewChatMember != nil {
				outLog := fmt.Sprintf("Member Joined: ID: %d Name: %s %s Username: %s", update.Message.NewChatMember.ID, update.Message.NewChatMember.FirstName, update.Message.NewChatMember.LastName, update.Message.NewChatMember.UserName)
				fmt.Println(outLog)
				//fmt.Println(update.Message.Chat.ID)
				ProcessNewMember(update, bot)
			} else if update.Message.LeftChatMember != nil {
				outLog := fmt.Sprintf("Member Left: ID: %d Name: %s %s Username: %s", update.Message.LeftChatMember.ID, update.Message.LeftChatMember.FirstName, update.Message.LeftChatMember.LastName, update.Message.LeftChatMember.UserName)
				fmt.Println(outLog)
				//fmt.Println(update.Message.Chat.ID)
				ProcessLeftMember(update, bot)
			} else {
				outLog := fmt.Sprintf("Message: %s %s>%s", update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)
				fmt.Println(outLog)
				//fmt.Println(update.Message.Chat.ID)
				ProcessMessage(update, bot)
			}
		} else if update.CallbackQuery != nil {
			fmt.Println("Cback")
			fmt.Println("Callback handled: " + update.CallbackQuery.Data)
			ProcessCallback(update, bot)
			//config := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			//fmt.Println(update.CallbackQuery.Data)
			//bot.AnswerCallbackQuery(config)
		} else {
			continue
		}
	}
}
