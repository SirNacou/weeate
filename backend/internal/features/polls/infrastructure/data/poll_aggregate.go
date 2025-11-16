package data

import (
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPollAggregate(db *gorm.DB) gorm.ChainInterface[domain.Poll] {
	pollsChain := gorm.G[domain.Poll](db).
		Joins(clause.JoinTarget{Association: "PollOptions.Votes"}, nil)

	return pollsChain
}
