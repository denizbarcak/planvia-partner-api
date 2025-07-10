package handlers

import (
	"context"
	"fmt"
	"time"

	"planvia-partner-api/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type PartnerHandler struct {
	collection *mongo.Collection
	validate   *validator.Validate
}

func NewPartnerHandler(db *mongo.Database) *PartnerHandler {
	return &PartnerHandler{
		collection: db.Collection("partners"),
		validate:   validator.New(),
	}
}

func (h *PartnerHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Parse request body
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek formatı",
		})
	}

	// Log received data
	fmt.Printf("Received registration request: %+v\n", req)

	// Validate request data
	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			errorMessages[i] = translateValidationError(e)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation error",
			"details": errorMessages,
		})
	}

	// Check if email already exists
	existingPartner := h.collection.FindOne(ctx, bson.M{"email": req.Email})
	if existingPartner.Err() != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Bu e-posta adresi zaten kullanımda",
		})
	}

	// Check if tax number already exists
	existingPartner = h.collection.FindOne(ctx, bson.M{"tax_number": req.TaxNumber})
	if existingPartner.Err() != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Bu vergi numarası zaten kullanımda",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Şifre işlenirken bir hata oluştu",
		})
	}

	// Create partner from request
	partner := req.ToPartner()
	partner.Password = string(hashedPassword)

	// Insert partner
	result, err := h.collection.InsertOne(ctx, partner)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Partner kaydedilirken bir hata oluştu",
		})
	}

	// Set ID from insert result
	partner.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "İşletme başarıyla kaydedildi",
		"partner": partner.ToResponse(),
	})
}

func (h *PartnerHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek formatı",
		})
	}

	// Validate login data
	if err := h.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			errorMessages[i] = translateValidationError(e)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation error",
			"details": errorMessages,
		})
	}

	// Find partner by email
	var partner models.Partner
	err := h.collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&partner)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "E-posta veya şifre hatalı",
		})
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(partner.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "E-posta veya şifre hatalı",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Giriş başarılı",
		"partner": partner.ToResponse(),
	})
}

func translateValidationError(e validator.FieldError) string {
	switch e.Field() {
	case "CompanyName":
		return "İşletme adı zorunludur"
	case "Email":
		if e.Tag() == "required" {
			return "E-posta adresi zorunludur"
		}
		return "Geçerli bir e-posta adresi giriniz"
	case "Password":
		if e.Tag() == "required" {
			return "Şifre zorunludur"
		}
		return "Şifre en az 6 karakter olmalıdır"
	case "PhoneNumber":
		return "Telefon numarası zorunludur"
	case "Address":
		return "Adres zorunludur"
	case "City":
		return "Şehir seçimi zorunludur"
	case "BusinessType":
		return "İşletme kategorisi seçimi zorunludur"
	case "TaxNumber":
		return "Vergi numarası zorunludur"
	case "ContactPerson":
		return "Yetkili kişi bilgisi zorunludur"
	default:
		return fmt.Sprintf("%s alanı için %s kuralı geçerli değil", e.Field(), e.Tag())
	}
} 