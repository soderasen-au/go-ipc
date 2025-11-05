# go-ipc

A lightweight Go utility library providing reusable modules for inter-process communication, messaging, database abstraction, and high availability.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

**go-ipc** is a collection of independent, composable packages designed to simplify common tasks in Go applications:

- **`msn`** - Email/Messaging system with SMTP support
- **`db`** - Multi-database abstraction layer built on GORM
- **`ha`** - High Availability peer monitoring and failover management

## Features

### 📧 Messaging (`msn`)

A comprehensive email sending module with:

- **SMTP Support**: Full-featured SMTP client with multiple encryption methods
- **Connection Pooling**: Keep-alive support with configurable timeouts
- **Async/Sync Modes**: Support for both synchronous and asynchronous email sending
- **Rich Email Features**:
  - HTML email support
  - Multiple recipients (To, Cc, Bcc) with automatic duplicate filtering
  - File attachments
  - Configurable NOOP heartbeat for long-lived connections
- **Popular Providers**: Pre-tested configurations for Gmail, Outlook/Office365
- **Flexible Architecture**: Agent pattern supporting multiple message senders

### 🗄️ Database (`db`)

A simple but effective database abstraction layer:

- **Multi-Dialect Support**: MySQL, PostgreSQL, SQL Server, SQLite
- **GORM Integration**: Built on top of the powerful GORM ORM
- **Unified Configuration**: Single interface for all database types
- **Connection String Examples**: Documented DSN formats for each dialect

### 🔄 High Availability (`ha`)

Peer monitoring system for active-passive failover scenarios:

- **Automatic Failover**: Switches to SERVING mode when peer is down
- **Health Checking**: Periodic HTTP-based health checks
- **Configurable Retry Logic**: Transient failure handling
- **Two-Mode System**:
  - `WAITING` - Passive/standby mode
  - `SERVING` - Active/primary mode
- **TLS Support**: Optional certificate verification

## Installation

```bash
go get github.com/soderasen-au/go-ipc
```

## Quick Start

### Email/Messaging

```go
package main

import (
    "github.com/soderasen-au/go-ipc/msn"
    "github.com/soderasen-au/go-common/util"
)

func main() {
    // Configure email server
    cfg := msn.Config{
        Email: &msn.EmailServerConfig{
            ServerType: msn.SMTP,
            Host:       util.Ptr("smtp.gmail.com"),
            Port:       util.Ptr(587),
            Username:   util.Ptr("your-email@gmail.com"),
            Password:   util.Ptr("your-app-password"),
            Encryption: util.Ptr(msn.STARTTLS),
        },
    }

    // Create agent
    agent, res := msn.NewAgent(cfg)
    if res != nil {
        panic(res)
    }

    // Send email
    message := msn.Message{
        From:  "your-email@gmail.com",
        To:    []string{"recipient@example.com"},
        Title: "Hello from go-ipc",
        Body:  "<h1>Hello World!</h1><p>This is a test email.</p>",
    }

    if res := agent.Send(message); res != nil {
        panic(res)
    }
}
```

#### Email Configuration Options

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `ServerType` | `EmailServerType` | Currently supports `"smtp"` | Yes |
| `Host` | `*string` | SMTP server hostname | Yes |
| `Port` | `*int` | SMTP port (default: 587) | No |
| `Username` | `*string` | SMTP authentication username | Yes |
| `Password` | `*string` | SMTP authentication password | Yes |
| `Encryption` | `*EmailServerEncryption` | `"none"`, `"ssl/tls"`, or `"starttls"` | No |
| `KeepAlive` | `*bool` | Enable connection pooling | No |
| `KeepAliveTimeout` | `*int` | Connection reuse timeout (seconds) | No |
| `NoopTimeout` | `*int` | NOOP heartbeat interval (seconds) | No |

#### Popular Email Provider Configurations

**Gmail:**
```go
&msn.EmailServerConfig{
    ServerType: msn.SMTP,
    Host:       util.Ptr("smtp.gmail.com"),
    Port:       util.Ptr(587),
    Username:   util.Ptr("your-email@gmail.com"),
    Password:   util.Ptr("your-app-password"),  // Use App Password, not regular password
    Encryption: util.Ptr(msn.STARTTLS),
}
```

