FROM golang:1.19-bullseye AS builder

RUN apt update
RUN apt dist-upgrade -y

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/main cmd/worker/main.go

FROM gcr.io/distroless/static-debian11

WORKDIR /app

COPY --from=builder /app/main /

ENTRYPOINT ["/main"]