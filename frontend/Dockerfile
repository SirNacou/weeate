FROM oven/bun:1 AS builder

WORKDIR /app

COPY ./.env ./

WORKDIR /app/frontend

COPY ./frontend/package.json ./
COPY ./frontend/bun.lock ./

RUN bun install --frozen-lockfile

COPY ./frontend/ ./
RUN bunx dotenvx run -f ../.env -- bun run build

FROM nginx:mainline-alpine

COPY --from=builder /app/frontend/dist /usr/share/nginx/html

COPY ./frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD [ "nginx", "-g", "daemon off;" ]