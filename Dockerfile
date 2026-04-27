FROM --platform=linux/amd64 golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /api ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /healthcheck ./cmd/healthcheck

FROM --platform=linux/amd64 scratch
COPY --from=build /api /api
COPY --from=build /healthcheck /healthcheck
ENTRYPOINT ["/api"]
