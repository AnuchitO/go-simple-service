# Task:
# [ ] build docker image from scratch with go binary embedded inside it

# Step 1: Use golang image as base image
FROM golang:1.22.0-alpine as builder

# Step 2: Set working directory
WORKDIR /app

# Step 3: Copy go.mod files
COPY go.mod ./

# Step 4: Download dependencies
RUN go mod download

# Step 5: Copy source code
COPY . .

# Step 6: Build the go binary
RUN go build -o demo .

# Step 7: Use alpine image as base image
FROM alpine:3.20

# Step 8: Copy the binary from builder image to alpine image
COPY --from=builder /app/demo /app/demo

# Step 9: Run the binary
CMD ["/app/demo"]
