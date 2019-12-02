package settings

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	BotToken             string `json:"BotToken"`
	ChannelID            int64  `json:"ChannelID"`
	ControlChannelID     int64  `json:"ControlChannelID"`
	SQLConnectString     string `json:"SQLConnectString"`
	RandomChance         int64  `json:"RandomChance"`
	DatabaseTimeout      string `json:"DatabaseTimeout"`
	AdminAnnounceChannel int64  `json:"AdminAnnounceChannel"`
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

func GetAnnounceChannel() int64 {
	return settings.AdminAnnounceChannel
}
