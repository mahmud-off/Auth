package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mahmud-off/auth/internal/domain"
)

type InfoRepository struct {
	db *sqlx.DB
}

func NewInfoRepository(db *sqlx.DB) *InfoRepository {
	return &InfoRepository{db: db}
}

func (r *InfoRepository) GetUsers(ctx *gin.Context) ([]domain.User, error) {

	var users []domain.User
	rows, err := r.db.DB.Query("SELECT id, name, email, registered_at FROM users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
