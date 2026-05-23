package service

import (
	"context"
	"errors"

	"github.com/example/repair-crm/internal/models"
	"github.com/example/repair-crm/internal/repository"
	"github.com/example/repair-crm/pkg/auth"
	"github.com/example/repair-crm/pkg/password"
)

type AuthService struct {
	masters *repository.MasterRepository
	jwt     *auth.JWTManager
}

type LoginResult struct {
	Token  string        `json:"token"`
	Master models.Master `json:"master"`
}

func NewAuthService(masters *repository.MasterRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{masters: masters, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, username, plainPassword string) (*LoginResult, error) {
	if username == "" || plainPassword == "" {
		return nil, ValidationError{Message: "укажите логин и пароль"}
	}

	master, err := s.masters.FindByUsername(ctx, username)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ValidationError{Message: "неверный логин или пароль"}
	}
	if err != nil {
		return nil, err
	}

	if err := password.Compare(master.PasswordHash, plainPassword); err != nil {
		return nil, ValidationError{Message: "неверный логин или пароль"}
	}

	token, err := s.jwt.Generate(master.ID, master.WorkshopID, master.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, Master: *master}, nil
}
