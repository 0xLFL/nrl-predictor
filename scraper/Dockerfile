# ----------- Build Stage -----------
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scraper .

# ----------- Runtime Stage -----------
FROM debian:bullseye-slim

# Install Chromium and dependencies
RUN apt-get update && apt-get install -y \
    chromium \
    fonts-liberation \
    libnss3 \
    libxss1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libgtk-3-0 \
    ca-certificates \
    wget \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy the built scraper binary from builder stage
COPY --from=builder /app/scraper .

# Tell chromedp where Chromium is
ENV CHROMEDP_CHROME_PATH=/usr/bin/chromium

# Default command
CMD ["./scraper"]
