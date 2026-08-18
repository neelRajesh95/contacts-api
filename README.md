# Contacts API

A small REST API for managing contacts. This project was built as part of a technical coding exercise using Go and PostgreSQL.

## Tech Stack

* Go
* PostgreSQL
* Chi router
* pgx PostgreSQL driver
* UUID

## APIs

### Create Contact

`POST /contacts`

Required header:

```text
X-Tenant-ID: tenant-123
```

Request body:

```json
{
  "first_name": "John",
  "last_name": "Smith",
  "email": "john@example.com",
  "company": "Acme"
}
```

The contact is created with `pending` status.

---

### Get Contacts

`GET /contacts`

Required header:

```text
X-Tenant-ID: tenant-123
```

The API returns contacts for the given tenant only.

---

### Enrich Contact

`POST /contacts/:id/enrich`

Required header:

```text
X-Tenant-ID: tenant-123
```

For this exercise, the enrichment is simulated by changing the contact status from `pending` to `enriched`.

The tenant ID is also checked while updating the contact so that a contact belonging to another tenant cannot be updated.

---

## Database

The application uses PostgreSQL.

Create a database called:

```text
contacts
```

Then run the SQL file:

```text
migrations/001_create_contacts.sql
```

The contacts table contains:

* id
* tenant_id
* first_name
* last_name
* email
* company
* status
* created_at
* updated_at

An index is added on `tenant_id` since all contact queries are tenant based.

## Configuration

Set the following environment variables:

```text
PORT=8080

DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/contacts?sslmode=disable
```

Example on PowerShell:

```powershell
$env:PORT="8080"
$env:DATABASE_URL="postgres://postgres:password@localhost:5432/contacts?sslmode=disable"
```

## Run Locally

Install the Go dependencies:

```bash
go mod tidy
```

Make sure PostgreSQL is running and the database/migration has been created.

Then start the application:

```bash
go run ./cmd/server
```

The API will be available at:

```text
http://localhost:8080
```

Health check:

```text
GET /health
```

Expected response:

```json
{
  "status": "ok"
}
```

## Testing

Run the tests with:

```bash
go test ./...
```

## Project Structure

```text
cmd/
  server/
    main.go

internal/
  database/
    postgres.go
  handler/
    contact_handler.go
  model/
    contact.go
  repository/
    contact_repository.go
  service/
    contact_service.go

migrations/
  001_create_contacts.sql
```

The HTTP handlers deal with requests and responses, the service layer contains validation/business logic, and the repository handles PostgreSQL queries.

