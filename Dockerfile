# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
# This includes the main.go, config/, database/, report/ packages, and templates/template.xlsx
COPY . .

# Build the Go application
# Using CGO_ENABLED=0 and GOOS=linux for a static binary suitable for distroless images
# -ldflags="-w -s" to strip debug information and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/main .

# Stage 2: Runtime
# Using a distroless static image as the binary is built statically
FROM gcr.io/distroless/static-debian11

# Set the working directory
WORKDIR /app

# Copy the compiled application binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the template.xlsx file from the builder stage
# The path in report_handler.go is "templates/template.xlsx" relative to the workdir /app
COPY --from=builder /app/templates/template.xlsx /app/templates/template.xlsx

# Expose the port the application will run on (default 3000)
EXPOSE 3000

# Set the command to run the application
CMD ["/app/main"]
