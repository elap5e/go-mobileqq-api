package mobileqq

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"

	"github.com/elap5e/go-mobileqq-api/log"
	"github.com/elap5e/go-mobileqq-api/mobileqq/rpc"
)

type Options struct {
	BaseDir  string
	CacheDir string
	LogDir   string
	LogLevel string
	Client   *Config
	Engine   *rpc.Config
}

func (opt *Options) init() {
	if opt.BaseDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		opt.BaseDir = path.Join(homeDir, "."+PackageName)
	}
	if opt.CacheDir == "" {
		opt.CacheDir = path.Join(opt.BaseDir, "cache")
	}
	if opt.LogDir == "" {
		opt.LogDir = path.Join(opt.BaseDir, "log")
	}
	if opt.LogLevel == "" {
		opt.LogLevel = "info"
	}
	logLevel, err := zerolog.ParseLevel(opt.LogLevel)
	if err != nil {
		log.Error().Msg(err.Error())
		opt.LogLevel = "info"
		logLevel = zerolog.InfoLevel
	} else {
		opt.LogLevel = logLevel.String()
	}
	for _, dir := range []string{opt.BaseDir, opt.CacheDir, opt.LogDir} {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
		}
		if err != nil {
			panic(err)
		}
	}
	logFile, err := os.OpenFile(path.Join(
		logDir,
		fmt.Sprintf(
			"goqq-%s.log",
			time.Now().Local().Format("20060102150405"),
		),
	), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	logFileWriter := log.ConsoleWriterNoColor
	logFileWriter.Out = logFile
	log.Logger = zerolog.New(
		zerolog.MultiLevelWriter(log.ConsoleWriter, logFileWriter),
	).With().Timestamp().Logger().Level(logLevel)
	opt.Engine = &rpc.Config{
		Network:     "tcp",
		Address:     "msfwifi.3g.qq.com:8080",
		FixID:       opt.Client.Engine.Client.AppID,
		AppID:       opt.Client.Engine.Client.AppID,
		NetworkType: 0x01,
		NetIPFamily: 0x03,
		IMEI:        opt.Client.Engine.Device.IMEI,
		IMSI:        opt.Client.Engine.Device.IMSI,
		Revision:    opt.Client.Engine.Client.Revision,
	}
}
