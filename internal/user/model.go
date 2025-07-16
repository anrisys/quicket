package user

import "time"

type User struct {
	ID        uint 		`gorm:"primary_key;column:id;autoIncrement"`
	PublicID  string	`gorm:"column:public_id;type:char(36);uniqueIndex"`
	Email     string	`gorm:"column:email;uniqueIndex"`
	Password  string	`gorm:"column:password;size:255"`
	Role      string	`gorm:"column:role;type:ENUM('user', 'organizer', 'admin');default:'user'"`
	CreatedAt time.Time	`gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time	`gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}