FROM docker.io/golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/server/main.go ./
COPY pkg/ ./pkg/

RUN GOOS=linux go build -o /server

EXPOSE 8080

ENV GIN_MODE=release
CMD ["/server"]
