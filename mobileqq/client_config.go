package mobileqq

import (
	"os"
	"path"

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
	NetworkType string
	NetIPFamily string
}

func NewClientConfig() *Config {
	return &Config{
		BaseDir:  baseDir,
		CacheDir: cacheDir,
		Client: &ClientConfig{
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
