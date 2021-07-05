package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/rpc"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+mobileqq.PackageName)
	username   string
	password   string
)

func init() {
	log.Printf("~v~ [init] Go MobileQQ API (%s)", mobileqq.PackageVersion)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(baseDir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("x_x [init] failed to load config.yaml")
	} else {
		username = viper.GetString("accounts.0.username")
		password = viper.GetString("accounts.0.password")
	}
}

func main() {
	c := mobileqq.NewClient(
		mobileqq.Option{
			Config: mobileqq.NewClientConfigFromViper(),
		},
	)
	if err := c.Auth(username, password); err != nil {
		log.Printf("x_x [auth] error: %s", err.Error())
	}
	if err := c.AccountUpdateStatus(
		username,
		rpc.AccountStatusOnline,
		false,
	); err != nil {
		log.Printf("x_x [auth] error: %s", err.Error())
	}
	for range time.NewTicker(300 * time.Second).C {
		if err := c.MessageSendMessage(
			username,
			"\n![[困]](goqq://res/marketFace?id=ipEfT7oeSIPz3SIM7j4u5A==&tabId=204112&key=MmJjMGE1M2NmZDYyZjNkZg==)"+time.Now().Local().String(),
		); err != nil {
			log.Printf("x_x [test] error: %s", err.Error())
		}
	}
}
