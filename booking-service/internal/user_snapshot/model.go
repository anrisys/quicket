package usersnapshot

import "time"

type UserSnapshot struct {
	ID        uint   `gorm:"primarykey"`
	PublicID  string `gorm:"column:public_id;type:char(36);uniqueIndex"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *UserSnapshot) TableName() string {
	return "users_snapshot"
}