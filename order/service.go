package order

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/user"
	"strconv"
)

type Service struct {
	orderRepo OrderRepository
}

type OrderRepository interface {
	Store(ctx context.Context, order domain.Order) (int64, error)
	GetByNum(ctx context.Context, num int64) (domain.Order, error)
	FindByUser(ctx context.Context, userId int64) ([]domain.Order, error)
}

func NewOrderService(o OrderRepository) *Service {
	return &Service{
		orderRepo: o,
	}
}

func (s *Service) Add(ctx context.Context, orderNumber []byte) (string, error) {
	orderNum, err := strconv.ParseInt(string(orderNumber), 10, 64)
	if err != nil {
		return "", domain.ErrBadParams
	}
	orderValid := validateOrderNum(orderNum)
	if !orderValid {
		return "", domain.ErrBadOrderNum
	}

	order := domain.Order{
		Number: orderNum,
		UserId: ctx.Value(user.ContextUserIdKey).(int64),
		Status: STATUS_NEW,
	}

	existedOrder, err := s.orderRepo.GetByNum(ctx, orderNum)
	if err != nil {
		internal.Logger.Infow("error in get order", "err", err)
		return "", domain.ErrInternalServerError
	}

	if existedOrder.ID != 0 {
		if existedOrder.UserId == order.UserId {
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
	orders, err := s.orderRepo.FindByUser(ctx, ctx.Value(user.ContextUserIdKey).(int64))
	if err != nil {
		internal.Logger.Infow("error findByUser orders", "err", err)
		return orders, "", domain.ErrInternalServerError
	}

	if len(orders) == 0 {
		return orders, domain.RespNoDataToResponse, nil
	}

	return orders, "", nil
}

func validateOrderNum(orderNum int64) bool {
	err := goluhn.Validate(strconv.FormatInt(orderNum, 10))
	if err != nil {
		return false
	}

	return true
}
