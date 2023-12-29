package models

import (
	"time"
)

type User struct {
	Id       *string `gorm:"primaryKey" json:"id,omitempty"`
	Name     *string `gorm:"name" json:"name,omitempty"`
	Email    *string `gorm:"role" json:"email,omitempty"`
	Password *string `gorm:"password" json:"password,omitempty"`
	Token    *string `gorm:"-" json:"token,omitempty"`
}

type Drug struct {
	Id                  *string    `gorm:"primaryKey" json:"id,omitempty"`
	Name                *string    `gorm:"name" json:"name,omitempty"`
	Approved            *bool      `gorm:"approved" json:"approved,omitempty"`
	MinDose             *int64     `gorm:"min_dose" json:"min_dose,omitempty"`
	MaxDose             *int64     `gorm:"max_dose" json:"max_dose,omitempty"`
	AvailableAtAsString *string    `gorm:"-" json:"available_at,omitempty"`
	AvailableAtAsDate   *time.Time `gorm:"column:available_at" json:"-"`
}
