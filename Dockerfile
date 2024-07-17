FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port (if needed)
# EXPOSE 8080

# Set the entrypoint to the compiled binary
ENTRYPOINT ["./main"]
