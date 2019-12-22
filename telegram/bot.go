package telegram

import (
	"fmt"
	"log"
	"regexp"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var commands []*BotCommand
var callbacks []*BotCommand

type BotCommand struct {
	MatchCmd   *regexp.Regexp
	Chan       int64
	HandleFunc func(tgbotapi.Update, *tgbotapi.BotAPI)
}

func TestCmd(updateIn tgbotapi.Update, botIn *tgbotapi.BotAPI) {
	c := tgbotapi.NewMessage(updateIn.Message.Chat.ID, "Pong")
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
	LogMessage(upd)
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

// AnswerCallback answers the callback query made in the provided update. Note: If text is not an empty string, it will be displayed to the user who initiated the callback query.
func AnswerCallback(upd tgbotapi.Update, bot *tgbotapi.BotAPI, text string) {
	callback := tgbotapi.NewCallback(upd.CallbackQuery.ID, text)
	bot.AnswerCallbackQuery(callback)
	LogCallback(upd, callback)
}

// LogCallback logs information from a processed callback
func LogCallback(upd tgbotapi.Update, callback tgbotapi.CallbackConfig) {
	if callback.Text != "" {
		log.Printf("Callback: %s\nFrom: %s %d\nText: %s", upd.CallbackQuery.Data, upd.CallbackQuery.From.String(), upd.CallbackQuery.From.ID, callback.Text)
	} else {
		log.Printf("Callback: %s\nFrom: %s %d", upd.CallbackQuery.Data, upd.CallbackQuery.From.String(), upd.CallbackQuery.From.ID)
	}
}

// LogCallback logs information from a processed callback
func LogMessage(upd tgbotapi.Update) {
	log.Printf("Message: %s\nFrom: %s %d", upd.Message.Text, upd.Message.From.String(), upd.Message.From.ID)
}

// LogCommand logs information from a processed command
func LogCommand(upd tgbotapi.Update, err error) {
	if err != nil {
		log.Printf("Command: %s\nFrom: %s %d\nError: %s", upd.Message.Command(), upd.Message.From.String(), upd.Message.From.ID, err.Error())
	} else {
		log.Printf("Command: %s\nFrom: %s %d\n", upd.Message.Command(), upd.Message.From.String(), upd.Message.From.ID)
	}
}

func InitBot(botToken string) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Failed to establish connection to bot with provided API token.")
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Failed to retrieve updates.")
		return
	}

	log.Print("Bot started successfully!")

	for update := range updates {
		if update.Message != nil {
			ProcessMessage(update, bot)
		} else if update.CallbackQuery != nil {
			ProcessCallback(update, bot)
		} else if update.ChannelPost != nil {
			fmt.Println(update.ChannelPost.Chat.ID)
		} else {
			continue
		}
	}
}
