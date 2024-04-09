package db

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var connStr string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error getting the current file path.")
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..")
	configPath := filepath.Join(projectRoot, "config", "test.yml")
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading testing configuration file: " + err.Error())
	}
	viper.SetDefault("environment.test.database.user", os.Getenv("TEST_USER"))
	viper.SetDefault("environment.test.database.password", os.Getenv("TEST_PASSWORD"))
	viper.SetDefault("environment.test.database.name", os.Getenv("TEST_DB"))
	viper.SetDefault("environment.test.database.connection_string", os.Getenv("TEST_CONN_STR"))

	connStr = viper.GetString("environment.test.database.connection_string")
	if connStr == "" {
		panic("Connection string not found in configuration")
	}
}
