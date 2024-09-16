package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	errwriter "github.com/Te8va/Tender/internal/pkg/errWriter"
	"github.com/Te8va/Tender/internal/tender/domain"
	"github.com/Te8va/Tender/pkg/logger"
)

type PingHandler struct {
	srv domain.TenderServicePingProvider
}

func NewPingProvider(srv domain.TenderServicePingProvider) *PingHandler {
	return &PingHandler{srv: srv}
}

func (h *PingHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.srv.Ping(r.Context())
	if err != nil {
		errwriter.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logger.Logger().Errorln("Error pinging service:", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("ok"))
	if err != nil {
		logger.Logger().Errorln("Error writing response:", err.Error())
	}
}

type TenderHandler struct {
	srv domain.TenderService
}

func NewTenderHandler(srv domain.TenderService) *TenderHandler {
	return &TenderHandler{srv: srv}
}

func (h *TenderHandler) ListTenderHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	serviceTypes := r.URL.Query()["service_type"]

	limit := 10
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	tenders, err := h.srv.ListTender(r.Context(), limit, offset, serviceTypes)
	if err != nil {
		errwriter.RespondWithError(w, http.StatusBadRequest, err.Error())
		logger.Logger().Errorln("Error fetching tender list:", err.Error())
		return
	}

	var responseTenders []domain.TenderResponse
	for _, tender := range tenders {
		responseTenders = append(responseTenders, domain.TenderResponse{
			ID:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			Version:     tender.Version,
			CreatedAt:   tender.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseTenders); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

func (h *TenderHandler) GetUserTendersHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Missing username")
		logger.Logger().Errorln("Error: Missing username in query parameters")
		return
	}

	tenders, err := h.srv.GetUserTenders(r.Context(), limit, offset, username)
	if err != nil {
		var statusCode int
		if err.Error() == "service.CreateTender: user does not exist" {
			statusCode = http.StatusBadRequest
		} else {
			statusCode = http.StatusInternalServerError
		}
		errwriter.RespondWithError(w, statusCode, err.Error())
		return
	}

	var responseTenders []domain.TenderResponse
	for _, tender := range tenders {
		responseTenders = append(responseTenders, domain.TenderResponse{
			ID:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			Version:     tender.Version,
			CreatedAt:   tender.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseTenders); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

func (h *TenderHandler) CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTenderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		logger.Logger().Errorln("Error decoding request payload:", err.Error())
		return
	}

	if req.Name == "" || req.ServiceType == "" || req.OrganizationId == "" || req.CreatorUsername == "" {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		logger.Logger().Errorln("Error: Missing required fields in request")
		return
	}

	newTender := domain.Tender{
		Name:            req.Name,
		Description:     req.Description,
		ServiceType:     req.ServiceType,
		Status:          "CREATED",
		OrganizationId:  req.OrganizationId,
		CreatorUsername: req.CreatorUsername,
		Version:         1,
		CreatedAt:       time.Now(),
	}

	createdTender, err := h.srv.CreateTender(r.Context(), newTender)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "service.CreateTender: user does not exist":
			statusCode = http.StatusUnauthorized
		case "service.CreateTender: user is not authorized to create tender for this organization":
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusInternalServerError
		}
		errwriter.RespondWithError(w, statusCode, err.Error())
		return
	}

	response := domain.TenderResponse{
		ID:          createdTender.ID,
		Name:        createdTender.Name,
		Description: createdTender.Description,
		Status:      createdTender.Status,
		ServiceType: createdTender.ServiceType,
		Version:     createdTender.Version,
		CreatedAt:   createdTender.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

func (h *TenderHandler) GetTenderStatusHandler(w http.ResponseWriter, r *http.Request) {
	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid tender ID")
		logger.Logger().Errorln("Error: Invalid tender ID")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		errwriter.RespondWithError(w, http.StatusUnauthorized, "Missing username")
		logger.Logger().Errorln("Error: Missing username in query parameters")
		return
	}

	status, err := h.srv.GetTenderStatus(r.Context(), tenderID, username)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "service.GetTenderByID: repository.GetTenderStatus: %!w(<nil>)":
			statusCode = http.StatusForbidden
		case "service.GetTenderByID: repository.GetTenderStatus: no rows in result set":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		errwriter.RespondWithError(w, statusCode, err.Error())
		return
	}

	response := map[string]string{"status": status}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

func (h *TenderHandler) UpdateTenderStatusHandler(w http.ResponseWriter, r *http.Request) {
	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid tender ID")
		logger.Logger().Errorln("Error: Invalid tender ID")
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		errwriter.RespondWithError(w, http.StatusUnauthorized, "Missing username")
		logger.Logger().Errorln("Error: Missing username in query parameters")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		errwriter.RespondWithError(w, http.StatusUnauthorized, "Missing username")
		logger.Logger().Errorln("Error: Missing username in query parameters")
		return
	}

	updatedTender, err := h.srv.UpdateTenderStatus(r.Context(), tenderID, status, username)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "service.UpdateTenderStatus: repository.UpdateTenderStatus: %!w(<nil>)":
			statusCode = http.StatusForbidden
		case "service.UpdateTenderStatus: repository.UpdateTenderStatus: ERROR: new row for relation \"tender\" violates check constraint \"tender_status_check\" (SQLSTATE 23514)":
			statusCode = http.StatusBadRequest
		case "service.UpdateTenderStatus: no rows updated; check the ID":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		errwriter.RespondWithError(w, statusCode, err.Error())
		return
	}

	response := domain.TenderResponse{
		ID:          updatedTender.ID,
		Name:        updatedTender.Name,
		Description: updatedTender.Description,
		Status:      updatedTender.Status,
		ServiceType: updatedTender.ServiceType,
		Version:     updatedTender.Version,
		CreatedAt:   updatedTender.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

func (h *TenderHandler) UpdatePartTenderHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		errwriter.RespondWithError(w, http.StatusUnauthorized, "Missing username")
		logger.Logger().Errorln("Error: Missing username in query parameters")
		return
	}

	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid tender ID")
		logger.Logger().Errorln("Error: Invalid tender ID")
		return
	}

	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		logger.Logger().Errorln("Error decoding request payload:", err.Error())
		return
	}

	updatedTender, err := h.srv.UpdatePartTender(r.Context(), tenderID, updates, username)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "failed to update tender in repository: repository.UpdatePartTender: %!w(<nil>)":
			statusCode = http.StatusForbidden
		case "service.UpdateTenderStatus: repository.UpdatePartTender: ERROR: new row for relation \"tender\" violates check constraint \"tender_status_check\" (SQLSTATE 23514)":
			statusCode = http.StatusBadRequest
		case "failed to update tender in repository: error fetching current version: no rows in result set":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		errwriter.RespondWithError(w, statusCode, err.Error())
		return
	}

	response := domain.TenderResponse{
		ID:          updatedTender.ID,
		Name:        updatedTender.Name,
		Description: updatedTender.Description,
		Status:      updatedTender.Status,
		ServiceType: updatedTender.ServiceType,
		Version:     updatedTender.Version,
		CreatedAt:   updatedTender.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}

