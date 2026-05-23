package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/example/repair-crm/internal/models"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *models.RepairTicket) (*models.RepairTicket, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	device, err := upsertDevice(ctx, tx, ticket.Device)
	if err != nil {
		return nil, err
	}

	const query = `
		INSERT INTO repair_tickets (
			short_hash, workshop_id, device_id, client_name, client_phone,
			status, defect_description, water_damage, warranty_days, price
		)
		VALUES (NULLIF($1, ''), $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, short_hash, workshop_id, device_id, client_name, client_phone,
			status, defect_description, water_damage, warranty_days, price,
			rating, review_text, created_at
	`

	created, err := scanTicket(tx.QueryRowContext(
		ctx,
		query,
		ticket.ShortHash,
		ticket.WorkshopID,
		device.ID,
		ticket.ClientName,
		ticket.ClientPhone,
		ticket.Status,
		ticket.DefectDescription,
		ticket.WaterDamage,
		ticket.WarrantyDays,
		ticket.Price,
	), device)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return created, nil
}

func (r *TicketRepository) ListByWorkshop(ctx context.Context, workshopID int64, status string) ([]models.RepairTicket, error) {
	const query = `
		SELECT t.id, t.short_hash, t.workshop_id, t.device_id, t.client_name, t.client_phone,
			t.status, t.defect_description, t.water_damage, t.warranty_days, t.price,
			t.rating, t.review_text, t.created_at,
			d.id, d.imei, d.brand, d.model
		FROM repair_tickets t
		JOIN devices d ON d.id = t.device_id
		WHERE t.workshop_id = $1 AND ($2 = '' OR t.status = $2)
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, workshopID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := make([]models.RepairTicket, 0)
	for rows.Next() {
		ticket, err := scanTicketWithDevice(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) GetByID(ctx context.Context, workshopID, id int64) (*models.RepairTicket, error) {
	const query = `
		SELECT t.id, t.short_hash, t.workshop_id, t.device_id, t.client_name, t.client_phone,
			t.status, t.defect_description, t.water_damage, t.warranty_days, t.price,
			t.rating, t.review_text, t.created_at,
			d.id, d.imei, d.brand, d.model
		FROM repair_tickets t
		JOIN devices d ON d.id = t.device_id
		WHERE t.workshop_id = $1 AND t.id = $2
	`

	ticket, err := scanTicketWithDevice(r.db.QueryRowContext(ctx, query, workshopID, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *TicketRepository) GetByHash(ctx context.Context, shortHash string) (*models.RepairTicket, error) {
	const query = `
		SELECT t.id, t.short_hash, t.workshop_id, t.device_id, t.client_name, t.client_phone,
			t.status, t.defect_description, t.water_damage, t.warranty_days, t.price,
			t.rating, t.review_text, t.created_at,
			d.id, d.imei, d.brand, d.model
		FROM repair_tickets t
		JOIN devices d ON d.id = t.device_id
		WHERE t.short_hash = $1
	`

	ticket, err := scanTicketWithDevice(r.db.QueryRowContext(ctx, query, shortHash))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *TicketRepository) Update(ctx context.Context, workshopID, id int64, ticket *models.RepairTicket) (*models.RepairTicket, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	device, err := upsertDevice(ctx, tx, ticket.Device)
	if err != nil {
		return nil, err
	}

	const query = `
		UPDATE repair_tickets
		SET device_id = $1,
			client_name = $2,
			client_phone = $3,
			status = $4,
			defect_description = $5,
			water_damage = $6,
			warranty_days = $7,
			price = $8,
			short_hash = CASE WHEN $9 = '' THEN short_hash ELSE $9 END
		WHERE workshop_id = $10 AND id = $11
		RETURNING id, short_hash, workshop_id, device_id, client_name, client_phone,
			status, defect_description, water_damage, warranty_days, price,
			rating, review_text, created_at
	`

	updated, err := scanTicket(tx.QueryRowContext(
		ctx,
		query,
		device.ID,
		ticket.ClientName,
		ticket.ClientPhone,
		ticket.Status,
		ticket.DefectDescription,
		ticket.WaterDamage,
		ticket.WarrantyDays,
		ticket.Price,
		ticket.ShortHash,
		workshopID,
		id,
	), device)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *TicketRepository) AddReview(ctx context.Context, ticketID int64, rating int, reviewText string) (*models.RepairTicket, error) {
	const query = `
		UPDATE repair_tickets
		SET rating = $1, review_text = $2
		WHERE id = $3
		RETURNING id, short_hash, workshop_id, device_id, client_name, client_phone,
			status, defect_description, water_damage, warranty_days, price,
			rating, review_text, created_at
	`

	ticket, err := scanTicket(r.db.QueryRowContext(ctx, query, rating, reviewText, ticketID), models.Device{})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	fullTicket, err := r.GetByID(ctx, ticket.WorkshopID, ticket.ID)
	if err != nil {
		return nil, err
	}

	return fullTicket, nil
}

func (r *TicketRepository) Delete(ctx context.Context, workshopID, id int64) error {
	const query = `DELETE FROM repair_tickets WHERE workshop_id = $1 AND id = $2`

	result, err := r.db.ExecContext(ctx, query, workshopID, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TicketRepository) ShortHashExists(ctx context.Context, shortHash string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM repair_tickets WHERE short_hash = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, shortHash).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func upsertDevice(ctx context.Context, tx *sql.Tx, device models.Device) (models.Device, error) {
	const query = `
		INSERT INTO devices (imei, brand, model)
		VALUES (NULLIF($1, ''), $2, $3)
		ON CONFLICT (imei) DO UPDATE SET
			brand = EXCLUDED.brand,
			model = EXCLUDED.model
		RETURNING id, COALESCE(imei, ''), brand, model
	`

	var saved models.Device
	err := tx.QueryRowContext(ctx, query, device.IMEI, device.Brand, device.Model).Scan(
		&saved.ID,
		&saved.IMEI,
		&saved.Brand,
		&saved.Model,
	)
	return saved, err
}

func scanTicket(scanner rowScanner, device models.Device) (*models.RepairTicket, error) {
	var ticket models.RepairTicket
	var shortHash sql.NullString
	var rating sql.NullInt64
	var reviewText sql.NullString
	err := scanner.Scan(
		&ticket.ID,
		&shortHash,
		&ticket.WorkshopID,
		&ticket.DeviceID,
		&ticket.ClientName,
		&ticket.ClientPhone,
		&ticket.Status,
		&ticket.DefectDescription,
		&ticket.WaterDamage,
		&ticket.WarrantyDays,
		&ticket.Price,
		&rating,
		&reviewText,
		&ticket.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if shortHash.Valid {
		ticket.ShortHash = shortHash.String
	}
	if rating.Valid {
		value := int(rating.Int64)
		ticket.Rating = &value
	}
	if reviewText.Valid {
		ticket.ReviewText = reviewText.String
	}

	ticket.Device = device
	return &ticket, nil
}

func scanTicketWithDevice(scanner rowScanner) (*models.RepairTicket, error) {
	var ticket models.RepairTicket
	var shortHash sql.NullString
	var rating sql.NullInt64
	var reviewText sql.NullString
	var imei sql.NullString
	err := scanner.Scan(
		&ticket.ID,
		&shortHash,
		&ticket.WorkshopID,
		&ticket.DeviceID,
		&ticket.ClientName,
		&ticket.ClientPhone,
		&ticket.Status,
		&ticket.DefectDescription,
		&ticket.WaterDamage,
		&ticket.WarrantyDays,
		&ticket.Price,
		&rating,
		&reviewText,
		&ticket.CreatedAt,
		&ticket.Device.ID,
		&imei,
		&ticket.Device.Brand,
		&ticket.Device.Model,
	)
	if err != nil {
		return nil, err
	}

	if shortHash.Valid {
		ticket.ShortHash = shortHash.String
	}
	if rating.Valid {
		value := int(rating.Int64)
		ticket.Rating = &value
	}
	if reviewText.Valid {
		ticket.ReviewText = reviewText.String
	}
	if imei.Valid {
		ticket.Device.IMEI = imei.String
	}

	return &ticket, nil
}
