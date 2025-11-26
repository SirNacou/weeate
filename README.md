# Weeate

A full-stack web application for collaborative food ordering and decision-making. Weeate helps groups decide what to eat through polls and streamlines the ordering process.

## Features

- **Food Management**: Browse and manage food listings with images and descriptions
- **Polls**: Create polls to help groups decide on food choices
- **Orders**: Manage and track food orders
- **Real-time Updates**: Live updates via WebSocket connections using Centrifugo
- **Authentication**: Secure user authentication powered by Supabase

## Tech Stack

### Backend

- **Language**: Go 1.25+
- **Framework**: [Fiber](https://gofiber.io/) with [Huma](https://huma.rocks/) for API documentation
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Authentication**: [Supabase Auth](https://supabase.com/auth) with JWT validation
- **Real-time**: [Centrifugo](https://centrifugal.dev/) for WebSocket connections
- **Message Bus**: [Watermill](https://watermill.io/) for event-driven architecture
- **Architecture**: Clean Architecture with domain-driven design

### Frontend

- **Framework**: React 19 with TypeScript
- **Build Tool**: [Vite](https://vitejs.dev/)
- **Routing**: [TanStack Router](https://tanstack.com/router)
- **Data Fetching**: [TanStack Query](https://tanstack.com/query)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) with [Shadcn UI](https://ui.shadcn.com/)
- **Forms**: [React Hook Form](https://react-hook-form.com/) with [Zod](https://zod.dev/) validation
- **Linting**: [Biome](https://biomejs.dev/)

### Infrastructure

- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL 18
- **Real-time Server**: Centrifugo
- **Tunneling**: Cloudflare Tunnel (optional)
- **Image Storage**: ImageKit

## Project Structure

```
weeate/
├── backend/
│   ├── cmd/                    # Application entry points
│   │   ├── main.go            # Main application setup
│   │   └── server.go          # HTTP server configuration
│   └── internal/
│       ├── common/            # Shared infrastructure
│       │   ├── api/           # Common API utilities
│       │   ├── domain/        # Shared domain types
│       │   ├── events/        # Event definitions
│       │   └── infrastructure/# Database, bus, config
│       └── features/          # Feature modules
│           ├── auth/          # Authentication
│           ├── foods/         # Food management
│           ├── orders/        # Order management
│           └── polls/         # Polling system
├── frontend/
│   └── src/
│       ├── api/               # API client
│       ├── components/        # UI components
│       ├── features/          # Feature-specific code
│       ├── hooks/             # Custom React hooks
│       ├── routes/            # TanStack Router routes
│       └── lib/               # Utility functions
├── centrifugo/                # Centrifugo configuration
├── docker-compose.yml         # Production compose file
├── docker-compose.dev.yml     # Development compose file
└── openapi.yaml               # API specification
```

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- [Go 1.25+](https://golang.org/dl/) (for local backend development)
- [Bun](https://bun.sh/) (for local frontend development)
- [Just](https://github.com/casey/just) (optional, for task automation)

### Environment Setup

1. Copy the example environment file:

   ```bash
   cp .env.example .env
   ```

2. Configure the required environment variables in `.env`:

   - **Supabase**: Set up a project at [supabase.com](https://supabase.com) and add credentials
   - **Database**: Configure PostgreSQL connection details
   - **Centrifugo**: Set HMAC secret key for real-time authentication
   - **ImageKit**: Configure for image storage (optional)

### Running with Docker

Start all services:

```bash
docker compose up --build
```

For development with hot reload:

```bash
just docker-dev-up
# or
docker compose -f docker-compose.dev.yml up --build
```

The application will be available at:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080

### Local Development

#### Backend

```bash
# Run the backend server
go run ./backend/cmd/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
bun install

# Start development server
bun run dev
```

#### Frontend Commands

```bash
bun run build      # Build for production
bun run test       # Run tests
bun run lint       # Run linter
bun run format     # Format code
bun run storybook  # Start Storybook
```

## API Documentation

The API follows OpenAPI 3.1 specification. See [`openapi.yaml`](./openapi.yaml) for the full specification.

### Main Endpoints

- `GET /` - Get current user information
- `GET /foods/` - List all foods
- `POST /foods/` - Add a new food item

## Architecture

The backend follows Clean Architecture principles with strict layer separation:

- **Domain Layer**: Core business entities and interfaces
- **Application Layer**: Business logic and command handlers
- **API Layer**: HTTP endpoints and request/response handling
- **Infrastructure Layer**: External dependencies (database, messaging, etc.)

Dependencies flow inward only, ensuring the domain layer remains independent of external concerns.

## License

This project is licensed under the GNU General Public License v3.0 - see the [COPYING](COPYING) file for details.