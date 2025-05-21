package services

import (
	"errors"

	user_models "github.com/nazrawigedion123/wallet-backend/auth/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type AuthService struct {
	db         *gorm.DB
	sessionSvc *SessionService
}

func NewAuthService(db *gorm.DB, sessionSvc *SessionService) *AuthService {
	return &AuthService{
		db:         db,
		sessionSvc: sessionSvc,
	}
}

func (s *AuthService) Login(email, password, ipAddress string) (string, *user_models.User, error) {
	var user user_models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {

		return "", nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		return "", nil, ErrInvalidPassword

	}

	token, err := s.sessionSvc.CreateSession(&user, ipAddress)
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *AuthService) Register(email, password string) (*user_models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := user_models.User{
		Email:    email,
		Password: string(hashedPassword),
		
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
