package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Partner struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyName     string            `json:"companyName" bson:"company_name"`
	Email           string            `json:"email" bson:"email"`
	Password        string            `json:"password,omitempty" bson:"password"`
	PhoneNumber     string            `json:"phoneNumber" bson:"phone_number"`
	Address         string            `json:"address" bson:"address"`
	City            string            `json:"city" bson:"city"`
	BusinessType    string            `json:"businessType" bson:"business_type"`
	TaxNumber       string            `json:"taxNumber" bson:"tax_number"`
	ContactPerson   string            `json:"contactPerson" bson:"contact_person"`
	CreatedAt       time.Time         `json:"createdAt" bson:"created_at"`
	UpdatedAt       time.Time         `json:"updatedAt" bson:"updated_at"`
}

type PartnerResponse struct {
	ID            primitive.ObjectID `json:"id"`
	CompanyName   string            `json:"companyName"`
	Email         string            `json:"email"`
	PhoneNumber   string            `json:"phoneNumber"`
	Address       string            `json:"address"`
	City          string            `json:"city"`
	BusinessType  string            `json:"businessType"`
	TaxNumber     string            `json:"taxNumber"`
	ContactPerson string            `json:"contactPerson"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
} 