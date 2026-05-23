package models

import "time"

const (
	StatusDraft      = "Заявка"
	StatusAccepted   = "Принято"
	StatusInProgress = "В работе"
	StatusDone       = "Готово"
	StatusIssued     = "Выдано"
)

var ticketStatuses = map[string]struct{}{
	StatusDraft:      {},
	StatusAccepted:   {},
	StatusInProgress: {},
	StatusDone:       {},
	StatusIssued:     {},
}

func TicketStatuses() []string {
	return []string{StatusDraft, StatusAccepted, StatusInProgress, StatusDone, StatusIssued}
}

func IsValidTicketStatus(status string) bool {
	_, ok := ticketStatuses[status]
	return ok
}

type Workshop struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Master struct {
	ID           int64  `json:"id"`
	WorkshopID   int64  `json:"workshop_id"`
	WorkshopName string `json:"workshop_name,omitempty"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type Device struct {
	ID    int64  `json:"id"`
	IMEI  string `json:"imei"`
	Brand string `json:"brand"`
	Model string `json:"model"`
}

type RepairTicket struct {
	ID                int64     `json:"id"`
	ShortHash         string    `json:"short_hash,omitempty"`
	WorkshopID        int64     `json:"workshop_id"`
	DeviceID          int64     `json:"device_id"`
	ClientName        string    `json:"client_name"`
	ClientPhone       string    `json:"client_phone"`
	Status            string    `json:"status"`
	DefectDescription string    `json:"defect_description"`
	WaterDamage       bool      `json:"water_damage"`
	WarrantyDays      int       `json:"warranty_days"`
	Price             int       `json:"price"`
	Rating            *int      `json:"rating,omitempty"`
	ReviewText        string    `json:"review_text,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	Device            Device    `json:"device"`
}

type PublicApplicationInput struct {
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	Brand             string `json:"brand"`
	Model             string `json:"model"`
	DefectDescription string `json:"defect_description"`
}

type CreateTicketInput struct {
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	IMEI              string `json:"imei"`
	Brand             string `json:"brand"`
	Model             string `json:"model"`
	DefectDescription string `json:"defect_description"`
	WaterDamage       bool   `json:"water_damage"`
	WarrantyDays      int    `json:"warranty_days"`
	Price             int    `json:"price"`
}

type UpdateTicketInput struct {
	ClientName        string `json:"client_name"`
	ClientPhone       string `json:"client_phone"`
	IMEI              string `json:"imei"`
	Brand             string `json:"brand"`
	Model             string `json:"model"`
	Status            string `json:"status"`
	DefectDescription string `json:"defect_description"`
	WaterDamage       bool   `json:"water_damage"`
	WarrantyDays      int    `json:"warranty_days"`
	Price             int    `json:"price"`
}

type ReviewInput struct {
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
}
