package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// WhatsAppMeowAccount represents a WhatsApp Meow account
type WhatsAppMeowAccount struct {
	ID               string                     `json:"id" db:"id"`
	OrganizationID   string                     `json:"organizationId" db:"organization_id"`
	DeviceID         string                     `json:"deviceId" db:"device_id"`
	SessionData      *SessionData               `json:"sessionData,omitempty" db:"session_data"`
	QRCode           *string                    `json:"qrCode,omitempty" db:"qr_code"`
	IsConnected      bool                       `json:"isConnected" db:"is_connected"`
	IsPaired         bool                       `json:"isPaired" db:"is_paired"`
	PhoneNumber      *string                    `json:"phoneNumber,omitempty" db:"phone_number"`
	DisplayName      *string                    `json:"displayName,omitempty" db:"display_name"`
	ProfilePicture   *string                    `json:"profilePicture,omitempty" db:"profile_picture"`
	LastSeen         *time.Time                 `json:"lastSeen,omitempty" db:"last_seen"`
	ConnectionStatus WhatsAppMeowConnectionStatus `json:"connectionStatus" db:"connection_status"`
	CreatedAt        time.Time                  `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time                  `json:"updatedAt" db:"updated_at"`
}

// WhatsAppMeowMessage represents a WhatsApp Meow message
type WhatsAppMeowMessage struct {
	ID                    string                    `json:"id" db:"id"`
	WhatsAppMeowAccountID   string                    `json:"whatsAppMeowAccountId" db:"whats_app_meow_account_id"`
	MessageID             string                    `json:"messageId" db:"message_id"`
	LeadID                *string                   `json:"leadId,omitempty" db:"lead_id"`
	FromJID               string                    `json:"fromJID" db:"from_jid"`
	ToJID                 string                    `json:"toJID" db:"to_jid"`
	MessageType           WhatsAppMeowMessageType   `json:"messageType" db:"message_type"`
	MessageText           *string                   `json:"messageText,omitempty" db:"message_text"`
	MediaURL              *string                   `json:"mediaUrl,omitempty" db:"media_url"`
	MediaType             *string                   `json:"mediaType,omitempty" db:"media_type"`
	IsSent                bool                      `json:"isSent" db:"is_sent"`
	IsDelivered           bool                      `json:"isDelivered" db:"is_delivered"`
	IsRead                bool                      `json:"isRead" db:"is_read"`
	Timestamp             time.Time                 `json:"timestamp" db:"timestamp"`
	SentAt                *time.Time                `json:"sentAt,omitempty" db:"sent_at"`
	DeliveredAt           *time.Time                `json:"deliveredAt,omitempty" db:"delivered_at"`
	ReadAt                *time.Time                `json:"readAt,omitempty" db:"read_at"`
	ErrorCode             *string                   `json:"errorCode,omitempty" db:"error_code"`
	ErrorMessage          *string                   `json:"errorMessage,omitempty" db:"error_message"`
	RetryCount            int                       `json:"retryCount" db:"retry_count"`
}

// SessionData represents encrypted session data
type SessionData struct {
	DeviceID    string                 `json:"deviceId"`
	SessionData map[string]interface{} `json:"sessionData"`
	Encrypted   bool                   `json:"encrypted"`
}

// Value implements driver.Valuer for database storage
func (sd *SessionData) Value() (driver.Value, error) {
	if sd == nil {
		return nil, nil
	}
	return json.Marshal(sd)
}

// Scan implements sql.Scanner for database retrieval
func (sd *SessionData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, sd)
}

// WhatsAppMeowConnectionStatus represents connection status
type WhatsAppMeowConnectionStatus string

const (
	ConnectionStatusDisconnected WhatsAppMeowConnectionStatus = "DISCONNECTED"
	ConnectionStatusConnecting    WhatsAppMeowConnectionStatus = "CONNECTING"
	ConnectionStatusConnected     WhatsAppMeowConnectionStatus = "CONNECTED"
	ConnectionStatusPairing       WhatsAppMeowConnectionStatus = "PAIRING"
	ConnectionStatusPaired        WhatsAppMeowConnectionStatus = "PAIRED"
	ConnectionStatusError         WhatsAppMeowConnectionStatus = "ERROR"
)

// WhatsAppMeowMessageType represents message types
type WhatsAppMeowMessageType string

const (
	MessageTypeText     WhatsAppMeowMessageType = "TEXT"
	MessageTypeImage    WhatsAppMeowMessageType = "IMAGE"
	MessageTypeVideo    WhatsAppMeowMessageType = "VIDEO"
	MessageTypeAudio    WhatsAppMeowMessageType = "AUDIO"
	MessageTypeDocument WhatsAppMeowMessageType = "DOCUMENT"
	MessageTypeSticker  WhatsAppMeowMessageType = "STICKER"
	MessageTypeLocation WhatsAppMeowMessageType = "LOCATION"
	MessageTypeContact  WhatsAppMeowMessageType = "CONTACT"
	MessageTypeSystem   WhatsAppMeowMessageType = "SYSTEM"
)

// Request/Response types
type SendMessageRequest struct {
	OrganizationID string `json:"organizationId"`
	ToJID          string `json:"toJID"`
	MessageType    string `json:"messageType"`
	MessageText    string `json:"messageText,omitempty"`
	MediaURL       string `json:"mediaUrl,omitempty"`
	MediaType      string `json:"mediaType,omitempty"`
	LeadID         string `json:"leadId,omitempty"`
}

type SendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

type ConnectionStatusResponse struct {
	Success bool                    `json:"success"`
	Account *WhatsAppMeowAccount    `json:"account,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

type QRCodeResponse struct {
	Success bool   `json:"success"`
	QRCode  string `json:"qrCode,omitempty"`
	Error   string `json:"error,omitempty"`
}
