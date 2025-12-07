package centrifugo

import (
	"context"
	"fmt"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo/apiproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Channel string

const (
	ChannelPublicPolls Channel = Channel("public:polls")
)

type CentrifugoClient struct {
	api  *apiproto.CentrifugoApiClient
	conn *grpc.ClientConn
}

func NewCentrifugoClient(host string, port int) (*CentrifugoClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%v", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := apiproto.NewCentrifugoApiClient(conn)

	return &CentrifugoClient{
		api:  &client,
		conn: conn,
	}, nil
}

func (c *CentrifugoClient) Close() error {
	return c.conn.Close()
}

func (c *CentrifugoClient) Publish(ctx context.Context, channel Channel, data []byte) (*apiproto.PublishResponse, error) {
	resp, err := (*c.api).Publish(ctx, &apiproto.PublishRequest{
		Channel: string(channel),
		Data:    data,
	})
	return resp, err
}

type WebsocketClient interface {
	Publish(ctx context.Context, channel Channel, data []byte) (*apiproto.PublishResponse, error)
	PublishPublicPolls(ctx context.Context, data *PublicPollsData) (*apiproto.PublishResponse, error)
}
