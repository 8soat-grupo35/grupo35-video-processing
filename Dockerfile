FROM golang:1.24

WORKDIR /app
COPY . /app

RUN go install github.com/air-verse/air@v1.61.1
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY go.mod go.sum ./
RUN go mod download


CMD ["air", "-c", ".air.toml"]