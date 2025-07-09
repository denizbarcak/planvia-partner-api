package handlers

import (
	"context"
	"time"

	"planvia-partner-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type PartnerHandler struct {
	collection *mongo.Collection
}

func NewPartnerHandler(db *mongo.Database) *PartnerHandler {
	return &PartnerHandler{
		collection: db.Collection("partners"),
	}
}

func (h *PartnerHandler) Register(c *fiber.Ctx) error {
	var partner models.Partner
	if err := c.BodyParser(&partner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if partner.CompanyName == "" || partner.Email == "" || partner.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Required fields are missing",
		})
	}

	// Check if email already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingPartner := h.collection.FindOne(ctx, bson.M{"email": partner.Email})
	if existingPartner.Err() != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(partner.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Set timestamps
	now := time.Now()
	partner.Password = string(hashedPassword)
	partner.CreatedAt = now
	partner.UpdatedAt = now

	// Insert partner
	result, err := h.collection.InsertOne(ctx, partner)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create partner",
		})
	}

	// Create response without sensitive data
	response := models.PartnerResponse{
		ID:            partner.ID,
		CompanyName:   partner.CompanyName,
		Email:         partner.Email,
		PhoneNumber:   partner.PhoneNumber,
		Address:       partner.Address,
		City:          partner.City,
		BusinessType:  partner.BusinessType,
		TaxNumber:     partner.TaxNumber,
		ContactPerson: partner.ContactPerson,
		CreatedAt:     partner.CreatedAt,
		UpdatedAt:     partner.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Partner registered successfully",
		"partner": response,
		"id":      result.InsertedID,
	})
} 