package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type SDiscordConfig struct {
	BotToken  string
	AppID     string
	GuildID   string
	ChannelID string
}

type SArrConfig struct {
	NickName string
	APIKey   string
	URI      string
	Type     string
}

type SJellyfinConfig struct {
	URI       string
	PublicURI string
	APIKey    string
	Username  string
	Password  string
}

type SDatabase struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// ConfigData struct
type SConfig struct {
	Port      int
	IP        string
	PublicURL string
	Database  SDatabase
	Jellyfin  SJellyfinConfig
	Discord   SDiscordConfig
	Services  []SArrConfig
}

// Config object
var Config SConfig

func init() {
	// Initilaize flags first
	initFlags()

	// Set default values for config if config does not exist
	Config = SConfig{
		Discord:  SDiscordConfig{},
		Services: []SArrConfig{{}},
	}

	dataDir := *Flags.DataDir

	configPath := fmt.Sprintf("%s/config.json", dataDir)

	if _, err := os.Stat(dataDir); err != nil {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			panic(err)
		}
	}

	initLogs()

	err := GetConfig(configPath, &Config)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetConfig get a config
func GetConfig(configFileName string, configPointer interface{}) error {
	if fileExists(configFileName) { // Get existing configuration from configFileName
		b, err := os.ReadFile(configFileName)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, configPointer)
		if err != nil {
			fmt.Println("Failed to unmarshal configuration file")
			return err
		}

		return nil
	}

	// If configFileName doesn't exist, create a new config file
	b, err := json.MarshalIndent(configPointer, "", " ")
	if err != nil {
		fmt.Println("Failed to marshal configuration file")
		return err
	}

	// Sevae default config into to file.
	err = os.WriteFile(configFileName, b, 0755)
	if err != nil {
		fmt.Println("Failed to write configuration file")
		return err
	}

	return nil
}
