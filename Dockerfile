# build stage
# FROM golang:1.19 AS builder
# WORKDIR /app
# COPY . .
# RUN go env -w GO111MODULE=auto
# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go build -o main main.go

# # run stage
# FROM scratch
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY *.env .
# EXPOSE 8083
# CMD ["./main"]

# Stage 1: Build the Go Application
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .
RUN go env -w GO111MODULE=auto
RUN go env -w GOPROXY=https://goproxy.cn,direct
# Build the Go application (you can replace "main.go" with your main source file)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# Stage 2: Create a lightweight container to run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

COPY *.env .
COPY db/migration ./db/migration
# Expose any necessary ports (replace 8080 with the port your app listens on)
EXPOSE 8080

# Define the command to run the application
CMD ["./main"]
