FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache git && go mod download

COPY . .

RUN go build -o main .

EXPOSE 8082

CMD ["./main"]
