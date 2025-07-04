# Use the official Golang image
FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o main .

EXPOSE 8081

CMD ["./main"]
