# Margin Requirement Calculator

This repository contains a simplified "initial margin preview" flow for a crypto derivatives exchange. It consists of a Go backend and a Next.js frontend.

## Project Structure

- `backend/`: Go API server (handles asset config and margin validation).
- `frontend/`: Next.js React application (UI for margin calculation).

## Prerequisites

- **Go**: Version 1.20 or higher.
- **Node.js**: Version 18 or higher (for Frontend).
- **npm**: Included with Node.js.

## Setup & Running Locally

### 1. Backend

The backend serves asset configuration and validates margin calculations.

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Run the server:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

### 2. Frontend

The frontend is a Next.js application that interacts with the backend.

1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Run the development server:
    ```bash
    npm run dev
    ```


## API Documentation

### 1. GET /config/assets
Returns metadata for supported assets.

**Response:**
```json
{
  "assets": [
    {
      "symbol": "BTC",
      "mark_price": 62000,
      "contract_value": 0.001,
      "allowed_leverage": [5, 10, 20, 50, 100]
    },
    ...
  ]
}
```

### 2. POST /margin/validate
Validates the margin calculated by the client.

**Request:**
```json
{
  "asset": "BTC",
  "order_size": 2,
  "side": "long",
  "leverage": 20,
  "margin_client": 6.2
}
```

**Response (Success):**
```json
{
  "status": "ok",
  "margin_required": 6.2
}
```

**Response (Error):**
```json
{
  "status": "error",
  "message": "Incorrect margin calculation",
  "margin_required": 6.25
}
```

## Assumptions

- **Margin Formula**: `(mark_price * order_size * contract_value) / leverage`
- **Rounding**: Margin is rounded to 2 decimal places.
- **Validation**: The backend validates the client's calculation within a small tolerance.
