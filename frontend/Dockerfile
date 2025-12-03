FROM oven/bun:1 AS builder

WORKDIR /app

COPY ./.env ./

COPY ./frontend/package.json ./
COPY ./frontend/bun.lock ./

RUN bun install --frozen-lockfile

COPY ./frontend/ ./
RUN bun run build

FROM oven/bun:1-alpine

WORKDIR /app

COPY --from=builder /app/package.json ./
COPY --from=builder /app/server.ts ./
COPY --from=builder /app/dist ./dist

EXPOSE 3000

CMD ["bun", "run", "start"]