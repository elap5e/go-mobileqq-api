package tcp

import (
	"net"
	"time"

	"github.com/elap5e/go-mobileqq-api/bytes"
	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
)

var bufPool = bytes.NewPool(0)

type clientCodec struct {
	conn net.Conn

	bufResp *bytes.Buffer
}

func NewClientCodec(conn net.Conn) codec.ClientCodec {
	return &clientCodec{conn: conn}
}

func (c *clientCodec) read() ([]byte, error) {
	p := make([]byte, 4)
	if _, err := c.conn.Read(p); err != nil {
		return p, err
	}
	l := int(p[0])<<24 | int(p[1])<<16 | int(p[2])<<8 | int(p[3])<<0
	p = append(p, make([]byte, l-4)...)
	i := 4
	for i < l {
		n, err := c.conn.Read(p[i:])
		if err != nil {
			return p, err
		}
		i += n
	}
	return p, nil
}

func (c *clientCodec) Close() error {
	return c.conn.Close()
}

func (c *clientCodec) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *clientCodec) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *clientCodec) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
