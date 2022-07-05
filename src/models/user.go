package models

import "time"

type User struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  *[]byte   `json:"-"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt"`
}

func (User) TableName() string {
	return "USERS"
}
