package user_entity

import (
	"regexp"

	user_const "github.com/Fi44er/cloud-store-api/internal/modules/user/pkg/constant"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string // identity id
	Username string
	Email    string

	QuotaMax   int64
	AvatarPath string
	BannerPath string
}

func (u *User) Validate() error {
	var err error

	if err = u.validateEmail(); err != nil {
		return err
	}

	return nil
}

func (u *User) validateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9]
		(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

	if u.Email == "" {
		return user_const.ErrEmailRequired
	}

	if len(u.Email) > 254 {
		return user_const.ErrEmailTooLong
	}

	if !emailRegex.MatchString(u.Email) {
		return user_const.ErrEmailInvalidFormat
	}

	return nil
}

func HashPassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
