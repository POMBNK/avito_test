FROM golang:1.20-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash make gcc musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./bin/segment-service cmd/api/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/segment-service /
COPY config.yml /config.yml

EXPOSE 8080
CMD ["/segment-service"]