FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN go build -o service-b .

EXPOSE 8081

CMD ["./service-b"]