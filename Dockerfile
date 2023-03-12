FROM golang:1.19-bullseye AS builder

RUN apt update
RUN apt dist-upgrade -y

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /worker cmd/worker/main.go

FROM gcr.io/distroless/static-debian11

WORKDIR /app/

COPY --chown=10001:10001 --from=0 /worker ./

ENTRYPOINT ["./worker"]