package domain

import (
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/gofrs/uuid/v5"
)

type Order struct {
	domain.Base
	PollID      uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex"`
	BuyerUserID uuid.UUID   `gorm:"type:uuid;not null;index"`
	OrderDate   time.Time   `gorm:"type:date;not null;index:,sort:desc"`
	TotalPrice  int64       `gorm:"not null"`
	OrderItems  []OrderItem `gorm:"serializer:json"`
}

type OrderItem struct {
	domain.Base
	FoodID       uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:idx_order_item_order_food"`
	PriceAtOrder int64             `gorm:"not null"`
	Details      []OrderItemDetail `gorm:"serializer:json"`
}

type OrderItemDetail struct {
	UserID   uuid.UUID
	Quantity int
}
