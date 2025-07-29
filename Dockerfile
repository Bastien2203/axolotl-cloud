########################################################################
# Frontend Build (Node + Vite)
########################################################################
FROM --platform=$BUILDPLATFORM node:22.15-alpine AS node-builder

WORKDIR /app
COPY front/package.json ./
COPY front/package-lock.json ./
RUN npm install

ENV VITE_APP_ENV=production

COPY front/ .
RUN npm run build

########################################################################
# Backend Build (Go)
########################################################################
FROM golang:1.24-alpine AS go-builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN apk add --no-cache --virtual .build-deps \
    gcc \
    musl-dev \
    sqlite-dev \
    linux-headers \
    libc6-compat

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN echo "TARGETARCH=$TARGETARCH VARIANT=$TARGETVARIANT"
ENV CGO_ENABLED=1 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH 

RUN if [ "$TARGETARCH" = "arm" ]; then export GOARM="${TARGETVARIANT#v}"; fi && \
    go build cmd/main.go

########################################################################
# 3) Final image (Alpine)
########################################################################
FROM alpine:3.21

RUN apk add --no-cache --virtual .runtime-deps sqlite
RUN apk add --no-cache su-exec

WORKDIR /app
USER root

COPY --from=go-builder /app/main    ./api
COPY --from=node-builder /app/dist ./dist

RUN mkdir /app/volumes

VOLUME ["/app/data", "/app/volumes"]
CMD ["./api"]


