# PostgreSQL Helper Package

A Go package for simplified PostgreSQL database operations with support for Unix timestamps.

## Features

- Connection management
- CRUD operations
- Unix timestamp handling
- Pagination and sorting
- Search functionality
- Health checks
- Basic backup operations
- Index management

## Installation

```bash
go get github.com/yourusername/postgres-helper
```

## Quick Start

```go
import "github.com/yourusername/postgres-helper/db"

// Configure database
config := db.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "your_password",
    DBName:   "your_database",
    SSLMode:  "disable",
}

// Connect
database, err := db.Connect(config)
if err != nil {
    log.Fatal(err)
}
defer database.Close()

// Create table
if err := db.CreateTable(database); err != nil {
    log.Fatal(err)
}

// Insert record
record := &db.Record{
    Name:          "Test Record",
    CreatedAtUnix: time.Now().Unix(),
}

if err := db.InsertRecord(database, record); err != nil {
    log.Fatal(err)
}
```

## Usage Examples

### Pagination
```go
opts := db.QueryOptions{
    Limit:  10,
    Offset: 0,
    SortBy: "created_at",
    Order:  "DESC",
}

records, err := db.GetRecords(database, opts)
```

### Search
```go
records, err := db.SearchRecords(database, "search term", opts)
```

### Health Check
```go
if err := db.HealthCheck(database); err != nil {
    log.Printf("Database unhealthy: %v", err)
}
```

## Contributing

Feel free to submit issues and pull requests.

## License

MIT License (or your chosen license)