func (h *TenderHandler) RollbackTenderHandler(w http.ResponseWriter, r *http.Request) {
	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid tender ID")
		logger.Logger().Errorln("Error: Invalid tender ID ")
		return
	}

	path := r.URL.Path

	parts := strings.Split(path, "/")

	if len(parts) < 6 {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		logger.Logger().Errorln("Error: Invalid URL format")
		return
	}

	rollbackVersion := parts[5]

	version, err := strconv.Atoi(rollbackVersion)
	if err != nil {
		errwriter.RespondWithError(w, http.StatusBadRequest, "Invalid version format")
		logger.Logger().Errorln("Error: Invalid version format")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		errwriter.RespondWithError(w, http.StatusUnauthorized, "Missing username")
		logger.Logger().Errorln("Error: Missing username in query parameters")
		return
	}

	updatedTender, err := h.srv.RollbackTenderVersion(r.Context(), tenderID, version, username)
	if err != nil {
		var statusCode int
		switch err.Error() {
		case "service.RollbackTenderVersion: repository.UpdatePartTender: %!w(<nil>)":
			statusCode = http.StatusForbidden
		case "service.RollbackTenderVersion: error fetching target version: no rows in result set":
			statusCode = http.StatusNotFound
		default:
			statusCode = http.StatusInternalServerError
		}
		errwriter.RespondWithError(w, statusCode, err.Error())
		return
	}

	response := domain.TenderResponse{
		ID:          updatedTender.ID,
		Name:        updatedTender.Name,
		Description: updatedTender.Description,
		Status:      updatedTender.Status,
		ServiceType: updatedTender.ServiceType,
		Version:     updatedTender.Version,
		CreatedAt:   updatedTender.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logger().Errorln("Error encoding JSON response:", err.Error())
		errwriter.RespondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
