# build stage
FROM golang:1.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY *.env .
EXPOSE 8080
CMD ["/app/main"]
