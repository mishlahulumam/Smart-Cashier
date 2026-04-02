package service

import (
	"errors"

	"smart-cashier-backend/internal/model"
	"smart-cashier-backend/internal/repository"
	jwtpkg "smart-cashier-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &model.LoginResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.User, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username sudah digunakan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hash),
		Name:     req.Name,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal membuat user")
	}

	return user, nil
}

func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}
