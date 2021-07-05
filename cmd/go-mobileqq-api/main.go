package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
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
  debug: true
  deviceInfo:
    randomSeed: %d
  protocol: android-tablet
`, time.Now().UnixNano())

func init() {
	log.Info().Msg(log.MsgInitInfof("Go MobileQQ API (%s)", mobileqq.PackageVersion))
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
			log.Fatal().Msg(log.MsgInitErrorf("create config.yaml in %s", configPath))
		} else {
			// Config file was found but another error was produced
			log.Fatal().Msg(log.MsgInitError("failed to load config.yaml"))
		}
	} else {
		username = viper.GetString("accounts.0.username")
		password = viper.GetString("accounts.0.password")
	}
}

func main() {
	cfg := mobileqq.NewClientConfigFromViper()
	ctx := context.Background()
	engine := rpc.NewEngine(&rpc.Config{
		Network:     "tcp",
		Address:     "msfwifi.3g.qq.com:8080",
		FixID:       cfg.RPC.Client.AppID,
		AppID:       cfg.RPC.Client.AppID,
		NetworkType: 0x01,
		NetIPFamily: 0x03,
		IMEI:        cfg.RPC.Device.IMEI,
		IMSI:        cfg.RPC.Device.IMSI,
		Revision:    cfg.RPC.Client.Revision,
	})
	for {
		engine.Start(ctx)
		err := <-engine.Error()
		log.Error().
			Err(err).
			Msg("x-x [conn] failed to start rpc engine, retry in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
