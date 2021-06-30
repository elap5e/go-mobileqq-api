package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/elap5e/go-mobileqq-api/mobileqq"
)

var (
	homeDir, _ = os.UserHomeDir()
	baseDir    = path.Join(homeDir, "."+mobileqq.PackageName)
	username   string
	password   string
)

const configYAML = `# Go MobileQQ API Configuration Template

accounts:
  - username: 10000
    password: 123456
    status: online

configs:
  auth:
    address: 127.0.0.1:0
    captcha: true
  database:
    dataSourceName: mqqapi.db
    driverName: sqlite
  netIPFamily: dual
  networkType: wifi
  protocol: android-tablet

plugins:
  echo:
    id: 0f73f3cd-edef-4b47-b94c-90bc48953694

servers:
  endpoints:
    - socket://msfwifi.3g.qq.com:8080
    - socket://msfwifiv6.3g.qq.com:8080
    - socket://msfxg.3g.qq.com:8080
    - socket://msfxg.3g.qq.com:80
    - https://msfhttp.3g.qq.com:80
  forceIPv6: false
  overwrite: false
`

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(baseDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			_ = ioutil.WriteFile(
				path.Join(baseDir, "config.yaml"),
				[]byte(configYAML),
				0644,
			)
			log.Fatalf("$_$ [init] create config.yaml")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("x_x [init] failed to load config.yaml")
		}
	} else {
		username = viper.GetString("accounts.0.username")
		password = viper.GetString("accounts.0.password")
	}
}

func main() {
	c := mobileqq.NewClient(
		mobileqq.Option{
			Config: mobileqq.NewClientConfigFromViper(),
		},
	)
	if err := c.HeartbeatAlive(); err != nil {
		log.Printf("x_x [test] error: %s", err.Error())
	}
	if err := c.Auth(username, password); err != nil {
		log.Printf("x_x [auth] error: %s", err.Error())
	}
}
