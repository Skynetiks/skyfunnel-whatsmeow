# WhatsApp Meow Service

A Go-based microservice that provides WhatsApp messaging capabilities using the [whatsmeow](https://github.com/tulir/whatsmeow) library. This service integrates with your main SkyFunnel application to provide free WhatsApp messaging as an alternative to the official WhatsApp Business API.

## Features

- 🚀 **Free WhatsApp Messaging** - No per-message charges
- 📱 **Personal Account Integration** - Use your personal WhatsApp account
- 🔄 **Real-time Messaging** - Send and receive messages instantly
- 📊 **Message Tracking** - Track message status and delivery
- 🛡️ **Secure** - Encrypted session management
- 🔌 **RESTful API** - Easy integration with your main application

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   SkyFunnel     │    │   Node.js API    │    │   Go Service    │
│   (Frontend)    │◄──►│   (tRPC/API)     │◄──►│   (whatsmeow)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │   PostgreSQL     │
                       │   Database       │
                       └──────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker (optional)

### Installation

1. **Clone the repository**
```bash
git clone <your-repo-url>
cd whatsmeow-service
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your database configuration
```

4. **Run the service**
```bash
go run main.go
```

### Docker Setup

```bash
# Build and run with Docker
docker-compose up -d
```

## API Endpoints

### Send Message
```http
POST /api/whatsmeow/send
Content-Type: application/json

{
  "organizationId": "org_123",
  "toJID": "1234567890@s.whatsapp.net",
  "messageType": "text",
  "messageText": "Hello from WhatsApp Meow!",
  "leadId": "lead_456"
}
```

### Get Connection Status
```http
GET /api/whatsmeow/status?organizationId=org_123
```

### Get QR Code
```http
GET /api/whatsmeow/qr?organizationId=org_123
```

### Connect Account
```http
POST /api/whatsmeow/connect
Content-Type: application/json

{
  "organizationId": "org_123",
  "deviceId": "device_456"
}
```

## Database Schema

The service uses the following database tables:

- `WhatsAppMeowAccount` - Stores account information and connection status
- `WhatsAppMeowMessage` - Stores message history and status

See `database/schema.sql` for the complete schema.

## Configuration

### Environment Variables

```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/skyfunnel

# Service
PORT=8081
LOG_LEVEL=info

# WhatsApp Meow
WHATSMEOW_SESSION_DIR=./sessions
WHATSMEOW_LOG_LEVEL=info
```

### Database Configuration

The service expects the following database tables to exist in your main SkyFunnel database:

```sql
-- See database/schema.sql for complete schema
CREATE TABLE "WhatsAppMeowAccount" (
  id VARCHAR PRIMARY KEY,
  organization_id VARCHAR NOT NULL,
  device_id VARCHAR UNIQUE NOT NULL,
  -- ... other fields
);
```

## Integration with SkyFunnel

### 1. Database Integration

Add the whatsmeow tables to your main SkyFunnel database by running the schema migration:

```sql
-- Run the schema from database/schema.sql
```

### 2. API Integration

In your SkyFunnel application, add the following API routes:

```typescript
// src/app/api/whatsmeow/send/route.ts
// src/app/api/whatsmeow/status/route.ts
// src/app/api/whatsmeow/qr/route.ts
```

### 3. tRPC Integration

Add the whatsmeow router to your tRPC configuration:

```typescript
// src/trpc/routes/whatsmeow.ts
export const whatsMeowRouter = router({
  // ... routes
});
```

### 4. Frontend Integration

Use the provider selector component to allow users to choose between WhatsApp Business API and WhatsApp Meow:

```typescript
// src/components/Dashboard/whatsapp/WhatsAppProviderSelector.tsx
```

## Message Types Supported

- **Text Messages** - Plain text messages
- **Image Messages** - Images with optional captions
- **Video Messages** - Video files with optional captions
- **Audio Messages** - Audio files
- **Document Messages** - Document files

## Security Considerations

- Session data is encrypted before storage
- Device IDs are unique per organization
- Message content is not logged
- Connection status is regularly updated

## Monitoring

The service provides the following monitoring endpoints:

- `/health` - Health check endpoint
- `/metrics` - Prometheus metrics (if enabled)
- `/status` - Service status and connection info

## Troubleshooting

### Common Issues

1. **Connection Failed**
   - Check if the device is properly paired
   - Verify database connection
   - Check session data integrity

2. **Messages Not Sending**
   - Verify WhatsApp account is connected
   - Check message format and JID format
   - Review error logs

3. **Database Errors**
   - Ensure all required tables exist
   - Check database permissions
   - Verify connection string

### Logs

The service logs important events including:
- Connection status changes
- Message send/receive events
- Error conditions
- Authentication events

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o whatsmeow-service main.go
```

### Code Generation

```bash
go generate ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue in the repository
- Check the documentation
- Review the logs for error details
