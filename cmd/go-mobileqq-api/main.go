package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/mobileqq"
)

var (
	username string
	password string
)

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	username = viper.GetString("GMA_USERNAME")
	password = viper.GetString("GMA_PASSWORD")
	log.Printf("username %s, password %s", username, password)
	// rpc.SetClientForAndroidPad()
}

func main() {
	c := mobileqq.NewClient()
	c.Auth(username, password)
}
