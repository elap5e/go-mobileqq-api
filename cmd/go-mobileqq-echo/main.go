package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
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
  logLevel: info
  protocol: android-tablet

targets:
  - uin: 0
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
			log.Fatal().Msgf("x_x [init] create config.yaml in %s", configPath)
		} else {
			// Config file was found but another error was produced
			log.Fatal().Msgf("x_x [init] failed to load config.yaml")
		}
	} else {
		username = viper.GetString("accounts.0.username")
		password = viper.GetString("accounts.0.password")
	}
}

func main() {
	ctx := context.Background()
	cfg := mobileqq.NewClientConfigFromViper()
	mqq := mobileqq.NewClient(&mobileqq.Options{
		BaseDir:  baseDir,
		LogLevel: viper.GetString("configs.logLevel"),
		Client:   cfg,
	})

	if err := mqq.Run(ctx, func(ctx context.Context, restart chan struct{}) error {
		if err := mqq.Auth(username, password); err != nil {
			return err
		}
		rpc := mqq.GetClient()
		uin, _ := strconv.ParseInt(username, 10, 64)
		if _, err := rpc.AccountUpdateStatus(ctx, client.NewAccountUpdateStatusRequest(
			uint64(uin),
			client.AccountStatusOnline,
			false,
		)); err != nil {
			return err
		}
		<-restart
		return nil
	}); err != nil {
		log.Panic().Err(err).Msg("client unexpected exit")
	}
}
