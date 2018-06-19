package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint `gorm:"primary_key"`
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

func HashPwd(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func (u *User) Create(db *gorm.DB) error {
	err := db.Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByName(db *gorm.DB, username string) (*User, error) {
	user := &User{}
	err := db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	user := &User{}
	err := db.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) Save(db *gorm.DB) error {
	err := db.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}
