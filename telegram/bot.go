package telegram

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"regexp"
)

var bot *tgbotapi.BotAPI
var commands []*BotCommand

type BotCommand struct {
	MatchCmd   *regexp.Regexp
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
func ProcessMessage(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := upd.Message.Text
	if message == "" {
		return
	}
	for _, cmd := range commands {
		if cmd.MatchCmd.Match([]byte(message)) {
			cmd.HandleFunc(upd, bot)
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
		if update.Message == nil {
			continue
		}
		fmt.Println(update.Message.Chat.ID)
		ProcessMessage(update, bot)
	}
}
