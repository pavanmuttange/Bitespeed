package model

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID             uint    `gorm:"primaryKey,not null,Unique"`
	PhoneNumber    *string `gorm:"size:15"`
	Email          *string `gorm:"size:100"`
	LinkedID       *uint
	LinkPrecedence string `gorm:"size:10"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (e *Contact) TableName() string {
	return "CONTACT"
}

type IdentifyRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type ContactResponse struct {
	PrimaryContactID    uint     `json:"primaryContactId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []uint   `json:"secondaryContactIds"`
}
