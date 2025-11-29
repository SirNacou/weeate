package centrifugo

import (
	"context"
	"encoding/json"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo/apiproto"
)

type PublicPollsDataType string

const (
	PollCreated PublicPollsDataType = PublicPollsDataType("poll_created")
	PollClosed  PublicPollsDataType = PublicPollsDataType("poll_closed")
	VoteAdded   PublicPollsDataType = PublicPollsDataType("vote_added")
	VoteMoved   PublicPollsDataType = PublicPollsDataType("vote_moved")
	VoteRemoved PublicPollsDataType = PublicPollsDataType("vote_removed")
)

type PublicPollsData struct {
	Type PublicPollsDataType `json:"type"`
	Data any                 `json:"data"`
}

func (c *CentrifugoClient) PublishPublicPolls(ctx context.Context, data *PublicPollsData) (*apiproto.PublishResponse, error) {
	// Serialize data to JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Publish(context.Background(), ChannelPublicPolls, payload)
}
