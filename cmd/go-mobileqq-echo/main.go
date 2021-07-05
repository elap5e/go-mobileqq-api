package main

import (
	"context"
	"os"
	"path"
	"strconv"

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

func init() {
	log.Info().Msgf("··· [init] Go MobileQQ API (%s)", mobileqq.PackageVersion)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(baseDir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Msg("x_x [init] failed to load config.yaml")
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

	if err := mqq.Run(ctx, func(ctx context.Context) error {
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
		select {}
	}); err != nil {
		log.Panic().Err(err).Msg("client unexpected exit")
	}
}
