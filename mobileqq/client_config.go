package mobileqq

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/rpc"
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
	Client   *ClientConfig
	RPC      *rpc.Config
}

type ClientConfig struct {
	AuthAddress string
	AuthCaptcha bool

	NetworkType string
	NetIPFamily string // None, IPv4, IPv6, Dual
}

func NewClientConfig() *Config {
	return NewClientConfigForAndroid()
}

func NewClientConfigForAndroid() *Config {
	return &Config{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Client: &ClientConfig{
			AuthAddress: "127.0.0.1:0",
			AuthCaptcha: true,
			NetworkType: "Wi-Fi",
			NetIPFamily: "Dual",
		},
		RPC: &rpc.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   rpc.NewClientConfig(),
			Device:   rpc.NewDeviceConfig(),
		},
	}
}

func NewClientConfigForAndroidTablet() *Config {
	return &Config{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Client: &ClientConfig{
			AuthAddress: "127.0.0.1:0",
			AuthCaptcha: true,
			NetworkType: "Wi-Fi",
			NetIPFamily: "Dual",
		},
		RPC: &rpc.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   rpc.NewClientConfigForAndroidTablet(),
			Device:   rpc.NewDeviceConfig(),
		},
	}
}

func NewClientConfigFromViper() *Config {
	cfg := &Config{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Client: &ClientConfig{
			AuthAddress: "127.0.0.1:0",
			AuthCaptcha: true,
			NetworkType: "Wi-Fi",
			NetIPFamily: "Dual",
		},
		RPC: &rpc.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			LogLevel: 0b00011111,
			Client:   rpc.NewClientConfigForAndroidTablet(),
			Device:   rpc.NewDeviceConfig(),
		},
	}
	if viper.IsSet("configs.auth.address") {
		cfg.Client.AuthAddress = viper.GetString("configs.auth.address")
	}
	if viper.IsSet("configs.auth.captcha") {
		cfg.Client.AuthCaptcha = viper.GetBool("configs.auth.captcha")
	}
	if viper.IsSet("configs.logLevel") {
		switch strings.ToLower(viper.GetString("configs.logLevel")) {
		case "error":
			cfg.RPC.LogLevel = 0b00000011
		case "warn":
			cfg.RPC.LogLevel = 0b00001111
		case "info":
			cfg.RPC.LogLevel = 0b00011111
		case "debug":
			cfg.RPC.LogLevel = 0b01111111
		case "trace":
			cfg.RPC.LogLevel = 0b11111111
		}
	}
	if viper.IsSet("configs.networkType") {
		switch strings.ToLower(viper.GetString("configs.networkType")) {
		case "wifi", "wi-fi":
			cfg.Client.NetworkType = "Wi-Fi"
		}
	}
	if viper.IsSet("configs.netIPFamily") {
		switch strings.ToLower(viper.GetString("configs.netIPFamily")) {
		case "none":
			cfg.Client.NetIPFamily = "None"
		case "ipv4":
			cfg.Client.NetIPFamily = "IPv4"
		case "ipv6":
			cfg.Client.NetIPFamily = "IPv6"
		case "dual":
			cfg.Client.NetIPFamily = "Dual"
		}
	}
	if viper.IsSet("configs.deviceInfo.randomSeed") {
		cfg.RPC.Device = rpc.NewDeviceConfigBySeed(viper.GetInt64("configs.deviceInfo.randomSeed"))
	}
	if viper.IsSet("configs.protocol") {
		switch strings.ToLower(viper.GetString("configs.protocol")) {
		case "android":
			cfg.RPC.Client = rpc.NewClientConfigForAndroid()
		case "android-tablet":
			cfg.RPC.Client = rpc.NewClientConfigForAndroidTablet()
		}
	}
	return cfg
}
