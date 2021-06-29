package rpc

func (c *Client) HeartbeatAlive() error {
	c2s := &ClientToServerMessage{
		Seq:      c.getNextSeq(),
		Username: "0",
		Buffer:   nil,
		Simple:   false,
	}
	s2c := new(ServerToClientMessage)
	if err := c.Call("Heartbeat.Alive", c2s, s2c); err != nil {
		return err
	}
	return nil
}
