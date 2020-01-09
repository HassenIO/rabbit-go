FROM golang:1.13.5

WORKDIR /app
COPY . .
RUN go build ./consumer/main.go

CMD ["./main"]