**Outlook/Office 365:**
```go
&msn.EmailServerConfig{
    ServerType: msn.SMTP,
    Host:       util.Ptr("smtp-mail.outlook.com"),
    Port:       util.Ptr(587),
    Username:   util.Ptr("your-email@outlook.com"),
    Password:   util.Ptr("your-password"),
    Encryption: util.Ptr(msn.STARTTLS),
}
```

### Database

```go
package main

import (
    "github.com/soderasen-au/go-ipc/db"
)

func main() {
    // MySQL example
    cfg := db.Config{
        Dialect: "mysql",
        DSN:     "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
    }

    database, err := db.NewDB(cfg)
    if err != nil {
        panic(err)
    }

    // Use GORM API
    // database.AutoMigrate(&YourModel{})
    // database.Create(&YourModel{})
}
```

#### Supported Database Dialects

| Database | Dialect String | DSN Example |
|----------|----------------|-------------|
| **MySQL** | `"mysql"` | `user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local` |
| **PostgreSQL** | `"postgresql"` | `user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai` |
| **SQL Server** | `"sqlserver"` | `sqlserver://gorm:password@localhost:9930?database=gorm` |
| **SQLite** | `"sqlite"` | `file.db?_busy_timeout=5000` |

### High Availability

```go
package main

import (
    "github.com/soderasen-au/go-ipc/ha"
    "time"
)

func main() {
    // Configure HA monitoring
    cfg := ha.Config{
        PeerEndpoint:   "https://peer-server.example.com/health",
        SkipVerifyCert: false,
        PeriodInSec:    30,    // Check every 30 seconds
        TimeoutInMs:    5000,  // 5 second timeout
        Retries:        3,     // Retry 3 times on failure
    }

    // Create HA agent
    agent, res := ha.NewAgent(cfg)
    if res != nil {
        panic(res)
    }

    // Start monitoring
    agent.Start()
    defer agent.Stop()

    // Check current mode
    for {
        if agent.Mode == ha.SERVING {
            // This instance is active - handle requests
            println("I am SERVING")
        } else {
            // This instance is passive - wait
            println("I am WAITING")
        }
        time.Sleep(10 * time.Second)
    }
}
```

## Package Details

### `msn` Package Structure

```
msn/
├── msn.go       - Agent orchestration and Sender interface
├── email.go     - SMTP mailer implementation
├── email_test.go - Email configuration tests
└── msn_test.go  - Integration tests
```

**Key Types:**
- `Agent` - Manages multiple message senders with thread-safe operations
- `Mailer` - SMTP email implementation
- `Message` - Email message structure
- `EmailServerConfig` - SMTP server configuration

### `db` Package Structure

```
db/
├── config.go        - Database factory and configuration
└── dialect/
    ├── dialect.go   - Dialect interface
    ├── mysql.go     - MySQL driver
    ├── postgres.go  - PostgreSQL driver
    ├── sqlite.go    - SQLite driver
    └── sqlserver.go - SQL Server driver
```

### `ha` Package Structure

```
ha/
└── client.go - HA agent with health checking logic
```

**Key Types:**
- `Agent` - HA monitoring agent
- `Config` - HA configuration
- `RunMode` - Operating mode (WAITING/SERVING)
- `Response` - Health check response

## Dependencies

- [GORM](https://gorm.io/) - Database ORM (v1.25.11)
- [go-simple-mail](https://github.com/xhit/go-simple-mail) - SMTP client (v2.16.0)
- [zerolog](https://github.com/rs/zerolog) - Structured logging (v1.33.0)
- [go-common](https://github.com/soderasen-au/go-common) - Common utilities (v0.6.3)

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./msn
```

### Building

```bash
# Build the module
go build ./...

# Verify dependencies
go mod verify

# Update dependencies
go get -u ./...
go mod tidy
```

## Architecture

### Design Patterns

- **Factory Pattern**: `NewDB()`, `NewAgent()`, `NewMailer()` for creating instances
- **Strategy Pattern**: `Sender` interface allows multiple message delivery implementations
- **Options Pattern**: `MailerOption` for flexible configuration
- **Concurrent Design**: Goroutine-based async operations with channel-based IPC

### Thread Safety

- The `msn.Agent` uses mutex locks for thread-safe concurrent access
- Background monitoring in `ha.Agent` uses channels for graceful shutdown

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Soderasen AU**

Copyright (c) 2023 Soderasen AU

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please file an issue on the GitHub repository.
