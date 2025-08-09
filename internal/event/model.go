package event

import (
	"time"

	"github.com/anrisys/quicket/internal/user"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	PublicID 		string 		`gorm:"column:public_id;type:char(36);uniqueIndex"`
	Title 			string 		`gorm:"column:title;size:256;not null"`
	Description 	*string 	`gorm:"type:text"`
	StartDate 		time.Time 	`gorm:"column:start_date;not null;index"`
	EndDate 		time.Time 	`gorm:"column:end_date;not null"`
	MaxSeats 		uint64 		`gorm:"column:max_seats"`
	AvailableSeats 	uint64 		`gorm:"column:available_seats"`
	OrganizerID 	uint 		`gorm:"column:organizer_id;not null"`
	Organizer 		user.User	`gorm:"foreignKey:OrganizerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *Event) TableName() string {
	return "events"
}