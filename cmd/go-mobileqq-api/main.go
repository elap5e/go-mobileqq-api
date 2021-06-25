package main

import (
	"log"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/rpc"
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
	rpc.SetClientForAndroidPad()
}

func main() {
	c := mobileqq.NewClient()
	if err := c.HeartbeatAlive(); err != nil {
		log.Printf("x_x [test] error: %s", err.Error())
	}
	if err := c.Auth(username, password); err != nil {
		log.Printf("x_x [auth] error: %s", err.Error())
	}
}
