package model

import "time"

type User struct {
	ID        string    `gorm:"type:text;primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Email     string    `gorm:"size:128;not null;uniqueIndex" json:"email"`
	Age       int       `gorm:"not null" json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
