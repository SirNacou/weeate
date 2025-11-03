set windows-shell := ["powershell", "-Command"]
set shell := ["bash", "-c"]

default:
  
docker-dev:
  @docker-compose -f docker-compose.dev.yml up --build -d