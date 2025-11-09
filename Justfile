set windows-shell := ["powershell", "-Command"]
set shell := ["bash", "-c"]

# Default task: list all available tasks
default:
  @just -l
  
# Start development environment with Docker Compose
docker-dev-up *ARGS:
  docker compose -f docker-compose.dev.yml up --build -d {{ARGS}}

# Stop development environment
docker-dev-down *ARGS:
  docker compose -f docker-compose.dev.yml down {{ARGS}}