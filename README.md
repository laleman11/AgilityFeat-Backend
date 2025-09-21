# AgilityFeat-Backend

This repository contains a starter Go project that follows a hexagonal architecture. The application exposes HTTP endpoints while keeping domain logic isolated from adapters.

## Requirements

- Go 1.22 or newer

## Project Layout

- `cmd/api`: entry-point and dependency composition.
- `internal/core`: pure domain logic.
- `internal/app`: use-cases and orchestration.
- `internal/adapter/http`: primary HTTP adapter.
- `internal/infra/memory`: in-memory adapters (demo persistence for underwriting history).
- `internal/port`: contracts that connect adapters and use-cases.

## Run It

```bash
go run ./cmd/api
```

Available endpoints:

- `GET /api/v1/ping` → `{"message":"pong"}`
- `POST /api/v1/underwriting`
  ```json
  {
    "user_id": "user-123",
    "monthly_income": 8000,
    "monthly_debts": 2000,
    "loan_amount": 250000,
    "property_value": 350000,
    "credit_score": 760,
    "occupancy_type": "primary"
  }
  ```
  Response example:
  ```json
  {
    "decision": "Approve",
    "dti": 0.25,
    "ltv": 0.71,
    "reasons": [
      "Meets all underwriting rules"
    ]
  }
  ```
- `GET /api/v1/underwriting/history/{user_id}` →
  ```json
  {
    "items": [
      {
        "user_id": "user-123",
        "request": {
          "user_id": "user-123",
          "monthly_income": 8000,
          "monthly_debts": 2000,
          "loan_amount": 250000,
          "property_value": 350000,
          "credit_score": 760,
          "occupancy_type": "primary"
        },
        "response": {
          "decision": "Approve",
          "dti": 0.25,
          "ltv": 0.71,
          "reasons": ["Meets automatic approval criteria"]
        }
      }
    ]
  }
  ```

## Tests

```bash
go test ./...
```
