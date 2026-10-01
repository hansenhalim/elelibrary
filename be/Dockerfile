FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /api ./cmd/api

FROM alpine:3
COPY --from=build /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
