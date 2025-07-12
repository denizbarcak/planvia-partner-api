package handlers

import (
	"context"
	"time"

	"github.com/denizbarcak/planvia-partner-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationHandler struct {
	db *mongo.Database
}

func NewReservationHandler(db *mongo.Database) *ReservationHandler {
	return &ReservationHandler{db: db}
}

// CreateReservation yeni bir rezervasyon oluşturur
func (h *ReservationHandler) CreateReservation(c *fiber.Ctx) error {
	// Partner ID'yi context'ten al
	partnerID := c.Locals("partnerId").(string)
	partnerObjID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz partner ID",
		})
	}

	// Request body'yi parse et
	var reservation models.Reservation
	if err := c.BodyParser(&reservation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek formatı",
		})
	}

	// Zorunlu alanları kontrol et
	if reservation.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rezervasyon adı zorunludur",
		})
	}

	// Tarihleri kontrol et
	if reservation.StartDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Başlangıç tarihi zorunludur",
		})
	}

	// Bitiş tarihi kontrolü
	if reservation.EndDate.IsZero() {
		reservation.EndDate = reservation.StartDate
	}

	// Kapasite kontrolü
	if reservation.Capacity < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Kapasite en az 1 olmalıdır",
		})
	}

	// Rezervasyon nesnesini hazırla
	now := time.Now()
	reservation.ID = primitive.NewObjectID()
	reservation.PartnerID = partnerObjID
	reservation.CreatedAt = now
	reservation.UpdatedAt = now

	// Veritabanına kaydet
	_, err = h.db.Collection("reservations").InsertOne(context.Background(), reservation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Rezervasyon kaydedilemedi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(reservation)
}

// GetPartnerReservations partner'a ait rezervasyonları getirir
func (h *ReservationHandler) GetPartnerReservations(c *fiber.Ctx) error {
	// Partner ID'yi context'ten al
	partnerID := c.Locals("partnerId").(string)
	partnerObjID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz partner ID",
		})
	}

	// Tarih filtrelerini al
	startStr := c.Query("start")
	endStr := c.Query("end")

	// Filtreleri oluştur
	filter := bson.M{"partnerId": partnerObjID}
	
	// Tarih filtreleri varsa ekle
	if startStr != "" && endStr != "" {
		startDate, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz başlangıç tarihi formatı",
			})
		}

		endDate, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz bitiş tarihi formatı",
			})
		}

		filter["$or"] = bson.A{
			bson.M{
				"startDate": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
			bson.M{
				"endDate": bson.M{
					"$gte": startDate,
					"$lte": endDate,
				},
			},
		}
	}

	// Rezervasyonları getir
	cursor, err := h.db.Collection("reservations").Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Rezervasyonlar getirilemedi",
		})
	}
	defer cursor.Close(context.Background())

	var reservations []models.Reservation
	if err := cursor.All(context.Background(), &reservations); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Rezervasyonlar parse edilemedi",
		})
	}

	return c.JSON(reservations)
} 