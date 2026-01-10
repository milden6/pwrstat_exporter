# Build the application from source
FROM golang:1.25.5-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY collector/ ./collector/
COPY pwrstat/ ./pwrstat/
COPY server/ ./server/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /pwrstat_exporter 

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:3.23.2 AS build-release-stage

WORKDIR /

# Deploy the application binary into a lean image
COPY --from=build-stage /pwrstat_exporter /pwrstat_exporter

RUN addgroup -g 1000 -S pwrstat_user && \
    adduser -u 1000 -S pwrstat_user -G pwrstat_user

USER pwrstat_user

EXPOSE 9101

ENTRYPOINT ["/pwrstat_exporter"]