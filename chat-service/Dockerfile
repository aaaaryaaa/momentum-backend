FROM golang:1.24-alpine

# Set environment variable to avoid CGO issues
ENV CGO_ENABLED=0 GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache git && go mod download

COPY . .

RUN go build -o main .

EXPOSE 8083


CMD ["./main"]
