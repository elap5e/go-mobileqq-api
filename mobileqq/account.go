package mobileqq

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/elap5e/go-mobileqq-api/rpc"
)

func (c *Client) AccountUpdateStatus(username string, status rpc.AccountStatusType, kick bool) error {
	uin, _ := strconv.ParseInt(username, 10, 64)
	resp, err := c.rpc.AccountUpdateStatus(
		c.ctx,
		rpc.NewAccountUpdateStatusRequest(uint64(uin), status, kick),
	)
	if err != nil {
		return err
	}
	jresp, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("AccountUpdateStatusResponse\n%s", jresp)
	return nil
}
