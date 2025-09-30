# WhatsApp Meow Integration Setup Instructions

This document provides step-by-step instructions for integrating the WhatsApp Meow service with your SkyFunnel application.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Node.js and npm/pnpm
- Docker (optional)

## Step 1: Database Setup

### 1.1 Add Database Schema

Run the following SQL commands in your SkyFunnel database:

```sql
-- Copy and run the contents of database/schema.sql
-- This will create the required tables for WhatsApp Meow
```

### 1.2 Update Prisma Schema

Add the following to your `prisma/schema.prisma` file:

```prisma
// WhatsApp Meow Account Management
model WhatsAppMeowAccount {
  id                String                     @id @default(cuid())
  organizationId    String
  organization      Organization               @relation(fields: [organizationId], references: [id], onDelete: Cascade)
  deviceId          String                     @unique
  sessionData       Json?
  qrCode            String?
  isConnected       Boolean                    @default(false)
  isPaired          Boolean                    @default(false)
  phoneNumber       String?
  displayName       String?
  profilePicture    String?
  lastSeen          DateTime?
  connectionStatus  WhatsAppMeowConnectionStatus @default(DISCONNECTED)
  createdAt         DateTime                   @default(now())
  updatedAt         DateTime                   @updatedAt
  
  whatsAppMeowMessages WhatsAppMeowMessage[]
  
  @@unique([organizationId])
}

// WhatsApp Meow Message Management
model WhatsAppMeowMessage {
  id                    String                @id @default(cuid())
  whatsAppMeowAccountId String
  whatsAppMeowAccount   WhatsAppMeowAccount   @relation(fields: [whatsAppMeowAccountId], references: [id], onDelete: Cascade)
  messageId             String                @unique
  leadId                String?
  lead                  Lead?                 @relation(fields: [leadId], references: [id], onDelete: SetNull)
  fromJID               String
  toJID                 String
  messageType           WhatsAppMeowMessageType
  messageText           String?
  mediaUrl              String?
  mediaType             String?
  isSent                Boolean               @default(false)
  isDelivered           Boolean                @default(false)
  isRead                Boolean                @default(false)
  timestamp             DateTime               @default(now())
  sentAt                DateTime?
  deliveredAt           DateTime?
  readAt                DateTime?
  errorCode             String?
  errorMessage          String?
  retryCount            Int                    @default(0)
  
  @@index([whatsAppMeowAccountId])
  @@index([leadId])
  @@index([fromJID])
  @@index([toJID])
}

// Enums
enum WhatsAppMeowConnectionStatus {
  DISCONNECTED
  CONNECTING
  CONNECTED
  PAIRING
  PAIRED
  ERROR
}

enum WhatsAppMeowMessageType {
  TEXT
  IMAGE
  VIDEO
  AUDIO
  DOCUMENT
  STICKER
  LOCATION
  CONTACT
  SYSTEM
}
```

### 1.3 Update Organization and Lead Models

Add these relations to your existing models:

```prisma
// Add to Organization model
model Organization {
  // ... existing fields
  whatsAppMeowAccounts WhatsAppMeowAccount[]
}

// Add to Lead model  
model Lead {
  // ... existing fields
  whatsAppMeowMessages WhatsAppMeowMessage[]
}
```

### 1.4 Run Database Migration

```bash
npx prisma db push
# or
npx prisma migrate dev --name add-whatsmeow-support
```

## Step 2: Backend Integration

### 2.1 Add API Routes

Create the following files in your SkyFunnel project:

```bash
# Create API routes
mkdir -p src/app/api/whatsmeow
```

Copy the contents from `integration/skyfunnel-api-routes.ts` to:
- `src/app/api/whatsmeow/send/route.ts`
- `src/app/api/whatsmeow/status/route.ts`
- `src/app/api/whatsmeow/qr/route.ts`
- `src/app/api/whatsmeow/connect/route.ts`
- `src/app/api/whatsmeow/disconnect/route.ts`

### 2.2 Add tRPC Routes

Copy the contents from `integration/skyfunnel-trpc-routes.ts` to:
- `src/trpc/routes/whatsmeow.ts`

### 2.3 Update tRPC Router

