package auth

import (
	"errors"
	"pLaNtS/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	hasher                 PasswordHasher
	jwtService             TokenService
	userRepository         UserRepository
	refreshTokenRepository RefreshTokenRepository
	refreshTokenTTL        time.Duration
}

func NewAuthService(hasher PasswordHasher, jwtService TokenService, refreshTokenRepository RefreshTokenRepository, userRepository UserRepository, refreshTokenTTL time.Duration) *Service {
	return &Service{
		hasher:                 hasher,
		jwtService:             jwtService,
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
		refreshTokenTTL:        refreshTokenTTL,
	}
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) bool
}

type TokenService interface {
	Generate(userID string, email string) (*domain.TokenPair, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(token *domain.RefreshToken) error
	FindRefreshTokenByToken(token string) (*domain.RefreshToken, error)
	DeleteRefreshTokenById(id uuid.UUID) error
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserById(id uuid.UUID) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

func (s *Service) Register(email, password, firstName, lastName string) (*domain.TokenResponse, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	err = s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtService.Generate(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken := s.newRefreshToken(user.ID, tokens.RefreshToken)

	err = s.refreshTokenRepository.CreateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *Service) Login(email, password string) (*domain.TokenResponse, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !s.hasher.Compare(user.PasswordHash, password) {
		return nil, errors.New("wrong password")
	}

	tokens, err := s.jwtService.Generate(user.ID.String(), user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken := s.newRefreshToken(user.ID, tokens.RefreshToken)

	err = s.refreshTokenRepository.CreateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *Service) Refresh(token string) (*domain.TokenResponse, error) {
	refreshToken, err := s.refreshTokenRepository.FindRefreshTokenByToken(token)
	if err != nil {
		return nil, err
	}

	if refreshToken.IsRevoked {
		return nil, errors.New("refresh token is revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("refresh token is expired")
	}

	userID := refreshToken.UserID
	user, err := s.userRepository.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	email := user.Email

	tokens, err := s.jwtService.Generate(userID.String(), email)
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenRepository.DeleteRefreshTokenById(refreshToken.ID)
	if err != nil {
		return nil, err
	}

	newRefreshToken := s.newRefreshToken(user.ID, tokens.RefreshToken)

	err = s.refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}, nil
}

func (s *Service) GetMe(userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepository.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (s *Service) newRefreshToken(userID uuid.UUID, token string) *domain.RefreshToken {
	now := time.Now().UTC()
	return &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(s.refreshTokenTTL),
		CreatedAt: now,
		IsRevoked: false,
	}
}
