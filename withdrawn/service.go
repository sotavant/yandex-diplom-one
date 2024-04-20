package withdrawn

import (
	"context"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/order"
	"github.com/sotavant/yandex-diplom-one/user"
)

type Service struct {
	wdRepo   WithdrawnRepository
	userRepo user.UserRepository
}

type WithdrawnRepository interface {
	Store(ctx context.Context, withdrawn domain.Withdrawn) error
	FindOne(ctx context.Context, orderNum string) (domain.Withdrawn, error)
	FindByUser(ctx context.Context, userID int64) ([]domain.Withdrawn, error)
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

	dbWd, err := s.wdRepo.FindOne(ctx, wd.OrderNum)
	if err != nil {
		return domain.ErrInternalServerError
	}
	if dbWd.ID > 0 {
		return domain.ErrOrderAlreadyUploaded
	}

	bdUser, err := s.userRepo.GetByID(ctx, wd.UserID)
	if err != nil {
		return domain.ErrUserNotAuthorized
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

func (s *Service) List(ctx context.Context, userID int64) ([]domain.Withdrawn, string, error) {
	wds, err := s.wdRepo.FindByUser(ctx, userID)
	if err != nil {
		internal.Logger.Infow("error findByUser wds", "err", err)
		return wds, "", domain.ErrInternalServerError
	}

	if len(wds) == 0 {
		return wds, domain.RespNoDataToResponse, nil
	}

	return wds, "", nil
}
