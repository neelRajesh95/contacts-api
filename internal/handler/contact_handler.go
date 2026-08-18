package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/neelRajesh95/contacts-api/internal/model"
	"github.com/neelRajesh95/contacts-api/internal/repository"
	"github.com/neelRajesh95/contacts-api/internal/service"
)

type ContactHandler struct {
	service *service.ContactService
}

func NewContactHandler(
	service *service.ContactService,
) *ContactHandler {
	return &ContactHandler{
		service: service,
	}
}

func (h *ContactHandler) CreateContact(
	w http.ResponseWriter,
	r *http.Request,
) {

	tenantID := strings.TrimSpace(
		r.Header.Get("X-Tenant-ID"),
	)

	if tenantID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"X-Tenant-ID header is required",
		)
		return
	}

	var request model.CreateContactRequest

	decoder := json.NewDecoder(
		r.Body,
	)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	contact, err := h.service.Create(
		r.Context(),
		tenantID,
		request,
	)

	if err != nil {

		if errors.Is(
			err,
			service.ErrInvalidContact,
		) {
			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to create contact",
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		contact,
	)
}

func (h *ContactHandler) GetContacts(
	w http.ResponseWriter,
	r *http.Request,
) {

	tenantID := strings.TrimSpace(
		r.Header.Get("X-Tenant-ID"),
	)

	if tenantID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"X-Tenant-ID header is required",
		)
		return
	}

	contacts, err := h.service.GetAll(
		r.Context(),
		tenantID,
	)

	if err != nil {

		if errors.Is(
			err,
			service.ErrInvalidContact,
		) {
			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to fetch contacts",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		contacts,
	)
}

func (h *ContactHandler) EnrichContact(
	w http.ResponseWriter,
	r *http.Request,
) {

	tenantID := strings.TrimSpace(
		r.Header.Get("X-Tenant-ID"),
	)

	if tenantID == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"X-Tenant-ID header is required",
		)
		return
	}

	idParam := chi.URLParam(
		r,
		"id",
	)

	id, err := uuid.Parse(idParam)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid contact ID",
		)
		return
	}

	contact, err := h.service.Enrich(
		r.Context(),
		tenantID,
		id,
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrContactNotFound,
		) {
			writeError(
				w,
				http.StatusNotFound,
				"contact not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to enrich contact",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		contact,
	)
}

func writeJSON(
	w http.ResponseWriter,
	statusCode int,
	data interface{},
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(
	w http.ResponseWriter,
	statusCode int,
	message string,
) {

	writeJSON(
		w,
		statusCode,
		map[string]string{
			"error": message,
		},
	)
}