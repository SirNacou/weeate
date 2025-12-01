FROM oven/bun:1 AS builder

WORKDIR /app

COPY ./frontend/package.json ./frontend/bun.lock ./
RUN bun install --frozen-lockfile

COPY ./.env ./.env
COPY ./frontend ./

RUN bun run build

FROM oven/bun:1-alpine

WORKDIR /app
COPY --from=builder /app/package.json ./
COPY --from=builder /app/.output ./.output

EXPOSE 3000
CMD [ "bun", "run", "start" ]