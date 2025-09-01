# Step 1: Build the binary with Go 1.24
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o task-binary main.go

# Step 2: Minimal runtime image
FROM scratch
COPY --from=builder /app/task-binary /task-binary
CMD ["/task-binary"]
