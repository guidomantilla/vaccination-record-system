package models

type UserCtxKey struct{}

type User struct {
	Id       *string `gorm:"primaryKey" json:"id,omitempty"`
	Name     *string `gorm:"name" json:"name,omitempty"`
	Email    *string `gorm:"role" json:"email,omitempty"`
	Password *string `gorm:"password" json:"password,omitempty"`
	Token    *string `gorm:"-" json:"token,omitempty"`
}
