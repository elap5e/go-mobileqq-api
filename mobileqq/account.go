package mobileqq

import (
	"strconv"

	"github.com/elap5e/go-mobileqq-api/rpc"
)

func (c *Client) AccountUpdateStatus(username string, status rpc.AccountStatusType, kick bool) error {
	uin, _ := strconv.ParseInt(username, 10, 64)
	_, err := c.rpc.AccountUpdateStatus(
		c.ctx,
		rpc.NewAccountUpdateStatusRequest(uint64(uin), status, kick),
	)
	if err != nil {
		return err
	}
	return nil
}
