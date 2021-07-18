package rpc

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elap5e/go-mobileqq-api/log"
)

var (
	serverListSocket = []string{
		"socket://203.205.255.224:443#46000_46002_46007_46001_46003:0:1",
		"socket://msfxg.3g.qq.com:80#46000_46002_46007_46001_46003:0:1",
		"socket://203.205.255.221:443#46000_46002_46007_46001_46003:0:1",
		"socket://183.3.235.162:8080#46000_46002_46007_46001_46003:0:1",
		"socket://163.177.89.195:8080#46000_46002_46007_46001_46003:0:1",
		"socket://183.232.94.44:443#46000_46002_46007_46001_46003:0:1",
	}
	serverListSocketCMCC = []string{
		"socket://120.232.18.27:14000#46000_46002_46007:0:1",
		"socket://msfxg.3g.qq.com:80#46000_46002_46007:0:1",
		"socket://111.30.178.75:443#46000_46002_46007:0:1",
		"socket://36.155.240.38:8080#46000_46002_46007:0:1",
		"socket://183.232.94.44:8080#46000_46002_46007:0:1",
		"socket://111.30.138.152:443#46000_46002_46007:0:1",
		"socket://117.144.244.33:443#46000_46002_46007:0:1",
		"socket://111.30.138.152:443#46000_46002_46007:0:1",
	} // 46000, 46002, 46007
	serverListSocketCUCC = []string{
		"socket://163.177.89.195:14000#46001:0:1",
		"socket://msfxg.3g.qq.com:80#46001:0:1",
		"socket://157.255.13.77:8080#46001:0:1",
		"socket://221.198.69.96:8080#46001:0:1",
		"socket://153.3.149.61:14000#46001:0:1",
		"socket://111.206.25.142:443#46001:0:1",
		"socket://153.3.50.58:8080#46001:0:1",
	} // 46001
	serverListSocketCTCC = []string{
		"socket://113.96.12.224:14000#46003:0:1",
		"socket://msfxg.3g.qq.com:80#46003:0:1",
		"socket://183.3.235.162:443#46003:0:1",
		"socket://42.81.169.100:8080#46003:0:1",
		"socket://114.221.144.89:443#46003:0:1",
		"socket://123.150.76.143:80#46003:0:1",
		"socket://61.129.6.101:14000#46003:0:1",
	} // 46003
	serverListSocketOther = []string{
		"socket://msfxg.3g.qq.com:8080#46000_46002_46007_46001_46003:0:1",
		"socket://113.96.12.224:80#46003:0:1",
		"socket://183.232.94.44:14000#46000_46002_46007:0:1",
		"socket://120.232.18.27:8080#46000_46002_46007:0:1",
		"socket://157.255.13.77:443#46001:0:1",
		"socket://203.205.255.224:8080#46000_46002_46007_46001_46003:0:1",
	} // 46000, 46001, 46002, 46003, 46007
	serverListSocketIPv6 = []string{
		"socket://msfxgv6.3g.qq.com:8080#00000:0:1",
		"socket://[240e:ff:f101:10::109]:14000",
	}
	serverListSocketWiFi = []string{
		"socket://msfwifi.3g.qq.com:8080#00000:0:1",
		"socket://14.215.138.110:8080#00000:0:1",
		"socket://113.96.12.224:8080#00000:0:1",
		"socket://157.255.13.77:14000#00000:0:1",
		"socket://120.232.18.27:443#00000:0:1",
		"socket://183.3.235.162:14000#00000:0:1",
		"socket://163.177.89.195:443#00000:0:1",
		"socket://183.232.94.44:80#00000:0:1",
		"socket://203.205.255.224:8080#00000:0:1",
		"socket://203.205.255.221:8080#00000:0:1",
	}
	serverListSocketWiFiIPv6 = []string{
		"socket://msfwifiv6.3g.qq.com:8080#00000:0:1",
		"socket://[240e:ff:f101:10::109]:14000",
	}
	serverListHTTP = []string{
		"https://msfhttp.3g.qq.com:80#00000:0:1",
	}
	serverListQUIC = []string{
		"quic://58.251.106.174:443#00000:0:1",
	}
)

