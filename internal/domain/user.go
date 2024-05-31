package domain

import (
	"time"
)

type User struct {
	Id           int64
	Email        string
	Password     string
	Phone        string
	NickName     string
	Birthday     time.Time
	Introduction string
}
