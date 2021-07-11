package main

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/auth"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+mobileqq.PackageName)
	config     mobileqq.Config
)

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
				[]byte(mobileqq.ConfigYAML),
				0600,
			)
			log.Fatal().Msgf("x_x [init] create config.yaml in %s", configPath)
		} else {
			// Config file was found but another error was produced
			log.Fatal().Msgf("x_x [init] failed to load config.yaml")
		}
	} else {
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal().Err(err).Msg("x_x [init] failed to unmarshal config")
		}
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
		errCh := make(chan error, 1)
		wg := sync.WaitGroup{}
		for _, account := range config.Accounts {
			wg.Add(1)
			go func(username, password string) {
				defer wg.Done()
				rpc := mqq.GetClient()
				if err := auth.NewFlow(&auth.FlowOptions{
					Username: username,
					Password: password,
					AuthAddr: cfg.AuthAddress,
					CacheDir: cfg.CacheDir,
				}, auth.NewHandler(&auth.HandlerOptions{
					BaseDir: cfg.BaseDir,
					Client:  cfg.Engine.Client,
					Device:  cfg.Engine.Device,
				}, rpc)).Run(ctx); err != nil {
					errCh <- err
					return
				}
				uin, _ := strconv.ParseInt(username, 10, 64)
				if _, err := rpc.AccountUpdateStatus(ctx, client.NewAccountUpdateStatusRequest(
					uint64(uin),
					client.AccountStatusOnline,
					false,
				)); err != nil {
					errCh <- err
					return
				}
			}(account.Username, account.Password)
		}
		wg.Wait()
		select {
		case err := <-errCh:
			return err
		case <-restart:
			return nil
		}
	}); err != nil {
		log.Panic().Err(err).Msg("client unexpected exit")
	}
}
