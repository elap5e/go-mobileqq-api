package mobileqq

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+PackageName)
	cacheDir   = path.Join(baseDir, "cache")
	logDir     = path.Join(baseDir, "log")
)

type Config struct {
	BaseDir  string
	CacheDir string

	AuthAddress string
	AuthCaptcha bool

	Engine *client.Config
}

func NewClientConfig() *Config {
	return NewClientConfigForAndroid()
}

func NewClientConfigForAndroid() *Config {
	return &Config{
		BaseDir:     baseDir,
		CacheDir:    cacheDir,
		AuthAddress: "127.0.0.1:0",
		AuthCaptcha: true,
		Engine: &client.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   client.NewClientConfig(),
			Device:   client.NewDeviceConfig(),
		},
	}
}

func NewClientConfigForAndroidTablet() *Config {
	return &Config{
		BaseDir:     baseDir,
		CacheDir:    cacheDir,
		AuthAddress: "127.0.0.1:0",
		AuthCaptcha: true,
		Engine: &client.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   client.NewClientConfigForAndroidTablet(),
			Device:   client.NewDeviceConfig(),
		},
	}
}

func NewClientConfigFromViper() *Config {
	cfg := &Config{
		BaseDir:     baseDir,
		CacheDir:    cacheDir,
		AuthAddress: "127.0.0.1:0",
		AuthCaptcha: true,
		Engine: &client.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   client.NewClientConfigForAndroidTablet(),
			Device:   client.NewDeviceConfig(),
		},
	}
	if viper.IsSet("configs.auth.address") {
		cfg.AuthAddress = viper.GetString("configs.auth.address")
	}
	if viper.IsSet("configs.auth.captcha") {
		cfg.AuthCaptcha = viper.GetBool("configs.auth.captcha")
	}
	if viper.IsSet("configs.deviceInfo.randomSeed") {
		cfg.Engine.Device = client.NewDeviceConfigBySeed(viper.GetInt64("configs.deviceInfo.randomSeed"))
	}
	if viper.IsSet("configs.protocol") {
		switch strings.ToLower(viper.GetString("configs.protocol")) {
		case "android":
			cfg.Engine.Client = client.NewClientConfigForAndroid()
		case "android-tablet":
			cfg.Engine.Client = client.NewClientConfigForAndroidTablet()
		}
	}
	return cfg
}
