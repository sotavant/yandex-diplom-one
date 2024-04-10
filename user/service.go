package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

const passwordMinLength = 6

type Service struct {
	userRepo UserRepository
}

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
	Store(ctx context.Context, user domain.User) (int64, error)
}

func NewService(u UserRepository) *Service {
	return &Service{
		userRepo: u,
	}
}

func (u *Service) Register(ctx context.Context, user domain.User) (string, error) {
	if user.Login == "" || user.Password == "" {
		return "", domain.ErrBadParams
	}

	user, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		internal.Logger.Infow("error in get by login", "err", err)
		return "", domain.ErrInternalServerError
	}

	if err == nil {
		return "", domain.ErrLoginExist
	}

	if len(user.Password) < passwordMinLength {
		return "", domain.ErrPasswordTooWeak
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		internal.Logger.Infow("error in crypt passwd", "err", err)
		return "", domain.ErrInternalServerError
	}

	userId, err := u.userRepo.Store(ctx, user)
	if err != nil {
		internal.Logger.Infow("error save user", "err", err)
		return "", domain.ErrInternalServerError
	}

	token, err := auth.BuildJWTString(userId)
	if err != nil {
		internal.Logger.Infow("error generation toke", "err", err)
		return "", domain.ErrInternalServerError
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}
