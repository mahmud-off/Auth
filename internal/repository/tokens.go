package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mahmud-off/auth/internal/domain"
)

type Tokens struct {
	db *sqlx.DB
}

func NewTokens(db *sqlx.DB) *Tokens {
	return &Tokens{
		db: db,
	}
}

func (r *Tokens) Create(ctx *gin.Context, t domain.RefreshSession) error {
	//TODO: добавить проверку на существование токена, инчае они начинают копиться в бд, при множественном sign-in'е
	_, err := r.db.DB.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)",
		t.UserId, t.Token, t.Expires_at)
	return err
}

func (r *Tokens) Get(ctx *gin.Context, refreshToken string) (domain.RefreshSession, error) {
	var t domain.RefreshSession
	err := r.db.DB.QueryRow("SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", refreshToken).
		Scan(&t.ID, &t.UserId, &t.Token, &t.Expires_at)
	if err != nil {
		return t, err
	}

	_, err = r.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", t.UserId)

	return t, err
}
