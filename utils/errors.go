package utils

import "errors"

const (
	ErrWrongSecretKey         = "Wrong secret key"
)

var (
	WrongSecretKey            = errors.New("Wrong secret key")
)

func NewWrongSecretKeyError() error {
	return WrongSecretKey
}