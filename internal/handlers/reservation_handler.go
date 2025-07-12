package handlers

import (
	"context"
	"time"

	"github.com/denizbarcak/planvia-partner-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// UpdateReservation günceller bir rezervasyonu
func (h *ReservationHandler) UpdateReservation(c *fiber.Ctx) error {
	// Partner ID'yi context'ten al
	partnerID := c.Locals("partnerId").(string)
	partnerObjID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz partner ID",
		})
	}

	// Rezervasyon ID'yi URL'den al
	reservationID := c.Params("id")
	reservationObjID, err := primitive.ObjectIDFromHex(reservationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz rezervasyon ID",
		})
	}

	// Request body'yi parse et
	var updateData models.Reservation
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek formatı",
		})
	}

	// Zorunlu alanları kontrol et
	if updateData.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rezervasyon adı zorunludur",
		})
	}

	// Tarihleri kontrol et
	if updateData.StartDate.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Başlangıç tarihi zorunludur",
		})
	}

	// Bitiş tarihi kontrolü
	if updateData.EndDate.IsZero() {
		updateData.EndDate = updateData.StartDate
	}

	// Kapasite kontrolü
	if updateData.Capacity < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Kapasite en az 1 olmalıdır",
		})
	}

	// Rezervasyonun mevcut olduğunu ve bu partner'a ait olduğunu kontrol et
	filter := bson.M{
		"_id":       reservationObjID,
		"partnerId": partnerObjID,
	}

	// Güncellenecek alanları hazırla
	update := bson.M{
		"$set": bson.M{
			"name":       updateData.Name,
			"startDate":  updateData.StartDate,
			"endDate":    updateData.EndDate,
			"isAllDay":   updateData.IsAllDay,
			"isMultiDay": updateData.IsMultiDay,
			"capacity":   updateData.Capacity,
			"recurrence": updateData.Recurrence,
			"updatedAt":  time.Now(),
		},
	}

	// Güncelleme işlemini gerçekleştir
	result := h.db.Collection("reservations").FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Rezervasyon bulunamadı veya bu partner'a ait değil",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Rezervasyon güncellenirken bir hata oluştu",
		})
	}

	// Güncellenmiş rezervasyonu döndür
	var updatedReservation models.Reservation
	if err := result.Decode(&updatedReservation); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Güncellenmiş rezervasyon alınamadı",
		})
	}

	return c.JSON(updatedReservation)
} 

// DeleteReservation bir rezervasyonu siler
func (h *ReservationHandler) DeleteReservation(c *fiber.Ctx) error {
	// Partner ID'yi context'ten al
	partnerID := c.Locals("partnerId").(string)
	partnerObjID, err := primitive.ObjectIDFromHex(partnerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz partner ID",
		})
	}

	// Rezervasyon ID'yi URL'den al
	reservationID := c.Params("id")
	reservationObjID, err := primitive.ObjectIDFromHex(reservationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz rezervasyon ID",
		})
	}

	// Rezervasyonun mevcut olduğunu ve bu partner'a ait olduğunu kontrol et
	filter := bson.M{
		"_id":       reservationObjID,
		"partnerId": partnerObjID,
	}

	// Silme işlemini gerçekleştir
	result, err := h.db.Collection("reservations").DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Rezervasyon silinirken bir hata oluştu",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rezervasyon bulunamadı veya bu partner'a ait değil",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rezervasyon başarıyla silindi",
	})
} 