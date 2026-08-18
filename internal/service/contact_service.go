package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"

	"github.com/neelRajesh95/contacts-api/internal/model"
)

var ErrInvalidContact = errors.New(
	"invalid contact",
)

type ContactRepository interface {
	Create(
		ctx context.Context,
		contact *model.Contact,
	) (*model.Contact, error)

	GetAll(
		ctx context.Context,
		tenantID string,
	) ([]model.Contact, error)

	Enrich(
		ctx context.Context,
		id uuid.UUID,
		tenantID string,
	) (*model.Contact, error)
}

type ContactService struct {
	repository ContactRepository
}

func NewContactService(
	repository ContactRepository,
) *ContactService {
	return &ContactService{
		repository: repository,
	}
}

func (s *ContactService) Create(
	ctx context.Context,
	tenantID string,
	request model.CreateContactRequest,
) (*model.Contact, error) {

	tenantID = strings.TrimSpace(tenantID)

	if tenantID == "" {
		return nil, fmt.Errorf(
			"%w: tenant ID is required",
			ErrInvalidContact,
		)
	}

	firstName := strings.TrimSpace(
		request.FirstName,
	)

	lastName := strings.TrimSpace(
		request.LastName,
	)

	email := strings.TrimSpace(
		request.Email,
	)

	company := strings.TrimSpace(
		request.Company,
	)

	if firstName == "" {
		return nil, fmt.Errorf(
			"%w: first_name is required",
			ErrInvalidContact,
		)
	}

	if lastName == "" {
		return nil, fmt.Errorf(
			"%w: last_name is required",
			ErrInvalidContact,
		)
	}

	if email == "" {
		return nil, fmt.Errorf(
			"%w: email is required",
			ErrInvalidContact,
		)
	}

	parsedEmail, err := mail.ParseAddress(email)

	if err != nil ||
		parsedEmail.Address != email {

		return nil, fmt.Errorf(
			"%w: invalid email",
			ErrInvalidContact,
		)
	}

	contact := &model.Contact{
		TenantID:  tenantID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Company:   company,
		Status:    "pending",
	}

	return s.repository.Create(
		ctx,
		contact,
	)
}

func (s *ContactService) GetAll(
	ctx context.Context,
	tenantID string,
) ([]model.Contact, error) {

	tenantID = strings.TrimSpace(tenantID)

	if tenantID == "" {
		return nil, fmt.Errorf(
			"%w: tenant ID is required",
			ErrInvalidContact,
		)
	}

	return s.repository.GetAll(
		ctx,
		tenantID,
	)
}

func (s *ContactService) Enrich(
	ctx context.Context,
	tenantID string,
	id uuid.UUID,
) (*model.Contact, error) {

	tenantID = strings.TrimSpace(tenantID)

	if tenantID == "" {
		return nil, fmt.Errorf(
			"%w: tenant ID is required",
			ErrInvalidContact,
		)
	}

	return s.repository.Enrich(
		ctx,
		id,
		tenantID,
	)
}