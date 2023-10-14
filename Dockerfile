FROM docker.io/golang:1.21 as build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/server/main.go ./
COPY pkg/ ./pkg/

RUN GOOS=linux go build -o /server

FROM docker.io/debian:bookworm-slim
RUN apt-get update && apt-get -y install ca-certificates
WORKDIR /
COPY ./static/ ./static/
COPY --from=build /server /server

EXPOSE 8080
RUN mkdir -p /data

ENV GIN_MODE=release
CMD ["/server"]
