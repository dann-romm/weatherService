FROM golang:1.18-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /app

RUN go build -o weatherBackend cmd/weatherBackend/main.go

EXPOSE 8900
EXPOSE 80

CMD ["./weatherBackend"]
