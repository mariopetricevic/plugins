


# Build stage
FROM golang:1.19-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app/go-sample-app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/plugmex .

# Runtime stage
FROM alpine:latest

# Install ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the build stage
COPY --from=builder /app/go-sample-app/out/plugmex /app/plugmex

# Set the working directory
WORKDIR /app

# Run the binary program
CMD ["./plugmex"]





