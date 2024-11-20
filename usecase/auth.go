package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mirasildev/chat_task/config"
	"github.com/mirasildev/chat_task/domain"
	"github.com/mirasildev/chat_task/internal/repository/inmemory"
	emailPkg "github.com/mirasildev/chat_task/pkg/email"
	"github.com/mirasildev/chat_task/pkg/utils"
)

type AuthService struct {
	userRepo UserRepository
	cfg      config.Config
	inMemory inmemory.InMemoryStorageI
}

func NewAuthService(g UserRepository) *AuthService {
	return &AuthService{
		userRepo: g,
	}
}

const (
	RegisterCodeKey = "register_code_"
)

func (s *AuthService) Register(req *domain.RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = s.inMemory.Set("user_"+user.Email, string(userData), 10*time.Minute)
	if err != nil {
		return err
	}

	go func() {
		err := s.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	return nil
}

func (s *AuthService) sendVerificationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = s.inMemory.Set(key+email, code, time.Minute)
	if err != nil {
		return err
	}

	err = emailPkg.SendEmail(&s.cfg, &emailPkg.SendEmailRequest{
		To:   []string{email},
		Type: "verification_email",
		Body: map[string]string{
			"code": code,
		},
		Subject: "Verification email",
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Verify(req *domain.VerifyRequest) (*domain.AuthResponse, error) {
	userData, err := s.inMemory.Get("user_" + req.Email)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, err
	}

	code, err := s.inMemory.Get(RegisterCodeKey + user.Email)
	if err != nil {
		return nil, fmt.Errorf("code_expired")
	}

	if req.Code != code {
		return nil, fmt.Errorf("incorrect_code")
	}

	result, err := s.userRepo.Store(&user)
	if err != nil {
		return nil, err
	}

	token, _, err := utils.CreateToken(&s.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Id:          result.ID,
		Email:       result.Email,
		Username:    result.Username,
		CreatedAt:   result.CreatedAt.Format(time.RFC3339),
		AccessToken: token,
	}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *domain.VerifyTokenRequest) (*domain.AuthPayload, error) {
	accessToken := req.AccessToken

	payload, err := utils.VerifyToken(&s.cfg, accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return &domain.AuthPayload{
		Id:        payload.ID.String(),
		UserID:    payload.UserID,
		Email:     payload.Email,
		Username:  payload.Username,
		IssuedAt:  payload.IssuedAt.Format(time.RFC3339),
		ExpiredAt: payload.ExpiredAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		log.Printf("failed to get user by email: %v", err)
		return nil, err
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("incorrect_password")
	}

	token, _, err := utils.CreateToken(&s.cfg, &utils.TokenParams{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Duration: time.Hour * 24,
	})
	if err != nil {
		log.Printf("failed to create token: %v", err)
		return nil, err
	}

	return &domain.AuthResponse{
		Id:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		AccessToken: token,
	}, nil
}
