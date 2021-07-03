package rpc

func (c *Client) HeartbeatAlive() error {
	c2s := &ClientToServerMessage{
		Seq:      c.getNextSeq(),
		Username: "0",
		Buffer:   nil,
		Simple:   false,
	}
	s2c := ServerToClientMessage{}
	if err := c.Call("Heartbeat.Alive", c2s, &s2c); err != nil {
		return err
	}
	return nil
}
