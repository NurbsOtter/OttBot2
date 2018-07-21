package settings

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	BotToken         string
	ChannelID        int64
	ControlChannelID int64
	SQLConnectString string
	RandomChance     int64
	DatabaseTimeout  string
}

var settings Settings

func LoadSettings() {
	data, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		panic(err)
	}
	settings = Settings{}
	json.Unmarshal(data, &settings)

}

func GetBotToken() string {
	return settings.BotToken
}
func GetChannelID() int64 {
	return settings.ChannelID
}
func GetControlID() int64 {
	return settings.ControlChannelID
}
func GetDBAddr() string {
	return settings.SQLConnectString
}
func GetRandomChance() int64 {
	return settings.RandomChance
}
func GetDatabaseTimeout() string {
	return settings.DatabaseTimeout
}
