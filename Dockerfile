# build stage
FROM golang:1.19 AS builder
WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=auto
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o main main.go

# run stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
COPY *.env .
EXPOSE 8083
CMD ["/app/main"]
