package mobileqq

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/mobileqq/client/config"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+PackageName)
	cacheDir   = path.Join(baseDir, "cache")
	logDir     = path.Join(baseDir, "log")
)

var ConfigYAML = fmt.Sprintf(`# Go MobileQQ API Configuration Template

accounts:
  - username: 10000
    password: 123456
    botToken: 10000:ABC-DEF1234ghIkl-zyx57W2v1u123ew11

configs:
  auth:
    address: localhost:0
    captcha: true
  database:
    dataSourceName: goqqd.db?parseTime=true
    driverName: sqlite3
  deviceInfo:
    randomSeed: %d
  logLevel: info
  protocol: android

servers:
  endpoints:
    - socket://msfwifi.3g.qq.com:8080
    - socket://msfwifiv6.3g.qq.com:8080
  forceIPv6: false
  overwrite: false

targets:
  - peerId: 0
    userId: 10000
`, time.Now().UnixNano())

type Config struct {
	Accounts []struct {
		Username string `json:"username"`
		Password string `json:"password"`
		BotToken string `json:"botToken"`
	} `json:"accounts"`
	Configs struct {
		Auth struct {
			Address string `json:"address"`
			Captcha bool   `json:"captcha"`
		} `json:"auth"`
		Database struct {
			DataSourceName string `json:"dataSourceName"`
			DriverName     string `json:"driverName"`
		} `json:"database"`
		DeviceInfo struct {
			RandomSeed int64 `json:"randomSeed"`
		} `json:"deviceInfo"`
		LogLevel string `json:"logLevel"`
		Protocol string `json:"protocol"`
	} `json:"configs"`
	Servers *struct {
		Endpoints []string `json:"endpoints,omitempty"`
		ForceIPv6 bool     `json:"forceIPv6,omitempty"`
		Overwrite bool     `json:"overwrite,omitempty"`
	} `json:"servers,omitempty"`
	Targets []*struct {
		ChatID string `json:"chatId,omitempty"`
		PeerID int64  `json:"peerId,omitempty"`
		UserID int64  `json:"ueerId,omitempty"`
	} `json:"targets,omitempty"`
}

type ClientConfig struct {
	BaseDir  string
	CacheDir string

	Engine *config.Config
}

func NewClientConfig() *ClientConfig {
	return NewClientConfigForAndroid()
}

func NewClientConfigForAndroid() *ClientConfig {
	return &ClientConfig{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Engine: &config.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   config.NewClientConfig(),
			Device:   config.NewDeviceConfig(),
		},
	}
}

func NewClientConfigForAndroidTablet() *ClientConfig {
	return &ClientConfig{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Engine: &config.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   config.NewClientConfigForAndroidTablet(),
			Device:   config.NewDeviceConfig(),
		},
	}
}

func NewClientConfigFromViper() *ClientConfig {
	cfg := &ClientConfig{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Engine: &config.Config{
			BaseDir:  baseDir,
			CacheDir: cacheDir,
			Client:   config.NewClientConfig(),
			Device:   config.NewDeviceConfig(),
		},
	}
	if viper.IsSet("configs.deviceInfo.randomSeed") {
		cfg.Engine.Device = config.NewDeviceConfigBySeed(viper.GetInt64("configs.deviceInfo.randomSeed"))
	}
	if viper.IsSet("configs.protocol") {
		switch strings.ToLower(viper.GetString("configs.protocol")) {
		case "android":
			cfg.Engine.Client = config.NewClientConfigForAndroid()
		case "android-tablet":
			cfg.Engine.Client = config.NewClientConfigForAndroidTablet()
		}
	}
	return cfg
}
