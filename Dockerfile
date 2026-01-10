# Build the application from source
FROM golang:1-alpine AS build-stage

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

# Deploy the application binary into a lean image
FROM ubuntu:24.04 AS build-release-stage

# Prepare system
RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y wget sed && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install pwrstat
RUN wget -O /tmp/PPL_64bit_v1.4.1.deb https://dl4jz3rbrsfum.cloudfront.net/software/PPL_64bit_v1.4.1.deb && \
    dpkg -i /tmp/PPL_64bit_v1.4.1.deb && \
    rm -rf /tmp/*

# Configure pwrstat
# ENV DELAY_BEFORE_SHUTDOWN_SEC="300" \
#     SYSTEM_SHUTDOWN="on" \
#     RUN_SCRIPT="off" \
#     SCRIPT_RUN_DURATION_SEC="0" \
#     POWERFAIL_SCRIPT_PATH="/etc/pwrstatd-powerfail.sh" \
#     LOWBATT_SCRIPT_PATH="/etc/pwrstatd-lowbatt.sh" \
#     LOWBATT_THRESHOLD="50" \
#     BATTERY_RUNTIME_THRESHOLD_SEC="300" \
#     SHUTDOWN_SUSTAIN_SEC="60"

COPY pwrstat.sh /usr/local/bin/pwrstat.sh
RUN chmod +x /usr/local/bin/pwrstat.sh && /usr/local/bin/pwrstat.sh

WORKDIR /

COPY --from=build-stage /pwrstat_exporter /pwrstat_exporter

EXPOSE 9101

ENTRYPOINT ["/pwrstat_exporter"]