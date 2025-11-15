FROM golang:1-alpine

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./backend ./backend

RUN go build -o /usr/local/bin/app ./backend/cmd/main.go
EXPOSE 8080
CMD ["app"]