Add the whatsmeow router to your main tRPC router:

```typescript
// src/trpc/index.ts
import { whatsMeowRouter } from "./routes/whatsmeow";

export const appRouter = router({
  // ... existing routers
  whatsMeow: whatsMeowRouter,
});
```

## Step 3: Frontend Integration

### 3.1 Add Provider Selector Component

Copy the contents from `integration/skyfunnel-frontend-component.tsx` to:
- `src/components/Dashboard/whatsapp/WhatsAppProviderSelector.tsx`

### 3.2 Update WhatsApp Dashboard

Update your WhatsApp dashboard to include the provider selector:

```typescript
// src/components/Dashboard/whatsapp/whatsAppComponent.tsx
import WhatsAppProviderSelector from "./WhatsAppProviderSelector";

// Add the provider selector to your dashboard
<WhatsAppProviderSelector 
  organizationId={organizationId}
  onProviderChange={(provider) => {
    // Handle provider change
    console.log("Selected provider:", provider);
  }}
/>
```

## Step 4: WhatsApp Meow Service Setup

### 4.1 Environment Variables

Add these environment variables to your `.env` file:

```env
# WhatsApp Meow Service
WHATSMEOW_SERVICE_URL=http://localhost:8081
```

### 4.2 Run WhatsApp Meow Service

#### Option A: Using Docker

```bash
# Build and run with Docker Compose
docker-compose -f docker-compose.whatsmeow.yml up -d
```

#### Option B: Manual Setup

```bash
# Navigate to whatsmeow-service directory
cd whatsmeow-service

# Install dependencies
go mod download

# Set environment variables
cp env.example .env
# Edit .env with your database configuration

# Run the service
go run main.go
```

### 4.3 Verify Service is Running

```bash
# Check if service is running
curl http://localhost:8081/health

# Expected response:
# {"status":"healthy","timestamp":"2024-01-01T00:00:00Z","service":"whatsmeow-service","version":"1.0.0"}
```

## Step 5: Testing the Integration

### 5.1 Create WhatsApp Meow Account

1. Go to your WhatsApp dashboard
2. Select "WhatsApp Meow" as the provider
3. Click "Setup WhatsApp Meow"
4. This will create a WhatsApp Meow account in your database

### 5.2 Connect WhatsApp Account

1. The service will generate a QR code
2. Scan the QR code with your personal WhatsApp account
3. Once connected, you can start sending messages

### 5.3 Send Test Message

```bash
# Send a test message
curl -X POST http://localhost:8081/api/whatsmeow/send \
  -H "Content-Type: application/json" \
  -d '{
    "organizationId": "your-org-id",
    "toJID": "1234567890@s.whatsapp.net",
    "messageType": "text",
    "messageText": "Hello from WhatsApp Meow!"
  }'
```

## Step 6: Production Deployment

### 6.1 Docker Deployment

```bash
# Build production image
docker build -t whatsmeow-service:latest ./whatsmeow-service

# Run in production
docker run -d \
  --name whatsmeow-service \
  -p 8081:8081 \
  -e DATABASE_URL="your-production-db-url" \
  whatsmeow-service:latest
```

### 6.2 Environment Configuration

Update your production environment variables:

```env
# Production environment
WHATSMEOW_SERVICE_URL=https://your-whatsmeow-service.com
```

## Troubleshooting

### Common Issues

1. **Service not starting**
   - Check database connection
   - Verify environment variables
   - Check port availability

2. **Database connection errors**
   - Ensure PostgreSQL is running
   - Check database URL format
   - Verify database permissions

3. **WhatsApp connection issues**
   - Check if QR code is generated
   - Verify device pairing
   - Check session data integrity

### Logs

Check service logs for debugging:

```bash
# Docker logs
docker logs whatsmeow-service

# Manual service logs
# Logs will be printed to stdout when running manually
```

## Support

For issues and questions:
1. Check the service logs
2. Verify database connectivity
3. Test API endpoints manually
4. Review the README.md in the whatsmeow-service directory

## Security Considerations

- Session data is encrypted before storage
- Device IDs are unique per organization
- Message content is not logged
- Use HTTPS in production
- Implement proper authentication for API endpoints
