# Chem-Factory

A Go-based REST API for a chemical factory simulation game. Build your chemical empire by mixing materials, managing inventory, and trading on the market.

## Features

- **User Authentication** - JWT-based authentication with registration and login
- **User Profiles** - Track balance, XP, and level progression
- **Material Mixing** - Combine materials to create new compounds
- **Inventory Management** - Track owned materials and quantities
- **Marketplace** - Buy and sell materials with other players
- **Material System** - Pre-defined materials with recipes and mix times
- **GORM & SQLite** - Robust ORM-backed SQLite storage with separate migration and seeding commands

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.25+ |
| Framework | Gin (HTTP) |
| ORM | GORM (gorm.io/gorm) |
| Database | SQLite (gorm.io/driver/sqlite) |
| Auth | JWT (golang-jwt/jwt/v4) |
| Config | godotenv |
| Security | bcrypt (golang.org/x/crypto) |

## Project Structure

```
chem-factory/
├── cmd/
│   ├── server/main.go      # HTTP server entry point
│   ├── migration/main.go   # Database migration entry point
│   └── seed/main.go        # Database seeding entry point
├── internal/
│   ├── database/sqlite/    # SQLite connection & GORM setup
│   ├── domain/             # Global domain models & interfaces
│   ├── modules/            # Feature modules (clean architecture with GORM private domains)
│   │   ├── auth/           # Authentication module
│   │   ├── user/           # User management
│   │   ├── inventory/      # Inventory management
│   │   ├── market/         # Marketplace
│   │   ├── mixer/          # Material mixing
│   │   └── material/       # Material definitions
│   └── routes/http/        # HTTP routing & middleware
├── pkg/
│   ├── constants/          # Application constants
│   └── material/           # Default material definitions
├── utils/                  # Utilities (JWT, hash, convert)
├── api-test/               # API test files (.http)
├── api-examples/           # API usage examples
├── docs/                   # Documentation & test scenarios
└── main.go                 # Main entry point (serve/migrate/seed)
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- SQLite3

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd Chem-Factory

# Install dependencies
go mod tidy

# Copy environment configuration
cp .env.example .env
```

### Configuration

Edit `.env` with your settings:

```env
PORT=8090
SECRET_KEY=your-super-secret-jwt-key-change-in-production
```

### Running the Application

**Run database migrations (required first run):**
```bash
go run . migrate
```

**Seed initial materials and shop data:**
```bash
go run . seed
```

**Start the HTTP server:**
```bash
go run . serve
```

The server will start on `http://localhost:8090` (or your configured PORT).

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login and receive JWT |

### User
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/user/profile` | Get authenticated user profile |

### Inventory
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/inventory/items` | List user's inventory |

### Mixer
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/mixer/mix` | Mix two materials |
| GET | `/mixer/check` | Check mix result without consuming |
| GET | `/mixer/mixes` | List user's mix history |
| POST | `/mixer/pick` | Pick up completed mix |
| POST | `/mixer/new-material` | Create new material (admin) |

### Market
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/market/all-items` | List all market listings |
| POST | `/market/buy` | Buy material from market |
| POST | `/market/set-for-sell` | List material for sale |

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o chem-factory .
./chem-factory serve
```

## Database Schema

Tables created via GORM auto-migration matching exact SQL constraints:
- **users** - User accounts with balance, XP, level
- **materials** - Material definitions with recipes
- **inventory** - User-owned materials
- **market** - Marketplace listings
- **mixes** - Active/completed mixing operations

## Default Materials

Pre-defined materials are loaded from `pkg/material/default_materials.json` during the `seed` command. Base materials (no ingredients) are automatically seeded to the market shop.

## License

MIT License - feel free to use and modify for your own projects.
