package internal

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PublicID string `gorm:"column:public_id;type:char(36);uniqueIndex"`
	Email string `gorm:"column:email;uniqueIndex"`
	Password string `gorm:"column:password;size:255"`
	Role string `gorm:"column:role;type:ENUM('user', 'organizer', 'admin');default:'user'"`
}

func (u *User) TableName() string {
	return "users"
}