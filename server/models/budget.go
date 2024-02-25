package models

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID             uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Price          float32         `json:"price"`
	Vehicle        string          `gorm:"not null" json:"vehicle"`
	LicensePlate   string          `gorm:"not null" json:"licensePlate"`
	Date           time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	UserID         uuid.UUID       `gorm:"not null" json:"userID"`
	ClientID       uint            `gorm:"not null" json:"clientID"`
	ClientName     string          `gorm:"not null" json:"clientName"`
	BudgetServices []BudgetService `json:"budgetServices"`
}


type CreateBudgetRequest struct {
	Price        float32   `json:"price"`
	Vehicle      string    `json:"vehicle"`
	LicensePlate string    `json:"licensePlate"`
	UserID       uuid.UUID `json:"userID"`
	ClientID     uint      `json:"clientID"`
	ClientName   string    `json:"clientName"`
	ServiceIDs   []uint    `json:"serviceIDs"`
}

type BudgetService struct {
	BudgetID  uint
	ServiceID uint
}
