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

type Vaccination struct {
	Id           *string    `gorm:"primaryKey" json:"id,omitempty"`
	Name         *string    `gorm:"name" json:"name,omitempty"`
	DrugId       *string    `gorm:"column:drugs_id" json:"drug_id,omitempty"`
	Drug         *Drug      `gorm:"foreignKey:DrugId" json:"drug,omitempty"`
	Dose         *int64     `gorm:"dose" json:"dose,omitempty"`
	DateAsString *string    `gorm:"-" json:"date,omitempty"`
	DateAsDate   *time.Time `gorm:"column:date" json:"-"`
}