func getServerList(network string) []string {
	switch strings.ToLower(network) {
	case "cmcc":
		return append(append(
			serverListSocketCMCC,
			serverListSocketOther...),
			serverListSocketIPv6...,
		)
	case "cmccv4":
		return append(
			serverListSocketCMCC,
			serverListSocketOther...,
		)
	case "ctcc":
		return append(append(
			serverListSocketCTCC,
			serverListSocketOther...),
			serverListSocketIPv6...,
		)
	case "ctccv4":
		return append(
			serverListSocketCTCC,
			serverListSocketOther...,
		)
	case "cucc":
		return append(append(
			serverListSocketCUCC,
			serverListSocketOther...),
			serverListSocketIPv6...,
		)
	case "cuccv4":
		return append(
			serverListSocketCUCC,
			serverListSocketOther...,
		)
	case "cxccv6":
		return serverListSocketIPv6
	case "http":
		return serverListHTTP
	case "quic":
		return serverListQUIC
	case "wifi":
		return append(
			serverListSocketWiFi,
			serverListSocketWiFiIPv6...,
		)
	case "wifiv4":
		return serverListSocketWiFi
	case "wifiv6":
		return serverListSocketWiFiIPv6
	default:
		return append(append(append(append(
			serverListSocket,
			serverListSocketOther...),
			serverListSocketWiFi...),
			serverListSocketIPv6...),
			serverListSocketWiFiIPv6...,
		)
	}
}

func dialing(addr *net.TCPAddr) error {
	for j := 0; j < 5; j++ {
		conn, err := net.DialTimeout("tcp", addr.String(), 5*time.Second)
		if err != nil {
			return err
		}
		conn.Close()
	}
	return nil
}

func (e *engine) SetServers(list []string) {
	e.tcpTesting(list)
}

func (e *engine) tcpTesting(list []string) {
	var addrs []*net.TCPAddr
	for _, item := range list {
		uri, err := url.Parse(item)
		if err != nil {
			log.Warn().Err(err).
				Msgf("x—x [conn] failed to parse raw uri %s", item)
			continue
		}
		ips, err := net.LookupIP(uri.Hostname())
		if err != nil {
			log.Warn().Err(err).
				Msgf("x—x [conn] failed to nslookup %s", uri.Hostname())
			continue
		}
		port, _ := strconv.Atoi(uri.Port())
		for _, ip := range ips {
			skip := false
			for _, addr := range addrs {
				if addr.IP.String() == ip.String() && addr.Port == port {
					skip = true
				}
			}
			if !skip {
				addrs = append(addrs, &net.TCPAddr{IP: ip, Port: port})
			}
		}
	}
	for _, item := range e.addrs {
		skip := false
		for _, addr := range addrs {
			if addr.String() == item {
				skip = true
			}
		}
		if !skip {
			tcpAddr, err := net.ResolveTCPAddr("tcp", item)
			if err != nil {
				log.Warn().Err(err).
					Msgf("x—x [conn] failed to parse raw uri %s", item)
				continue
			}
			addrs = append(addrs, tcpAddr)
		}
	}

	e.addrs = []string{}
	wg := sync.WaitGroup{}
	wg.Add(len(addrs))
	log.Info().Msg("<-- [conn] testing tcp connections...")
	for i := range addrs {
		go func(addr *net.TCPAddr) {
			defer wg.Done()
			log.Debug().
				Msg("··· [conn] dialing tcp " + addr.String())
			if err := dialing(addr); err != nil {
				log.Warn().Err(err).
					Msg("x—x [conn] failed to dial tcp connection")
			} else {
				e.addrs = append(e.addrs, addr.String())
			}
		}(addrs[i])
	}
	wg.Wait()
	if len(e.addrs) > 10 {
		e.addrs = e.addrs[:10]
	}

	if len(e.addrs) != 0 {
		log.Info().Msg("--> [conn] tcp connections tested")
	} else {
		log.Error().Msg("x—x [conn] unavailable tcp connections")
		e.Close()
	}
}
