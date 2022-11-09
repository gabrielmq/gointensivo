package usecase

import "github.com/gabrielmq/gointensivo/internal/order/entity"

type GetTotalOutputDTO struct {
	Total int
}

func NewGetTotalOutputDTO(total int) *GetTotalOutputDTO {
	return &GetTotalOutputDTO{Total: total}
}

type GetTotalUseCase struct {
	OrderGateway entity.OrderGateway
}

func NewTotalUseCase(orderGategay entity.OrderGateway) *GetTotalUseCase {
	return &GetTotalUseCase{OrderGateway: orderGategay}
}

func (g *GetTotalUseCase) Execute() (*GetTotalOutputDTO, error) {
	total, err := g.OrderGateway.GetTotal()
	if err != nil {
		return nil, err
	}
	return NewGetTotalOutputDTO(total), nil
}
