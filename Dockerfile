# syntax=docker/dockerfile:1

FROM golang:1.22.3

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /dockarr cmd/dockarr/main.go

# Run
CMD ["/dockarr"]