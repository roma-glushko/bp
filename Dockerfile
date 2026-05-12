# syntax=docker/dockerfile:1
FROM golang:1.26.3-alpine as build

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

ENV GOOS=linux

WORKDIR /build

COPY . /build/
RUN go mod download

RUN COMMIT_SHA=${COMMIT} && \
    BUILD_DATE_FINAL=${BUILD_DATE:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")} && \
    go build -ldflags "-s -w \
    -X github.com/roma-glushko/bp/internal/version.Version=${VERSION} \
    -X github.com/roma-glushko/bp/internal/version.GitCommit=${COMMIT_SHA} \
    -X github.com/roma-glushko/bp/internal/version.BuildDate=${BUILD_DATE_FINAL}" \
    -o /build/dist/bp

FROM gcr.io/distroless/static-debian12:nonroot as release

WORKDIR /bin
COPY --from=build /build/dist/bp /bin/

ENTRYPOINT ["/bin/bp", "hello"]
