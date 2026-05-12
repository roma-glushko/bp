# syntax=docker/dockerfile:1
FROM node:22-alpine AS frontend

WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.26.3-alpine AS build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=linux

WORKDIR /build

COPY . /build/
COPY --from=frontend /frontend/dist /build/frontend/dist
RUN go mod download

RUN COMMIT_SHA=${COMMIT} && \
    BUILD_DATE_FINAL=${BUILD_DATE:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")} && \
    go build -ldflags "-s -w \
    -X github.com/roma-glushko/bp/internal/version.Version=${VERSION} \
    -X github.com/roma-glushko/bp/internal/version.GitCommit=${COMMIT_SHA} \
    -X github.com/roma-glushko/bp/internal/version.BuildDate=${BUILD_DATE_FINAL}" \
    -o /build/dist/bp

FROM gcr.io/distroless/static-debian12:nonroot AS release

WORKDIR /bin
COPY --from=build /build/dist/bp /bin/

ENTRYPOINT ["/bin/bp", "serve", "--no-open"]
