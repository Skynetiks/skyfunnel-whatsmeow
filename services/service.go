package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"whatsmeow-service/config"
	"whatsmeow-service/models"
)

type WhatsAppMeowService struct {
	config *config.Config
	db     *sql.DB
	client *whatsmeow.Client
}

func NewWhatsAppMeowService(cfg *config.Config, db *sql.DB) *WhatsAppMeowService {
	return &WhatsAppMeowService{
		config: cfg,
		db:     db,
	}
}

// SendMessage sends a WhatsApp message
func (s *WhatsAppMeowService) SendMessage(req models.SendMessageRequest) (string, error) {
	// Get account for organization
	account, err := s.getAccount(req.OrganizationID)
	if err != nil {
		return "", fmt.Errorf("failed to get account: %w", err)
	}

	if !account.IsConnected {
		return "", fmt.Errorf("account is not connected")
	}

	// Initialize client if needed
	if s.client == nil {
		if err := s.initializeClient(account); err != nil {
			return "", fmt.Errorf("failed to initialize client: %w", err)
		}
	}

	// Parse JID
	toJID, err := types.ParseJID(req.ToJID)
	if err != nil {
		return "", fmt.Errorf("invalid JID: %w", err)
	}

	// Send message based on type
	var messageID string
	switch req.MessageType {
	case "text":
		messageID, err = s.sendTextMessage(toJID, req.MessageText)
	case "image":
		messageID, err = s.sendImageMessage(toJID, req.MessageText, req.MediaURL)
	case "video":
		messageID, err = s.sendVideoMessage(toJID, req.MessageText, req.MediaURL)
	case "audio":
		messageID, err = s.sendAudioMessage(toJID, req.MediaURL)
	case "document":
		messageID, err = s.sendDocumentMessage(toJID, req.MessageText, req.MediaURL)
	default:
		return "", fmt.Errorf("unsupported message type: %s", req.MessageType)
	}

	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// Save message to database
	if err := s.saveMessage(account.ID, req, messageID); err != nil {
		log.Printf("Failed to save message: %v", err)
	}

	return messageID, nil
}

// GetAccount retrieves account information
func (s *WhatsAppMeowService) GetAccount(organizationID string) (*models.WhatsAppMeowAccount, error) {
	return s.getAccount(organizationID)
}

// GetQRCode retrieves QR code for pairing
func (s *WhatsAppMeowService) GetQRCode(organizationID string) (string, error) {
	account, err := s.getAccount(organizationID)
	if err != nil {
		return "", err
	}

	if account.QRCode == nil {
		return "", fmt.Errorf("QR code not available")
	}

	return *account.QRCode, nil
}

// Connect initiates connection
func (s *WhatsAppMeowService) Connect(organizationID, deviceID string) error {
	// Update account status
	_, err := s.db.Exec(`
		UPDATE "WhatsAppMeowAccount" 
		SET connection_status = 'CONNECTING', updated_at = $1 
		WHERE organization_id = $2
	`, time.Now(), organizationID)
	
	if err != nil {
		return fmt.Errorf("failed to update connection status: %w", err)
	}

	// Initialize client and start connection process
	account, err := s.getAccount(organizationID)
	if err != nil {
		return err
	}

	return s.initializeClient(account)
}

// Disconnect disconnects the client
func (s *WhatsAppMeowService) Disconnect(organizationID string) error {
	if s.client != nil {
		s.client.Disconnect()
		s.client = nil
	}

	// Update account status
	_, err := s.db.Exec(`
		UPDATE "WhatsAppMeowAccount" 
		SET connection_status = 'DISCONNECTED', is_connected = false, updated_at = $1 
		WHERE organization_id = $2
	`, time.Now(), organizationID)
	
	return err
}

