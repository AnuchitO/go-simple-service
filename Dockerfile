# Step 1: Use golang image as base image
FROM golang:1.22.0-alpine as builder

# Step 2: Set working directory
WORKDIR /app
ARG GIT_COMMIT
ARG VERSION

# Step 3: Copy go.mod files
COPY go.mod ./

# Step 4: Download dependencies
RUN go mod download

# Step 5: Copy source code
COPY . .

# Step 6: Build the go binary with ldflags to build a static binary with git commit hash
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.commit=$GIT_COMMIT -X main.version=$VERSION" -o demo .

# Step 7: Use alpine image as base image
#FROM scratch
FROM alpine:3.14

# Step 8: Copy the binary from builder image to scratch image
COPY --from=builder /app/demo /app/demo

# Step 9: Run the binary
CMD ["/app/demo"]
