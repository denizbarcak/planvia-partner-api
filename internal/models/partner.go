package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Partner represents a business partner in the system
type Partner struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyName   string            `bson:"company_name" json:"companyName" validate:"required"`
	Email         string            `bson:"email" json:"email" validate:"required,email"`
	Password      string            `bson:"password" json:"password" validate:"required,min=6"`
	PhoneNumber   string            `bson:"phone_number" json:"phoneNumber" validate:"required"`
	Address       string            `bson:"address" json:"address" validate:"required"`
	City          string            `bson:"city" json:"city" validate:"required"`
	BusinessType  string            `bson:"business_type" json:"businessType" validate:"required"`
	TaxNumber     string            `bson:"tax_number" json:"taxNumber" validate:"required"`
	ContactPerson string            `bson:"contact_person" json:"contactPerson" validate:"required"`
	CreatedAt     time.Time         `bson:"created_at" json:"createdAt,omitempty"`
	UpdatedAt     time.Time         `bson:"updated_at" json:"updatedAt,omitempty"`
}

// RegisterRequest represents the registration request data
type RegisterRequest struct {
	CompanyName   string `json:"companyName" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=6"`
	PhoneNumber   string `json:"phoneNumber" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	BusinessType  string `json:"businessType" validate:"required"`
	TaxNumber     string `json:"taxNumber" validate:"required"`
	ContactPerson string `json:"contactPerson" validate:"required"`
}

// LoginRequest represents the login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// PartnerResponse represents the partner data that is safe to send to the client
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

// ToResponse converts a Partner to a PartnerResponse
func (p *Partner) ToResponse() PartnerResponse {
	return PartnerResponse{
		ID:            p.ID,
		CompanyName:   p.CompanyName,
		Email:         p.Email,
		PhoneNumber:   p.PhoneNumber,
		Address:       p.Address,
		City:          p.City,
		BusinessType:  p.BusinessType,
		TaxNumber:     p.TaxNumber,
		ContactPerson: p.ContactPerson,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// ToPartner converts a RegisterRequest to a Partner
func (r *RegisterRequest) ToPartner() Partner {
	now := time.Now()
	return Partner{
		CompanyName:   r.CompanyName,
		Email:         r.Email,
		Password:      r.Password,
		PhoneNumber:   r.PhoneNumber,
		Address:       r.Address,
		City:          r.City,
		BusinessType:  r.BusinessType,
		TaxNumber:     r.TaxNumber,
		ContactPerson: r.ContactPerson,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
} 