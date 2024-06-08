FROM golang:1.22 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go mod tidy
RUN go build -v -o /usr/local/bin/app ./cmd/server/main.go

# Run
FROM alpine
COPY --from=build ./usr/local/bin/app ./bin/app
RUN mkdir config
COPY ./config/config.yaml ./config/config.yaml

CMD ["app"]