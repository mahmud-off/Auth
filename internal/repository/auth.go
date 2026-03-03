package repository

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mahmud-off/auth/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UsersRepository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewUsersRepository(db *sqlx.DB, rdb *redis.Client) *UsersRepository {
	return &UsersRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *UsersRepository) Create(ctx *gin.Context, user domain.User) error {
	_, err := r.db.DB.Exec("INSERT INTO users (name, email, password, registered_at) values ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)
	return err
}

func (r *UsersRepository) GetByCredentials(cxt *gin.Context, email string, password string) (domain.User, error) {
	var user domain.User
	err := r.db.DB.QueryRow("SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}

func (r *UsersRepository) AddToBlackList(ctx *gin.Context, accessToken string) error {
	_, err := r.rdb.Set(ctx, accessToken, "blocked", time.Duration(time.Minute*15)).Result()

	return err
}

func (r *UsersRepository) TokenInBlackList(ctx *gin.Context, token string) bool {
	val := r.rdb.Exists(ctx, token).Val()
	if val == 1 {
		return true
	}
	return false
}
