package withdrawn

import (
	"context"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/order"
	"github.com/sotavant/yandex-diplom-one/user"
)

type Service struct {
	wdRepo   WithdrawnRepository
	userRepo user.UserRepository
}

type WithdrawnRepository interface {
	Store(ctx context.Context, withdrawn domain.Withdrawn) error
}

func NewService(wd WithdrawnRepository, ur user.UserRepository) *Service {
	return &Service{
		wdRepo:   wd,
		userRepo: ur,
	}
}

func (s *Service) Add(ctx context.Context, wd *domain.Withdrawn) error {
	valid := order.ValidateOrderNum(wd.OrderNum)
	if !valid {
		return domain.ErrBadOrderNum
	}

	bdUser, err := s.userRepo.GetById(ctx, wd.UserId)
	if err != nil {
		return domain.ErrInternalServerError
	}

	if bdUser.Current < wd.Sum {
		return domain.ErrNotEnoughCurrent
	}

	err = s.wdRepo.Store(ctx, *wd)
	if err != nil {
		return domain.ErrInternalServerError
	}

	return nil
}