// Private methods
func (s *WhatsAppMeowService) getAccount(organizationID string) (*models.WhatsAppMeowAccount, error) {
	query := `
		SELECT id, organization_id, device_id, session_data, qr_code, is_connected, is_paired, 
		       phone_number, display_name, profile_picture, last_seen, connection_status, 
		       created_at, updated_at
		FROM "WhatsAppMeowAccount" 
		WHERE organization_id = $1
	`
	
	var account models.WhatsAppMeowAccount
	var sessionDataJSON sql.NullString
	var qrCode sql.NullString
	var phoneNumber sql.NullString
	var displayName sql.NullString
	var profilePicture sql.NullString
	var lastSeen sql.NullTime
	
	err := s.db.QueryRow(query, organizationID).Scan(
		&account.ID,
		&account.OrganizationID,
		&account.DeviceID,
		&sessionDataJSON,
		&qrCode,
		&account.IsConnected,
		&account.IsPaired,
		&phoneNumber,
		&displayName,
		&profilePicture,
		&lastSeen,
		&account.ConnectionStatus,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if sessionDataJSON.Valid {
		account.SessionData = &models.SessionData{}
		// Parse session data JSON here
	}
	if qrCode.Valid {
		account.QRCode = &qrCode.String
	}
	if phoneNumber.Valid {
		account.PhoneNumber = &phoneNumber.String
	}
	if displayName.Valid {
		account.DisplayName = &displayName.String
	}
	if profilePicture.Valid {
		account.ProfilePicture = &profilePicture.String
	}
	if lastSeen.Valid {
		account.LastSeen = &lastSeen.Time
	}

	return &account, nil
}

func (s *WhatsAppMeowService) initializeClient(account *models.WhatsAppMeowAccount) error {
	// Initialize device store
	deviceStore, err := sqlstore.New("postgres", s.config.DatabaseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create device store: %w", err)
	}

	// Get or create device
	device, err := deviceStore.GetDevice(account.DeviceID)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Create client
	client := whatsmeow.NewClient(device, nil)
	
	// Set up event handlers
	client.AddEventHandler(s.eventHandler)

	// Connect
	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	s.client = client
	return nil
}

func (s *WhatsAppMeowService) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		s.handleIncomingMessage(v)
	case *events.Connected:
		s.handleConnected()
	case *events.Disconnected:
		s.handleDisconnected()
	case *events.LoggedOut:
		s.handleLoggedOut()
	case *events.QR:
		s.handleQRCode(v)
	}
}

func (s *WhatsAppMeowService) handleIncomingMessage(msg *events.Message) {
	log.Printf("Received message from %s: %s", msg.Info.Sender, msg.Message.GetConversation())
	// Handle incoming message logic here
}

func (s *WhatsAppMeowService) handleConnected() {
	log.Println("Connected to WhatsApp")
	// Update connection status in database
}

func (s *WhatsAppMeowService) handleDisconnected() {
	log.Println("Disconnected from WhatsApp")
	// Update connection status in database
}

func (s *WhatsAppMeowService) handleLoggedOut() {
	log.Println("Logged out from WhatsApp")
	// Handle logout logic
}

func (s *WhatsAppMeowService) handleQRCode(qr *events.QR) {
	log.Println("QR code received")
	// Save QR code to database
}

func (s *WhatsAppMeowService) sendTextMessage(toJID types.JID, text string) (string, error) {
	msg := &whatsmeow.TextMessage{
		Text: text,
	}
	
	resp, err := s.client.SendMessage(context.Background(), toJID, msg)
	if err != nil {
		return "", err
	}
	
	return resp.ID, nil
}

func (s *WhatsAppMeowService) sendImageMessage(toJID types.JID, caption, mediaURL string) (string, error) {
	// Implementation for image messages
	// This would involve downloading the media and uploading to WhatsApp
	return "", fmt.Errorf("image messages not yet implemented")
}

func (s *WhatsAppMeowService) sendVideoMessage(toJID types.JID, caption, mediaURL string) (string, error) {
	// Implementation for video messages
	return "", fmt.Errorf("video messages not yet implemented")
}

func (s *WhatsAppMeowService) sendAudioMessage(toJID types.JID, mediaURL string) (string, error) {
	// Implementation for audio messages
	return "", fmt.Errorf("audio messages not yet implemented")
}

func (s *WhatsAppMeowService) sendDocumentMessage(toJID types.JID, caption, mediaURL string) (string, error) {
	// Implementation for document messages
	return "", fmt.Errorf("document messages not yet implemented")
}

func (s *WhatsAppMeowService) saveMessage(accountID string, req models.SendMessageRequest, messageID string) error {
	query := `
		INSERT INTO "WhatsAppMeowMessage" 
		(whats_app_meow_account_id, message_id, lead_id, from_jid, to_jid, message_type, message_text, is_sent, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := s.db.Exec(query, 
		accountID, 
		messageID, 
		req.LeadID, 
		"", // fromJID - will be set by the client
		req.ToJID, 
		req.MessageType, 
		req.MessageText, 
		true, 
		time.Now(),
	)
	
	return err
}
