package main

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/mobileqq/api"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client/auth"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+mobileqq.PackageName)
	config     mobileqq.Config
	tokens     = make(map[string]string)
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
		for _, account := range config.Accounts {
			tokens[account.Username] = account.BotToken
		}
	}
}

func main() {
	cfg := mobileqq.NewClientConfigFromViper()
	mqq := mobileqq.NewClient(&mobileqq.Options{
		BaseDir:  baseDir,
		LogLevel: viper.GetString("configs.logLevel"),
		Client:   cfg,
	})

	dbConfig := config.Configs.Database
	if strings.HasPrefix(dbConfig.DriverName, "sqlite") {
		dbConfig.DataSourceName = cfg.BaseDir + path.Join("/", dbConfig.DataSourceName)
	}
	db, err := sql.Open(dbConfig.DriverName, dbConfig.DataSourceName)
	if err != nil {
		log.Panic().Err(err).Msg("db unexpected exit")
	}
	defer db.Close()
	mqq.GetClient().SetDB(db)

	srv := api.NewServer(tokens)
	go func() {
		if err := srv.Run(context.Background()); err != nil {
			log.Panic().Err(err).Msg("bot api server unexpected exit")
		}
	}()

	if err := mqq.Run(
		context.Background(),
		func(ctx context.Context, once bool, restart chan struct{}) error {
			errCh := make(chan error, 1)
			wg := sync.WaitGroup{}
			for _, account := range config.Accounts {
				if once {
					wg.Add(1)
					go func(username, password string) {
						defer wg.Done()
						rpc := mqq.GetClient()
						if err := auth.NewFlow(&auth.FlowOptions{
							Username: username,
							Password: password,
							AuthAddr: config.Configs.Auth.Address,
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
						if _, err := rpc.AccountSetStatus(
							ctx,
							client.NewAccountSetStatusRequest(
								uint64(uin),
								client.AccountStatusOnline,
								false,
							),
						); err != nil {
							errCh <- err
							return
						}
						if _, err := rpc.FriendListGetFriendGroupList(ctx,
							client.NewFriendListGetFriendGroupListRequest(
								uin, 0, 100, 0, 100,
							),
						); err != nil {
							errCh <- err
							return
						}
						if _, err := rpc.FriendListGetGroupList(ctx,
							client.NewFriendListGetGroupListRequest(
								uin, nil,
							),
						); err != nil {
							errCh <- err
							return
						}
						if _, err := rpc.MessageGetMessage(ctx,
							username,
							client.NewMessageGetMessageRequest(
								0x00000000, nil,
							),
						); err != nil {
							errCh <- err
							return
						}
					}(account.Username, account.Password)
				}
				wg.Wait()
				select {
				case err := <-errCh:
					return err
				default:
				}
			}

			srv.ResetClient(ctx, mqq.GetClient())
			select {
			case err := <-errCh:
				return err
			case <-restart:
				return nil
			}
		},
	); err != nil {
		log.Panic().Err(err).Msg("api client unexpected exit")
	}
}
