FROM oven/bun:1 AS builder

WORKDIR /app

COPY ./.env ./

COPY ./frontend/package.json ./frontend/bun.lock ./

RUN bun install

COPY ./frontend ./
RUN bun run build

FROM oven/bun:1-alpine

WORKDIR /app

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/.env ./

EXPOSE 3000

CMD ["bun", "run", "start"]