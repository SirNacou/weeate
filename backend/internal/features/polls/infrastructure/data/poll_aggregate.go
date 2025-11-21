package data

import (
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"gorm.io/gorm"
)

func GetPollAggregate(db *gorm.DB) gorm.ChainInterface[domain.Poll] {
	pollsChain := gorm.G[domain.Poll](db).
		Preload("PollOptions.Votes", nil)

	return pollsChain
}
