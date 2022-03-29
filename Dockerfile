FROM golang:1.18-alpine

WORKDIR /app

COPY ./src/ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main

FROM scratch

WORKDIR /root/

COPY --from=0 /main ./

CMD ["./main"]