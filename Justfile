set windows-shell := ["powershell", "-Command"]
set shell := ["bash", "-c"]

# Default task: list all available tasks
default:
  @just -l
  
# Start development environment with Docker Compose
docker-dev-up *ARGS:
  @docker compose -f docker-compose.dev.yml down {{ARGS}}
  docker compose -f docker-compose.dev.yml up --build -d {{ARGS}}

docker-dev-up-force *ARGS:
  @docker compose -f docker-compose.dev.yml down {{ARGS}}
  docker compose -f docker-compose.dev.yml up --build -d --force-recreate {{ARGS}}

docker-dev-restart *ARGS:
  docker compose -f docker-compose.dev.yml restart {{ARGS}}

# Stop development environment
docker-dev-down *ARGS:
  docker compose -f docker-compose.dev.yml down {{ARGS}}

docker-up *ARGS:
  docker compose up --build -d {{ARGS}}

docker-down *ARGS:
  docker compose down {{ARGS}}

docker-restart *ARGS:
  docker compose restart {{ARGS}}

go-integration-test *ARGS:
  go test -tags=integration ./backend/tests/integration/... {{ARGS}}