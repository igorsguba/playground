FROM golang:1.14 as builder
WORKDIR /go/consumer
COPY main.go .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates && \
    addgroup -g 99 appuser && \
    adduser -D -u 99 -G appuser appuser
USER appuser
WORKDIR /app
COPY --from=builder /go/consumer/main .
CMD ["./main"]