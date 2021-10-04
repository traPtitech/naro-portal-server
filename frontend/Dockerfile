FROM node:16.10.0-alpine AS builder
WORKDIR /build
RUN apk add git

RUN git clone https://github.com/Q-n-A-dev-team/Q-n-A_UI
WORKDIR /build/Q-n-A_UI
RUN npm ci --unsafe-perm
RUN npm run build

FROM caddy:2.4.5-alpine AS runner
EXPOSE 80

COPY --from=builder /build/Q-n-A_UI/dist /usr/share/caddy
