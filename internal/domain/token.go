package domain

import "time"

type RefreshSession struct {
	ID         int
	UserId     int
	Token      string
	Expires_at time.Time
}
