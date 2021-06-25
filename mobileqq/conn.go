package mobileqq

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func (c *Client) createConn(ctx context.Context) (io.ReadWriteCloser, error) {
	var addrs []*net.TCPAddr

	ips, err := net.LookupIP("msfwifi.3g.qq.com")
	if err != nil {
		log.Printf("x_x [conn] failed to nslookup msfwifi.3g.qq.com, with error %s", err.Error())
	}
	for _, ip := range ips {
		addrs = append(addrs, &net.TCPAddr{IP: ip, Port: 8080})
	}
	ips, err = net.LookupIP("msfwifiv6.3g.qq.com")
	if err != nil {
		log.Printf("x_x [conn] failed to nslookup msfwifiv6.3g.qq.com, with error %s", err.Error())
	}
	for _, ip := range ips {
		addrs = append(addrs, &net.TCPAddr{IP: ip, Port: 8080})
	}

	addrs = append(addrs, []*net.TCPAddr{
		{IP: net.IPv4(14, 215, 138, 110), Port: 8080},
		{IP: net.IPv4(113, 96, 12, 224), Port: 8080},
		{IP: net.IPv4(157, 255, 13, 77), Port: 14000},
		{IP: net.IPv4(120, 232, 18, 27), Port: 443},
		{IP: net.IPv4(183, 3, 235, 162), Port: 14000},
		{IP: net.IPv4(163, 177, 89, 195), Port: 443},
		{IP: net.IPv4(183, 232, 94, 44), Port: 80},
		{IP: net.IPv4(203, 205, 255, 224), Port: 8080},
		{IP: net.IPv4(203, 205, 255, 221), Port: 8080},
		{IP: net.IP{0x24, 0x0e, 0x00, 0xff, 0xf1, 0x01, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x09}, Port: 14000},
	}...)

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
