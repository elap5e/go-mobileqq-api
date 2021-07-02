package main

import (
	"fmt"
	"io/ioutil"
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

var configYAML = fmt.Sprintf(`# Go MobileQQ API Configuration Template

accounts:
  - username: 10000
    password: 123456

configs:
  auth:
    address: 127.0.0.1:0
    captcha: true
  deviceInfo:
    randomSeed: %d
  protocol: android-tablet
`, time.Now().UnixNano())

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(baseDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			configPath := path.Join(baseDir, "config.yaml")
			_ = ioutil.WriteFile(
				configPath,
				[]byte(configYAML),
				0600,
			)
			log.Fatalf("$_$ [init] create config.yaml in %s", configPath)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("x_x [init] failed to load config.yaml")
		}
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
	if err := c.HeartbeatAlive(); err != nil {
		log.Printf("x_x [test] error: %s", err.Error())
	}
	if err := c.Auth(username, password); err != nil {
		log.Printf("x_x [auth] error: %s", err.Error())
	}
	if err := c.AccountUpdateStatus(
		username,
		rpc.PushRegisterInfoStatusOnline,
		false,
	); err != nil {
		log.Printf("x_x [auth] error: %s", err.Error())
	}
	select {}
}
