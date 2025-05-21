package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	user_models "github.com/nazrawigedion123/wallet-backend/auth/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisUnavailable = errors.New("redis service unavailable")
	ErrInvalidToken     = errors.New("invalid token")
)

type SessionService struct {
	redisClient *redis.Client
	secretKey   []byte
	sessionTTL  time.Duration
}

func NewSessionService(redisClient *redis.Client, secretKey string, ttl time.Duration) *SessionService {
	return &SessionService{
		redisClient: redisClient,
		secretKey:   []byte(secretKey),
		sessionTTL:  ttl,
	}
}

func (s *SessionService) CreateSession(user *user_models.User, ipAddress string) (string, error) {
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.sessionTTL).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	// Convert struct to a map for Redis HSET
	metadata := map[string]interface{}{
		"UserID":    user.ID,
		"Email":     user.Email,
		"Tier":      user.Tier,
		"LastLogin": time.Now().Format(time.RFC3339), // store as string
		"IPAddress": ipAddress,
	}

	ctx := context.Background()
	err = s.redisClient.HSet(ctx, "session:"+tokenString, metadata).Err()
	if err != nil {
		return "", err
	}

	// Set TTL for the session
	err = s.redisClient.Expire(ctx, "session:"+tokenString, s.sessionTTL).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *SessionService) ValidateSession(tokenString string) (*user_models.SessionMetadata, error) {
	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check Redis for session data
	ctx := context.Background()
	//metadata := &user_models.SessionMetadata{}

	result, err := s.redisClient.HGetAll(ctx, "session:"+tokenString).Result()
	if err != nil {

		if err == redis.Nil {
			return nil, ErrInvalidToken
		}
		return nil, ErrRedisUnavailable
	}

	if len(result) == 0 {
		return nil, ErrInvalidToken
	}

	metadata := &user_models.SessionMetadata{}
	if uid, ok := result["UserID"]; ok {
		parsedID, _ := strconv.ParseUint(uid, 10, 64)
		metadata.UserID = uint(parsedID)
	}
	metadata.Email = result["Email"]
	metadata.Tier = result["Tier"]
	metadata.IPAddress = result["IPAddress"]
	if lastLoginStr, ok := result["LastLogin"]; ok {
		t, _ := time.Parse(time.RFC3339, lastLoginStr)
		metadata.LastLogin = t
	}

	// Verify token matches stored session
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || uint(claims["sub"].(float64)) != metadata.UserID {
		return nil, ErrInvalidToken
	}

	return metadata, nil
}

func (s *SessionService) InvalidateSession(tokenString string) error {
	ctx := context.Background()
	_, err := s.redisClient.Del(ctx, "session:"+tokenString).Result()
	return err
}
