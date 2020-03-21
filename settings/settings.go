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
	DatabaseTimeout      string `json:"DatabaseTimeout"`
	AdminAnnounceChannel int64  `json:"AdminAnnounceChannel"`
	CombotID             int64  `json:CombotID`
}

var settings Settings

func LoadSettings() error {
	data, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		return err
	}
	settings = Settings{}

	json.Unmarshal(data, &settings)
	return nil
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
func GetDatabaseTimeout() string {
	return settings.DatabaseTimeout
}
func GetAnnounceChannel() int64 {
	return settings.AdminAnnounceChannel
}
func GetCombotID() int64 {
	return settings.CombotID
}
