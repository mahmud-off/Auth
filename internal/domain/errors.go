package domain

import "errors"

var (
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
