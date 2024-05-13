package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

const passwordMinLength = 6

type ContextUserIDKey struct{}

type Service struct {
	userRepo UserRepository
}

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
	Store(ctx context.Context, user domain.User) (int64, error)
	GetByID(ctx context.Context, userID int64) (domain.User, error)
}

func NewService(u UserRepository) *Service {
	return &Service{
		userRepo: u,
	}
}

func (u *Service) GetByID(ctx context.Context, userID int64) (domain.User, error) {
	dbUser, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		internal.Logger.Infow("error in get by id", "err", err)
		return domain.User{}, domain.ErrInternalServerError
	}

	if dbUser.ID == 0 {
		return domain.User{}, domain.ErrBadUserData
	}

	return dbUser, nil
}

func (u *Service) Register(ctx context.Context, user domain.User) (string, error) {
	dbUser, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		internal.Logger.Infow("error in get by login", "err", err)
		return "", domain.ErrInternalServerError
	}

	if (err != nil && errors.Is(err, pgx.ErrNoRows)) || (err == nil && dbUser.ID != 0) {
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

	userID, err := u.userRepo.Store(ctx, user)
	if err != nil {
		internal.Logger.Infow("error save user", "err", err)
		return "", domain.ErrInternalServerError
	}

	token, err := auth.BuildJWTString(userID)
	if err != nil {
		internal.Logger.Infow("error generation token", "err", err)
		return "", domain.ErrInternalServerError
	}

	return token, nil
}

func (u *Service) Login(ctx context.Context, user domain.User) (string, error) {
	dbUser, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		internal.Logger.Infow("error in get by login", "err", err)
		return "", domain.ErrInternalServerError
	}

	if dbUser.ID == 0 {
		return "", domain.ErrBadUserData
	}

	fmt.Println(user.Password)
	passwordCorrect, err := checkPassword(user.Password, dbUser.Password)
	if err != nil {
		internal.Logger.Infow("error in check passwd", "err", err)
		return "", domain.ErrInternalServerError
	}

	if !passwordCorrect {
		return "", domain.ErrBadUserData
	}

	token, err := auth.BuildJWTString(dbUser.ID)
	if err != nil {
		internal.Logger.Infow("error generation token", "err", err)
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

func checkPassword(password, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, err
}
