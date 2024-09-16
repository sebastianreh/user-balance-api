# User Balance API

## **INFO**: In order to run the project an SMTP_PASSWORD environment configuration is required, this was provided via email.

## Overview

The User Balance API is a backend system for managing users and their financial transactions. It includes endpoints for
creating users, managing their financial transactions, checking their current balance, and processing CSV files for bulk
user transactions migration.
It was designed using Domain Driven Design, Hexagonal Design and Clean Architecture principles in mind.
---

## Key Features

- **User Management**: Create, update, delete, and fetch user information.
- **Transaction Handling**: Allows for creation, update, and deletion of transactions.
- **Balance Inquiry**: Fetch the current balance for a user, with optional date range filters.
- **CSV-Based Migration**: Upload CSV files to process bulk user transaction data and generate migration reports.
- **Email Notifications**: Sends a migration report via email to specified recipients.

---

## Endpoints

### User Endpoints

- `/users/create`: Create a new user (POST request with user data in JSON).
- `/users/:id`: Get user details by ID (GET), update user (PUT), delete user (DELETE).
- `/users/:user_id/balance`: Get user balance, with optional `from` and `to` date filters for balance calculation (GET).

### Transaction Endpoints

- `/transactions/create`: Create a new transaction for a user (POST request with transaction data in JSON).
- `/transactions/:id`: Get transaction by ID (GET), update transaction (PUT), delete transaction (DELETE).

### Migration Endpoints

- `/migrate`: Upload a CSV file to process bulk transactions and generate a migration report (POST request with CSV
  file).

---

## Setup Guide

### Prerequisites

- Docker and Docker Compose
- Go (v1.23 or later)

### Step-by-Step Guide

1. **Build Docker Images**:

   This step will create Docker images for the HTTP service.
   Run:
   ```bash
   make build-server-image
   ```

2. **Start the Server**:

   Run Docker Compose to start the services.
   ```bash
   make start-compose
   ```

3. **Shut Down Services (Optional)**:
   ```bash
   make down-compose
   ```

---

## Example Request

### Creating a User

```http
POST /user-balance-api/users/create
```

### JSON Body:

```json
{
  "first_name": "Sebastian",
  "last_name": "Reh",
  "email": "sebastianreh@example.com"
}
```

### Example Response:

```json
{
  "user_id": "12345"
}
```

---

### Migrate Transactions

```http
POST /user-balance-api/migrate
```

### File:

```
input_data.csv
```

### Header

``` 
"X-User-Emails": sebastian.reh@gmail.com, test@example.com
```

### Example Response:

```
200 OK
```

---

## Testing

To run the tests, use this command:

```bash
make run-test
```

Tests cover unit and integration tests for the user, transaction, and migration services.

---

## Generate Testing data

To create test data, use this command, please remember to run the server before using it:

```bash
make create-migration-csv
```

#### 1. It Creates random users using the API

#### 2. It creates two csv file in directory `scripts/generate_transactions/files`

- **[input_data.csv](scripts/generate_transactions/files/input_data.csv)** Input csv file to test the migration process.
- **[expected_output_data.csv](scripts/generate_transactions/files/expected_output_data.csv)**: Expected balance output
  for every user.
- **[output_count_file.csv](scripts/generate_transactions/files/output_count_file.csv)**: Counts how many
  transactions to avoid transaction IDs conflicts on every migration

---

### Make commands:

- `code-format-check`: prints unformatted files.

- `lint-install`: runs installs go linter.

- `lint`: runs go linter.

- `test`: runs go tests with `go test ./...`

- `create-swag-docs:` creates all swagger required docs

- `build-server-image`: builds Dockerfile image

- `start-compose`: runs the docker-compose file detached

- `down-compose`: stops the docker-compose containers

---

## API Documentation

Swagger is integrated for API documentation. To view the API documentation, visit the
`user-balance-api/swagger/index.html` endpoint after running the server.
