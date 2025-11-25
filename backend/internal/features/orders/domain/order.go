package domain

import (
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	polls_domain "github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
)

type Order struct {
	domain.Base
	PollID      uuid.UUID                 `gorm:"type:uuid;not null;uniqueIndex"`
	BuyerUserID string                    `gorm:"not null;uniqueIndex:idx_order_buyer_date"`
	OrderDate   time.Time                 `gorm:"type:date;not null;uniqueIndex:idx_order_buyer_date,sort:desc"`
	Strategy    polls_domain.PollStrategy `gorm:"not null"`
	TotalPrice  int64                     `gorm:"not null"`
	OrderItems  []OrderItem               `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	OrderID      uuid.UUID         `gorm:"type:uuid;not null;primaryKey"`
	FoodID       uuid.UUID         `gorm:"type:uuid;not null;primaryKey"`
	PriceAtOrder int64             `gorm:"not null"`
	Details      []OrderItemDetail `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE"`
}

type OrderItemDetail struct {
	OrderItemID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	UserID      string    `gorm:"not null;primaryKey"`
	Quantity    int64     `gorm:"not null"`
}

func NewOrder(pollID uuid.UUID, buyerUserID string, orderDate time.Time, strategy polls_domain.PollStrategy, orderItems []OrderItem) (*Order, error) {
	totalPrice := lo.SumBy(orderItems, func(item OrderItem) int64 {
		return item.PriceAtOrder * lo.SumBy(item.Details, func(detail OrderItemDetail) int64 {
			return detail.Quantity
		})
	})

	order := &Order{
		Base:        domain.NewBase(),
		PollID:      pollID,
		BuyerUserID: buyerUserID,
		OrderDate:   orderDate,
		Strategy:    strategy,
		TotalPrice:  totalPrice,
		OrderItems:  orderItems,
	}

	return order, nil
}

func NewOrderItem(foodID uuid.UUID, priceAtOrder int64, details []OrderItemDetail) *OrderItem {
	return &OrderItem{
		FoodID:       foodID,
		PriceAtOrder: priceAtOrder,
		Details:      details,
	}
}

func NewOrderItemDetail(userID string, quantity int64) *OrderItemDetail {
	return &OrderItemDetail{
		UserID:   userID,
		Quantity: quantity,
	}
}
