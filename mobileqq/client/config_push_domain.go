package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"google.golang.org/protobuf/proto"

	"github.com/elap5e/go-mobileqq-api/mobileqq/codec"
	"github.com/elap5e/go-mobileqq-api/pb"
)

func (c *Client) handleConfigPushDomain(
	ctx context.Context,
	s2c *codec.ServerToClientMessage,
) (*codec.ClientToServerMessage, error) {
	data := pb.ConfigPushDomain{}
	if err := proto.Unmarshal(s2c.Buffer, &data); err != nil {
		return nil, err
	}
	tdata, err := json.MarshalIndent(data.GetDomainConfig(), "", "  ")
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path.Join(
		c.GetCacheByUsernameDir(s2c.Username), "domian-config.json",
	), append(tdata, '\n'), 0600); err != nil {
		return nil, err
	}
	return nil, nil
}
