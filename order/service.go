package order

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/user"
)

type Service struct {
	orderRepo OrderRepository
}

type OrderRepository interface {
	Store(ctx context.Context, order domain.Order) (int64, error)
	GetByNum(ctx context.Context, num string) (domain.Order, error)
	FindByUser(ctx context.Context, userID int64) ([]domain.Order, error)
}

func NewOrderService(o OrderRepository) *Service {
	return &Service{
		orderRepo: o,
	}
}

func (s *Service) Add(ctx context.Context, orderNumber []byte) (string, error) {
	orderNum := string(orderNumber)
	orderValid := ValidateOrderNum(orderNum)
	if !orderValid {
		return "", domain.ErrBadOrderNum
	}

	order := domain.Order{
		Number: orderNum,
		UserID: ctx.Value(user.ContextUserIDKey{}).(int64),
		Status: StatusNew,
	}

	existedOrder, err := s.orderRepo.GetByNum(ctx, orderNum)
	if err != nil {
		internal.Logger.Infow("error in get order", "err", err)
		return "", domain.ErrInternalServerError
	}

	if existedOrder.ID != 0 {
		if existedOrder.UserID == order.UserID {
			return domain.RespOrderAlreadyUploaded, nil
		} else {
			return "", domain.ErrOrderAlreadyUploaded
		}
	}

	_, err = s.orderRepo.Store(ctx, order)
	if err != nil {
		internal.Logger.Infow("error save order", "err", err)
		return "", domain.ErrInternalServerError
	}

	return domain.RespOrderAdmitted, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Order, string, error) {
	orders, err := s.orderRepo.FindByUser(ctx, ctx.Value(user.ContextUserIDKey{}).(int64))
	if err != nil {
		internal.Logger.Infow("error findByUser orders", "err", err)
		return orders, "", domain.ErrInternalServerError
	}

	if len(orders) == 0 {
		return orders, domain.RespNoDataToResponse, nil
	}

	return orders, "", nil
}

func ValidateOrderNum(orderNum string) bool {
	return goluhn.Validate(orderNum) == nil
}
