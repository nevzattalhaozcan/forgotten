FROM --platform=linux/amd64 golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

FROM --platform=linux/amd64 alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/bin/server .

COPY --from=builder /app/data ./data

EXPOSE 8080

CMD ["./server"]
