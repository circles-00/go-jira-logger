package config

import (
	"errors"
	"fmt"
	"os"

	assert "go_jira_logger/pkg/utils/assert"

	"github.com/spf13/viper"
)

const (
	ConfigPath = ".config/go_jira_logger"
	ConfigName = "config"
)

func GetConfigDirPath() string {
	homeDirectory, err := os.UserHomeDir()
	assert.NoError(err, "User Home directory cannot be obtained!")

	return fmt.Sprintf("%s/%s", homeDirectory, ConfigPath)
}

func CheckIfConfigDirExists() bool {
	file, err := os.OpenFile(GetConfigDirPath(), os.O_RDONLY, 0o644)
	file.Close()

	return !errors.Is(err, os.ErrNotExist)
}

func SetConfigParams() {
	viper.SetConfigName(ConfigName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(GetConfigDirPath())
}

func SetDefaultConfigParams() {
	SetConfigParams()

	viper.SetDefault("jira.board_url", "https://<some-board>.atlassian.net")
	viper.SetDefault("jira.email", "john_doe@example.com")
	viper.SetDefault("jira.token", "<some_token>")

	checkIfConfigDirExists := CheckIfConfigDirExists()

	if !checkIfConfigDirExists {
		err := os.Mkdir(GetConfigDirPath(), 0o755)
		assert.NoError(err, "Could not create configuration directory!")
	}
}

func ReadConfigFile() {
	SetConfigParams()

	err := viper.ReadInConfig()
	assert.NoError(err, "Could not read configuration file")
}

func GetBoardUrl() string {
	return viper.GetString("jira.board_url")
}

func GetJiraEmail() string {
	return viper.GetString("jira.email")
}

func GetJiraToken() string {
	return viper.GetString("jira.token")
}
