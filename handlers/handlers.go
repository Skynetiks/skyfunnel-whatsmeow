package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"whatsmeow-service/config"
	"whatsmeow-service/models"
	"whatsmeow-service/services"
)

type Handlers struct {
	config  *config.Config
	service *services.WhatsAppMeowService
}

func NewHandlers(cfg *config.Config, svc *services.WhatsAppMeowService) *Handlers {
	return &Handlers{
		config:  cfg,
		service: svc,
	}
}

// SendMessage handles message sending requests
func (h *Handlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.OrganizationID == "" || req.ToJID == "" || req.MessageType == "" {
		h.sendErrorResponse(w, "Missing required fields", fmt.Errorf("organizationId, toJID, and messageType are required"), http.StatusBadRequest)
		return
	}

	// Send message via service
	messageID, err := h.service.SendMessage(req)
	if err != nil {
		h.sendErrorResponse(w, "Failed to send message", err, http.StatusInternalServerError)
		return
	}

	response := models.SendMessageResponse{
		Success:   true,
		MessageID: messageID,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetStatus handles status requests
func (h *Handlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	organizationID := r.URL.Query().Get("organizationId")
	if organizationID == "" {
		h.sendErrorResponse(w, "Organization ID is required", fmt.Errorf("organizationId parameter is required"), http.StatusBadRequest)
		return
	}

	account, err := h.service.GetAccount(organizationID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get account", err, http.StatusInternalServerError)
		return
	}

	response := models.ConnectionStatusResponse{
		Success: true,
		Account: account,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// GetQR handles QR code requests
func (h *Handlers) GetQR(w http.ResponseWriter, r *http.Request) {
	organizationID := r.URL.Query().Get("organizationId")
	if organizationID == "" {
		h.sendErrorResponse(w, "Organization ID is required", fmt.Errorf("organizationId parameter is required"), http.StatusBadRequest)
		return
	}

	qrCode, err := h.service.GetQRCode(organizationID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to get QR code", err, http.StatusInternalServerError)
		return
	}

	response := models.QRCodeResponse{
		Success: true,
		QRCode:  qrCode,
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// Connect handles connection requests
func (h *Handlers) Connect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OrganizationID string `json:"organizationId"`
		DeviceID       string `json:"deviceId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err, http.StatusBadRequest)
		return
	}

	if req.OrganizationID == "" || req.DeviceID == "" {
		h.sendErrorResponse(w, "Missing required fields", fmt.Errorf("organizationId and deviceId are required"), http.StatusBadRequest)
		return
	}

	err := h.service.Connect(req.OrganizationID, req.DeviceID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to connect", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Connection initiated",
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// Disconnect handles disconnection requests
func (h *Handlers) Disconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OrganizationID string `json:"organizationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err, http.StatusBadRequest)
		return
	}

	if req.OrganizationID == "" {
		h.sendErrorResponse(w, "Organization ID is required", fmt.Errorf("organizationId is required"), http.StatusBadRequest)
		return
	}

	err := h.service.Disconnect(req.OrganizationID)
	if err != nil {
		h.sendErrorResponse(w, "Failed to disconnect", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Disconnected successfully",
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// Health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "whatsmeow-service",
		"version":   "1.0.0",
	}

	h.sendJSONResponse(w, response, http.StatusOK)
}

// Helper methods
func (h *Handlers) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func (h *Handlers) sendErrorResponse(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("Error: %s - %v", message, err)
	
	response := models.SendMessageResponse{
		Success: false,
		Error:   fmt.Sprintf("%s: %v", message, err),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
