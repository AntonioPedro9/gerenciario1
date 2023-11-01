package models

import "github.com/google/uuid"

type Client struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	CPF    string `gorm:"unique"`
	Name   string `gorm:"not null"`
	Email  string
	Phone  string    `gorm:"not null"`
	UserID uuid.UUID `gorm:"not null"`
}

type CreateClientRequest struct {
	CPF    string    `json:"cpf"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Phone  string    `json:"phone"`
	UserID uuid.UUID `json:"userID"`
}

type UpdateClientRequest struct {
	ID     uint      `json:"id"`
	CPF    string    `json:"cpf"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Phone  string    `json:"phone"`
	UserID uuid.UUID `json:"userID"`
}
