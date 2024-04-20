package accrual

import (
	"context"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/order"
)

type Service struct {
	ordersRepo OrdersRepository
}

type OrdersRepository interface {
	FindByStatus(ctx context.Context, states []string) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, order2 domain.Order) error
	SetAccrual(ctx context.Context, dbOrder domain.Order) error
}

func NewAccrualService(repo OrdersRepository) *Service {
	return &Service{
		ordersRepo: repo,
	}
}

func (s *Service) GetNotProcessedOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := s.ordersRepo.FindByStatus(ctx, order.GetNotProcessedStates())
	if err != nil {
		internal.Logger.Infow("err in find orders by statuses", "err", err)
		return make([]domain.Order, 0), nil
	}

	return orders, nil
}

func (s *Service) UpdateOrderState(ctx context.Context, dbOrder domain.Order, accrual domain.OrderAccrual) error {
	dbOrder.Status = accrual.Status
	dbOrder.Accrual = &accrual.Accrual

	if dbOrder.Number != accrual.Order || accrual.Status == order.StatusRegistered {
		return nil
	}

	if dbOrder.Status != order.StatusProcessed || accrual.Accrual <= 0 {
		err := s.ordersRepo.UpdateStatus(ctx, dbOrder)
		if err != nil {
			internal.Logger.Infow("error in update order status", "err", err)
			return err
		}

		return nil
	}

	err := s.ordersRepo.SetAccrual(ctx, dbOrder)
	if err != nil {
		internal.Logger.Infow("error in set accrual status", "err", err)
	}

	return err
}
