package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/neelRajesh95/contacts-api/internal/model"
)

var ErrContactNotFound = errors.New("contact not found")

type ContactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(
	db *pgxpool.Pool,
) *ContactRepository {
	return &ContactRepository{
		db: db,
	}
}

func (r *ContactRepository) Create(
	ctx context.Context,
	contact *model.Contact,
) (*model.Contact, error) {

	query := `
		INSERT INTO contacts (
			tenant_id,
			first_name,
			last_name,
			email,
			company,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			tenant_id,
			first_name,
			last_name,
			email,
			company,
			status,
			created_at,
			updated_at
	`

	var result model.Contact

	err := r.db.QueryRow(
		ctx,
		query,
		contact.TenantID,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.Company,
		contact.Status,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Company,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create contact: %w",
			err,
		)
	}

	return &result, nil
}

func (r *ContactRepository) GetAll(
	ctx context.Context,
	tenantID string,
) ([]model.Contact, error) {

	query := `
		SELECT
			id,
			tenant_id,
			first_name,
			last_name,
			email,
			company,
			status,
			created_at,
			updated_at
		FROM contacts
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch contacts: %w",
			err,
		)
	}

	defer rows.Close()

	contacts := make([]model.Contact, 0)

	for rows.Next() {

		var contact model.Contact

		err := rows.Scan(
			&contact.ID,
			&contact.TenantID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Email,
			&contact.Company,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan contact: %w",
				err,
			)
		}

		contacts = append(
			contacts,
			contact,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate contacts: %w",
			err,
		)
	}

	return contacts, nil
}

func (r *ContactRepository) Enrich(
	ctx context.Context,
	id uuid.UUID,
	tenantID string,
) (*model.Contact, error) {

	query := `
		UPDATE contacts
		SET
			status = 'enriched',
			updated_at = NOW()
		WHERE id = $1
		  AND tenant_id = $2
		RETURNING
			id,
			tenant_id,
			first_name,
			last_name,
			email,
			company,
			status,
			created_at,
			updated_at
	`

	var contact model.Contact

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		tenantID,
	).Scan(
		&contact.ID,
		&contact.TenantID,
		&contact.FirstName,
		&contact.LastName,
		&contact.Email,
		&contact.Company,
		&contact.Status,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContactNotFound
		}

		return nil, fmt.Errorf(
			"failed to enrich contact: %w",
			err,
		)
	}

	return &contact, nil
}