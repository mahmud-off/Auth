package service

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mahmud-off/auth/internal/domain"
	"github.com/mahmud-off/auth/pkg/logger"
	"github.com/spf13/viper"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Repository interface {
	Create(ctx *gin.Context, user domain.User) error
	GetByCredentials(cxt *gin.Context, email string, password string) (domain.User, error)
}

type TokensRepository interface {
	Create(ctx *gin.Context, t domain.RefreshSession) error
	Get(ctx *gin.Context, refreshToken string) (domain.RefreshSession, error)
}

type UsersService struct {
	repo   Repository
	hasher PasswordHasher
	token  TokensRepository

	hmacSecret []byte
}

func NewUsersService(repo Repository, hasher PasswordHasher, tokener TokensRepository, secret []byte) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
		token:  tokener,

		hmacSecret: secret,
	}
}

func (s *UsersService) SignUp(ctx *gin.Context, inp domain.SignUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIn(ctx *gin.Context, inp domain.SignInInput) (string, string, error) {

	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(ctx, int64(user.ID))
}

func (s *UsersService) ParseToken(ctx *gin.Context, token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return s.hmacSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		logger.Debug("Invalid Token")
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}

func (s *UsersService) generateTokens(ctx *gin.Context, userId int64) (string, string, error) {

	AccessTokenExpiresAt := viper.GetInt("access-token.expires_at")
	RefreshTokenExpiresAt := viper.GetInt("refresh-token.expires_at")

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(AccessTokenExpiresAt)).Unix(),
	})

	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.token.Create(ctx, domain.RefreshSession{
		UserId:     int(userId),
		Token:      refreshToken,
		Expires_at: time.Now().Add(time.Hour * time.Duration(RefreshTokenExpiresAt)),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, domain.LEN_OF_REFRESH_TOKEN)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *UsersService) RefreshTokens(ctx *gin.Context, refreshToken string) (string, string, error) {
	session, err := s.token.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.Expires_at.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return s.generateTokens(ctx, int64(session.UserId))
}
