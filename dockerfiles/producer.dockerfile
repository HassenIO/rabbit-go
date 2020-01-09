FROM golang:1.13.5

WORKDIR /app
COPY . .
RUN go build ./producer/main.go

EXPOSE 3000
CMD ["./main"]
