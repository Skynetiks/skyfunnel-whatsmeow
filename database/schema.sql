-- WhatsApp Meow Service Database Schema
-- Add these tables to your main SkyFunnel database

-- WhatsApp Meow Account Management
CREATE TABLE IF NOT EXISTS "WhatsAppMeowAccount" (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    organization_id VARCHAR(255) NOT NULL,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    session_data JSONB,
    qr_code TEXT,
    is_connected BOOLEAN DEFAULT false,
    is_paired BOOLEAN DEFAULT false,
    phone_number VARCHAR(20),
    display_name VARCHAR(255),
    profile_picture TEXT,
    last_seen TIMESTAMP,
    connection_status VARCHAR(20) DEFAULT 'DISCONNECTED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES "Organization"(id) ON DELETE CASCADE
);

-- WhatsApp Meow Message Management
CREATE TABLE IF NOT EXISTS "WhatsAppMeowMessage" (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    whats_app_meow_account_id VARCHAR(255) NOT NULL,
    message_id VARCHAR(255) UNIQUE NOT NULL,
    lead_id VARCHAR(255),
    from_jid VARCHAR(255) NOT NULL,
    to_jid VARCHAR(255) NOT NULL,
    message_type VARCHAR(20) NOT NULL,
    message_text TEXT,
    media_url TEXT,
    media_type VARCHAR(50),
    is_sent BOOLEAN DEFAULT false,
    is_delivered BOOLEAN DEFAULT false,
    is_read BOOLEAN DEFAULT false,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    error_code VARCHAR(50),
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    
    CONSTRAINT fk_account FOREIGN KEY (whats_app_meow_account_id) REFERENCES "WhatsAppMeowAccount"(id) ON DELETE CASCADE,
    CONSTRAINT fk_lead FOREIGN KEY (lead_id) REFERENCES "Lead"(id) ON DELETE SET NULL
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_whatsmeow_account_org ON "WhatsAppMeowAccount"(organization_id);
CREATE INDEX IF NOT EXISTS idx_whatsmeow_account_device ON "WhatsAppMeowAccount"(device_id);
CREATE INDEX IF NOT EXISTS idx_whatsmeow_message_account ON "WhatsAppMeowMessage"(whats_app_meow_account_id);
CREATE INDEX IF NOT EXISTS idx_whatsmeow_message_lead ON "WhatsAppMeowMessage"(lead_id);
CREATE INDEX IF NOT EXISTS idx_whatsmeow_message_from ON "WhatsAppMeowMessage"(from_jid);
CREATE INDEX IF NOT EXISTS idx_whatsmeow_message_to ON "WhatsAppMeowMessage"(to_jid);
CREATE INDEX IF NOT EXISTS idx_whatsmeow_message_timestamp ON "WhatsAppMeowMessage"(timestamp);

-- Enums (if your database supports them)
-- For PostgreSQL, you can create these as custom types
DO $$ BEGIN
    CREATE TYPE whatsapp_meow_connection_status AS ENUM (
        'DISCONNECTED',
        'CONNECTING', 
        'CONNECTED',
        'PAIRING',
        'PAIRED',
        'ERROR'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE whatsapp_meow_message_type AS ENUM (
        'TEXT',
        'IMAGE',
        'VIDEO',
        'AUDIO',
        'DOCUMENT',
        'STICKER',
        'LOCATION',
        'CONTACT',
        'SYSTEM'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Update the tables to use the enums
ALTER TABLE "WhatsAppMeowAccount" 
ALTER COLUMN connection_status TYPE whatsapp_meow_connection_status 
USING connection_status::whatsapp_meow_connection_status;

ALTER TABLE "WhatsAppMeowMessage" 
ALTER COLUMN message_type TYPE whatsapp_meow_message_type 
USING message_type::whatsapp_meow_message_type;

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_whatsmeow_account_updated_at 
    BEFORE UPDATE ON "WhatsAppMeowAccount" 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
