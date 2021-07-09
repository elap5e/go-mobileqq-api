package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/encoding/mark"
	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq"
	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
	"github.com/elap5e/go-mobileqq-api/pb"
	"github.com/elap5e/go-mobileqq-api/util"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+mobileqq.PackageName)
	username   string
	password   string
	chatID     uint64
	peerID     uint64
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
  - id: 0:10000
`, time.Now().UnixNano())

var reader = bufio.NewReader(os.Stdin)

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
		id := strings.Split(viper.GetString("targets.0.id"), ":")
		_ = id[1]
		chatID, _ = strconv.ParseUint(id[0], 10, 64)
		peerID, _ = strconv.ParseUint(id[1], 10, 64)
	}
}

func send(ctx context.Context, rpc *client.Client, text string) error {
	fromID, _ := strconv.ParseUint(username, 10, 64)
	chatName := strconv.Itoa(int(chatID))
	peerName := strconv.Itoa(int(peerID))
	fromName := strconv.Itoa(int(fromID))
	seq := rpc.GetNextMessageSeq(fmt.Sprintf("%d:%d", chatID, peerID))
	routingHead := &pb.RoutingHead{}
	if peerID == 0 {
		routingHead = &pb.RoutingHead{Group: &pb.Group{Code: chatID}}
	} else if chatID == 0 {
		routingHead = &pb.RoutingHead{C2C: &pb.C2C{Uin: peerID}}
	} else {
		routingHead = &pb.RoutingHead{
			GroupTemp: &pb.GroupTemp{Code: chatID, ToUin: peerID},
		}
	}

	msg := pb.Message{}
	if err := mark.Unmarshal([]byte(text), &msg); err != nil {
		return err
	}
	resp, err := rpc.MessageSendMessage(
		ctx, username, client.NewMessageSendMessageRequest(
			routingHead,
			msg.GetContentHead(),
			msg.GetMessageBody(),
			seq,
			nil,
		),
	)
	if err != nil {
		return err
	}

	log.PrintMessage(
		time.Unix(resp.GetSendTime(), 0),
		chatName, peerName, fromName, chatID, peerID, fromID, seq, text,
	)
	return nil
}

func main() {
	cfg := mobileqq.NewClientConfigFromViper()
	mqq := mobileqq.NewClient(&mobileqq.Options{
		BaseDir:  baseDir,
		LogLevel: viper.GetString("configs.logLevel"),
		Client:   cfg,
	})

	if err := mqq.Run(context.Background(), func(ctx context.Context, restart chan struct{}) error {
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
		errCh := make(chan error, 1)
		go func() {
			for {
				text, _ := util.ReadLine(reader)
				if err := send(ctx, rpc, text); err != nil {
					errCh <- err
				}
			}
		}()
		go func() {
			for range time.NewTicker(300 * time.Second).C {
				text := "![[困]](goqq://res/marketFace?id=ipEfT7oeSIPz3SIM7j4u5A==&tabId=204112&key=MmJjMGE1M2NmZDYyZjNkZg==)" + time.Now().Local().String()
				if err := send(ctx, rpc, text); err != nil {
					errCh <- err
				}
			}
		}()
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
