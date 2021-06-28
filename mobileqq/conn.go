package mobileqq

import (
	"context"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	connSocketMobileIPv4Default = []string{
		"socket://203.205.255.224:443#46000_46002_46007_46001_46003:0:1",
		"socket://msfxg.3g.qq.com:80#46000_46002_46007_46001_46003:0:1",
		"socket://203.205.255.221:443#46000_46002_46007_46001_46003:0:1",
		"socket://183.3.235.162:8080#46000_46002_46007_46001_46003:0:1",
		"socket://163.177.89.195:8080#46000_46002_46007_46001_46003:0:1",
		"socket://183.232.94.44:443#46000_46002_46007_46001_46003:0:1",
	}
	connSocketMobileIPv4CMCC = []string{
		"socket://120.232.18.27:14000#46000_46002_46007:0:1",
		"socket://msfxg.3g.qq.com:80#46000_46002_46007:0:1",
		"socket://111.30.178.75:443#46000_46002_46007:0:1",
		"socket://36.155.240.38:8080#46000_46002_46007:0:1",
		"socket://183.232.94.44:8080#46000_46002_46007:0:1",
		"socket://111.30.138.152:443#46000_46002_46007:0:1",
		"socket://117.144.244.33:443#46000_46002_46007:0:1",
		"socket://111.30.138.152:443#46000_46002_46007:0:1",
	} // 46000, 46002, 46007
	connSocketMobileIPv4CUCC = []string{
		"socket://163.177.89.195:14000#46001:0:1",
		"socket://msfxg.3g.qq.com:80#46001:0:1",
		"socket://157.255.13.77:8080#46001:0:1",
		"socket://221.198.69.96:8080#46001:0:1",
		"socket://153.3.149.61:14000#46001:0:1",
		"socket://111.206.25.142:443#46001:0:1",
		"socket://153.3.50.58:8080#46001:0:1",
	} // 46001
	connSocketMobileIPv4CTCC = []string{
		"socket://113.96.12.224:14000#46003:0:1",
		"socket://msfxg.3g.qq.com:80#46003:0:1",
		"socket://183.3.235.162:443#46003:0:1",
		"socket://42.81.169.100:8080#46003:0:1",
		"socket://114.221.144.89:443#46003:0:1",
		"socket://123.150.76.143:80#46003:0:1",
		"socket://61.129.6.101:14000#46003:0:1",
	} // 46003
	connSocketMobileIPv4Other = []string{
		"socket://msfxg.3g.qq.com:8080#46000_46002_46007_46001_46003:0:1",
		"socket://113.96.12.224:80#46003:0:1",
		"socket://183.232.94.44:14000#46000_46002_46007:0:1",
		"socket://120.232.18.27:8080#46000_46002_46007:0:1",
		"socket://157.255.13.77:443#46001:0:1",
		"socket://203.205.255.224:8080#46000_46002_46007_46001_46003:0:1",
	} // 46000, 46001, 46002, 46003, 46007
	connSocketMobileIPv6Default = []string{
		"socket://msfxgv6.3g.qq.com:8080#00000:0:1",
		"socket://[240e:ff:f101:10::109]:14000",
	}
	connSocketMobileWiFiIPv4Default = []string{
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
	connSocketMobileWiFiIPv6Default = []string{
		"socket://msfwifiv6.3g.qq.com:8080#00000:0:1",
		"socket://[240e:ff:f101:10::109]:14000",
	}
	connHTTPMobileWiFiIPv4Default = []string{
		"https://msfhttp.3g.qq.com:80#00000:0:1",
	}
	connQUICMobileWiFiIPv4Default = []string{
		"quic://58.251.106.174:443#00000:0:1",
	}
)

func (c *Client) createConn(ctx context.Context) (io.ReadWriteCloser, error) {
	rawURIs := append(connSocketMobileWiFiIPv4Default, connSocketMobileWiFiIPv6Default...)

	var addrs []*net.TCPAddr
	for _, rawURI := range rawURIs {
		uri, err := url.Parse(rawURI)
		if err != nil {
			log.Printf("x_x [conn] failed to parse raw uri %s, with error %s", rawURI, err.Error())
			continue
		}
		ips, err := net.LookupIP(uri.Hostname())
		if err != nil {
			log.Printf("x_x [conn] failed to nslookup %s, with error %s", uri.Hostname(), err.Error())
			continue
		}
		port, _ := strconv.Atoi(uri.Port())
		for _, ip := range ips {
			addrs = append(addrs, &net.TCPAddr{IP: ip, Port: port})
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(addrs))
	log.Printf("<== [conn] send tcp dialings")
	for i := range addrs {
		go func(addr *net.TCPAddr) {
			defer wg.Done()
			log.Printf("<-- [conn] send dial tcp %s", addr)
			if err := tcping(addr); err != nil {
				log.Printf("x_x [conn] %s", err.Error())
				return
			}
			log.Printf("--> [conn] recv dial tcp %s", addr)
			c.addrs = append(c.addrs, addr)
		}(addrs[i])
	}
	wg.Wait()
	log.Printf("==> [conn] recv tcp dialings")

	defer log.Printf("^_^ [conn] connected to server %s", c.addrs[0].String())
	return net.Dial("tcp", c.addrs[0].String())
}

func tcping(addr *net.TCPAddr) error {
	for j := 0; j < 5; j++ {
		conn, err := net.DialTimeout("tcp", addr.String(), 5*time.Second)
		if err != nil {
			return err
		}
		conn.Close()
	}
	return nil
}
