package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/smazmi/team-task-manager-assignment/backend/internal/models"
	"github.com/smazmi/team-task-manager-assignment/backend/internal/repository"
	"github.com/smazmi/team-task-manager-assignment/backend/pkg/apperror"
	jwtutil "github.com/smazmi/team-task-manager-assignment/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService struct {
	users     repository.UserRepository
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(users repository.UserRepository, jwtSecret string, jwtTTL time.Duration) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*models.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	name := strings.TrimSpace(input.Name)

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil, "", apperror.Conflict("email is already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", apperror.Internal("failed to verify user uniqueness")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", apperror.Internal("failed to hash password")
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, "", apperror.Internal("failed to create user")
	}

	token, err := jwtutil.GenerateToken(s.jwtSecret, s.jwtTTL, user.ID, user.Email)
	if err != nil {
		return nil, "", apperror.Internal("failed to generate token")
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*models.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperror.Unauthorized("invalid email or password")
		}
		return nil, "", apperror.Internal("failed to fetch user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", apperror.Unauthorized("invalid email or password")
	}

	token, err := jwtutil.GenerateToken(s.jwtSecret, s.jwtTTL, user.ID, user.Email)
	if err != nil {
		return nil, "", apperror.Internal("failed to generate token")
	}

	return user, token, nil
}
