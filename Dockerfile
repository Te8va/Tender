FROM golang:latest AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/tender/bin/main ./cmd/tender/

FROM alpine:latest
WORKDIR /tender
RUN mkdir /tender/logs
COPY --from=build /build/cmd/tender/bin/main .
COPY --from=build /build/migrations /tender/migrations
CMD ["/tender/main"]