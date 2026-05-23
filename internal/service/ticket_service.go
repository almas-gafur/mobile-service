package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/example/repair-crm/internal/models"
	"github.com/example/repair-crm/internal/repository"
	shorthash "github.com/example/repair-crm/pkg/hash"
)

var (
	imeiPattern  = regexp.MustCompile(`^[0-9]{15}$`)
	hashPattern  = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)
	phonePattern = regexp.MustCompile(`^[0-9+()\-\s]{7,20}$`)
)

type TicketService struct {
	tickets           *repository.TicketRepository
	defaultWorkshopID int64
}

func NewTicketService(tickets *repository.TicketRepository, defaultWorkshopID int64) *TicketService {
	return &TicketService{tickets: tickets, defaultWorkshopID: defaultWorkshopID}
}

func (s *TicketService) CreatePublicApplication(ctx context.Context, input models.PublicApplicationInput) (*models.RepairTicket, error) {
	if s.defaultWorkshopID <= 0 {
		return nil, ValidationError{Message: "мастерская для публичных заявок не настроена"}
	}
	if err := validatePublicApplication(input); err != nil {
		return nil, err
	}

	ticket := &models.RepairTicket{
		WorkshopID:        s.defaultWorkshopID,
		ClientName:        clean(input.ClientName),
		ClientPhone:       clean(input.ClientPhone),
		Status:            models.StatusDraft,
		DefectDescription: clean(input.DefectDescription),
		Device: models.Device{
			Brand: normalizeBrand(input.Brand),
			Model: clean(input.Model),
		},
	}

	return s.tickets.Create(ctx, ticket)
}

func (s *TicketService) Create(ctx context.Context, workshopID int64, input models.CreateTicketInput) (*models.RepairTicket, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	shortHash, err := s.generateUniqueShortHash(ctx)
	if err != nil {
		return nil, err
	}

	ticket := &models.RepairTicket{
		ShortHash:         shortHash,
		WorkshopID:        workshopID,
		ClientName:        clean(input.ClientName),
		ClientPhone:       clean(input.ClientPhone),
		Status:            models.StatusAccepted,
		DefectDescription: clean(input.DefectDescription),
		WaterDamage:       input.WaterDamage,
		WarrantyDays:      input.WarrantyDays,
		Price:             input.Price,
		Device: models.Device{
			IMEI:  clean(input.IMEI),
			Brand: normalizeBrand(input.Brand),
			Model: clean(input.Model),
		},
	}

	return s.tickets.Create(ctx, ticket)
}

func (s *TicketService) List(ctx context.Context, workshopID int64, status string) ([]models.RepairTicket, error) {
	status = clean(status)
	if status != "" && !models.IsValidTicketStatus(status) {
		return nil, ValidationError{Message: "некорректный статус заявки"}
	}

	return s.tickets.ListByWorkshop(ctx, workshopID, status)
}

func (s *TicketService) Get(ctx context.Context, workshopID, id int64) (*models.RepairTicket, error) {
	return s.tickets.GetByID(ctx, workshopID, id)
}

func (s *TicketService) GetPublic(ctx context.Context, shortHash string) (*models.RepairTicket, error) {
	shortHash = clean(shortHash)
	if !hashPattern.MatchString(shortHash) {
		return nil, ValidationError{Message: "некорректная ссылка отслеживания"}
	}

	return s.tickets.GetByHash(ctx, shortHash)
}

func (s *TicketService) Update(ctx context.Context, workshopID, id int64, input models.UpdateTicketInput) (*models.RepairTicket, error) {
	if err := validateUpdateInput(input); err != nil {
		return nil, err
	}

	existing, err := s.tickets.GetByID(ctx, workshopID, id)
	if err != nil {
		return nil, err
	}

	nextHash := existing.ShortHash
	if existing.ShortHash == "" && input.Status != models.StatusDraft {
		nextHash, err = s.generateUniqueShortHash(ctx)
		if err != nil {
			return nil, err
		}
	}

	ticket := &models.RepairTicket{
		ShortHash:         nextHash,
		WorkshopID:        workshopID,
		ClientName:        clean(input.ClientName),
		ClientPhone:       clean(input.ClientPhone),
		Status:            clean(input.Status),
		DefectDescription: clean(input.DefectDescription),
		WaterDamage:       input.WaterDamage,
		WarrantyDays:      input.WarrantyDays,
		Price:             input.Price,
		Device: models.Device{
			IMEI:  clean(input.IMEI),
			Brand: normalizeBrand(input.Brand),
			Model: clean(input.Model),
		},
	}

	return s.tickets.Update(ctx, workshopID, id, ticket)
}

