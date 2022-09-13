package domain

import "errors"

var (
	ErrBookNotFound        = errors.New("movie not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrUserAlreadyExists   = errors.New("user with this email already exists")
)
