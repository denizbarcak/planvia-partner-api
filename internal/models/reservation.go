package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecurrencePattern tekrarlama desenini tanımlar
type RecurrencePattern struct {
	Enabled      bool     `json:"enabled" bson:"enabled"`
	Type         string   `json:"type" bson:"type"`           // weekly, monthly, yearly
	DaysOfWeek   []int    `json:"daysOfWeek" bson:"daysOfWeek"` // 0-6 (Pazar-Cumartesi)
	EndType      string   `json:"endType" bson:"endType"`     // never, after, on
	EndAfter     int      `json:"endAfter" bson:"endAfter"`   // tekrar sayısı
	EndDate      *time.Time `json:"endDate" bson:"endDate"`     // bitiş tarihi
}

// Reservation modeli
type Reservation struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PartnerID       primitive.ObjectID `json:"partnerId" bson:"partnerId"`
	Name            string            `json:"name" bson:"name"`
	StartDate       time.Time         `json:"startDate" bson:"startDate"`
	EndDate         time.Time         `json:"endDate" bson:"endDate"`
	IsAllDay        bool              `json:"isAllDay" bson:"isAllDay"`
	IsMultiDay      bool              `json:"isMultiDay" bson:"isMultiDay"`
	Capacity        int               `json:"capacity" bson:"capacity"`
	Recurrence      RecurrencePattern `json:"recurrence" bson:"recurrence"`
	CreatedAt       time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt" bson:"updatedAt"`
} 