func (s *TicketService) AddReview(ctx context.Context, shortHash string, input models.ReviewInput) (*models.RepairTicket, error) {
	if err := validateReviewInput(input); err != nil {
		return nil, err
	}

	ticket, err := s.GetPublic(ctx, shortHash)
	if err != nil {
		return nil, err
	}
	if ticket.Status != models.StatusIssued {
		return nil, ValidationError{Message: "отзыв можно оставить только после выдачи аппарата"}
	}
	if ticket.Rating != nil {
		return nil, ValidationError{Message: "отзыв по этой заявке уже отправлен"}
	}

	return s.tickets.AddReview(ctx, ticket.ID, input.Rating, clean(input.ReviewText))
}

func (s *TicketService) Delete(ctx context.Context, workshopID, id int64) error {
	return s.tickets.Delete(ctx, workshopID, id)
}

func (s *TicketService) generateUniqueShortHash(ctx context.Context) (string, error) {
	for i := 0; i < 12; i++ {
		value, err := shorthash.Generate(8)
		if err != nil {
			return "", err
		}

		exists, err := s.tickets.ShortHashExists(ctx, value)
		if err != nil {
			return "", err
		}
		if !exists {
			return value, nil
		}
	}

	return "", errors.New("cannot generate unique short hash")
}

func validatePublicApplication(input models.PublicApplicationInput) error {
	if len([]rune(clean(input.ClientName))) < 2 {
		return ValidationError{Message: "укажите имя клиента"}
	}
	if !isValidPhone(input.ClientPhone) {
		return ValidationError{Message: "укажите корректный телефон от 7 до 20 символов"}
	}
	if clean(input.Model) == "" {
		return ValidationError{Message: "укажите модель телефона"}
	}
	if len([]rune(clean(input.DefectDescription))) < 5 {
		return ValidationError{Message: "опишите поломку подробнее"}
	}

	return nil
}

func validateCreateInput(input models.CreateTicketInput) error {
	if len([]rune(clean(input.ClientName))) < 2 {
		return ValidationError{Message: "укажите имя клиента"}
	}
	if !isValidPhone(input.ClientPhone) {
		return ValidationError{Message: "укажите корректный телефон от 7 до 20 символов"}
	}
	if clean(input.IMEI) != "" && !imeiPattern.MatchString(clean(input.IMEI)) {
		return ValidationError{Message: "IMEI должен содержать ровно 15 цифр"}
	}
	if clean(input.Model) == "" {
		return ValidationError{Message: "укажите модель телефона"}
	}
	if input.WarrantyDays < 0 || input.WarrantyDays > 730 {
		return ValidationError{Message: "гарантия должна быть от 0 до 730 дней"}
	}
	if input.Price < 0 {
		return ValidationError{Message: "стоимость не может быть отрицательной"}
	}

	return nil
}

func validateUpdateInput(input models.UpdateTicketInput) error {
	if len([]rune(clean(input.ClientName))) < 2 {
		return ValidationError{Message: "укажите имя клиента"}
	}
	if !isValidPhone(input.ClientPhone) {
		return ValidationError{Message: "укажите корректный телефон от 7 до 20 символов"}
	}
	if clean(input.IMEI) != "" && !imeiPattern.MatchString(clean(input.IMEI)) {
		return ValidationError{Message: "IMEI должен содержать ровно 15 цифр"}
	}
	if clean(input.Model) == "" {
		return ValidationError{Message: "укажите модель телефона"}
	}
	if !models.IsValidTicketStatus(clean(input.Status)) {
		return ValidationError{Message: "некорректный статус заявки"}
	}
	if input.WarrantyDays < 0 || input.WarrantyDays > 730 {
		return ValidationError{Message: "гарантия должна быть от 0 до 730 дней"}
	}
	if input.Price < 0 {
		return ValidationError{Message: "стоимость не может быть отрицательной"}
	}

	return nil
}

func validateReviewInput(input models.ReviewInput) error {
	if input.Rating < 1 || input.Rating > 5 {
		return ValidationError{Message: "оценка должна быть от 1 до 5"}
	}
	textLength := len([]rune(clean(input.ReviewText)))
	if textLength < 5 || textLength > 1000 {
		return ValidationError{Message: "отзыв должен быть от 5 до 1000 символов"}
	}

	return nil
}

func isValidPhone(phone string) bool {
	return phonePattern.MatchString(clean(phone))
}

func normalizeBrand(brand string) string {
	brand = clean(brand)
	if brand == "" {
		return "Не указан"
	}
	return brand
}

func clean(value string) string {
	return strings.TrimSpace(value)
}
