# Planvia Partner API

This is the backend API for the Planvia Partner platform built with Go, Fiber, and MongoDB.

## Prerequisites

- Go 1.21 or later
- MongoDB
- Make sure MongoDB is running on localhost:27017

## Setup

1. Clone the repository
2. Create a `.env` file in the root directory with the following content:

   ```
   MONGO_URI=mongodb://localhost:27017
   DB_NAME=planvia
   PORT=5000
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

The server will start on port 5000 (or the port specified in your .env file).

## API Endpoints

### Partner Registration

- **POST** `/api/partners/register`
- **Body:**
  ```json
  {
    "companyName": "string",
    "email": "string",
    "password": "string",
    "phoneNumber": "string",
    "address": "string",
    "city": "string",
    "businessType": "string",
    "taxNumber": "string",
    "contactPerson": "string"
  }
  ```

## Development

The project structure follows standard Go project layout:

```
.
├── cmd/
│   └── api/            # Application entrypoint
├── config/             # Configuration
├── internal/
│   ├── database/       # Database connection
│   ├── handlers/       # HTTP handlers
│   └── models/         # Data models
└── .env               # Environment variables
```
