FROM node:22-alpine AS builder

WORKDIR /app

COPY ./.env ./

WORKDIR /app/frontend

COPY ./frontend/package.json ./
COPY ./frontend/package-lock.json ./

RUN npm ci

COPY ./frontend/ ./
RUN npx dotenvx run -f ../.env -- npm run build

FROM nginx:mainline-alpine

COPY --from=builder /app/frontend/dist /usr/share/nginx/html

COPY ./frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD [ "nginx", "-g", "daemon off;" ]