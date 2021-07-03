package rpc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) handlePushConfigDomain(
	ctx context.Context,
	s2c *ServerToClientMessage,
) (*ClientToServerMessage, error) {
	data := pb.ConfigDomain{}
	if err := proto.Unmarshal(s2c.Buffer, &data); err != nil {
		return nil, err
	}
	tdata, err := json.MarshalIndent(data.GetDomainList(), "", "  ")
	if err != nil {
		return nil, err
	}
	if ioutil.WriteFile(path.Join(
		c.cfg.CacheDir, s2c.Username, "domian-list.json",
	), append(tdata, '\n'), 0600); err != nil {
		return nil, err
	}
	return nil, nil
}
