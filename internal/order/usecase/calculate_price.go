package usecase

import (
	"github.com/gabrielmq/gointensivo/internal/order/entity"
)

type OrderInputDTO struct {
	ID    string
	Price float64
	Tax   float64
}

type OrderOutputDTO struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrderOutputDTOFrom(order *entity.Order) *OrderOutputDTO {
	return &OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
}

type CalculateFinalPriceUseCase struct {
	OrderGateway entity.OrderGateway
}

func NewCalculateFinalPriceUseCase(orderGateway entity.OrderGateway) *CalculateFinalPriceUseCase {
	return &CalculateFinalPriceUseCase{
		OrderGateway: orderGateway,
	}
}

func (c *CalculateFinalPriceUseCase) Execute(input OrderInputDTO) (*OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	if err = order.CalculateFinalPrice(); err != nil {
		return nil, err
	}

	if err = c.OrderGateway.Save(order); err != nil {
		return nil, err
	}
	return NewOrderOutputDTOFrom(order), nil
}
