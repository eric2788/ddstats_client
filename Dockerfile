FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN apk add git

RUN go mod download

RUN go build -o /go/bin/ddstats_client

FROM alpine:latest

COPY --from=builder /go/bin/ddstats_client /ddstats_client
RUN chmod +x /ddstats_client

ENV GIN_MODE=release

EXPOSE 9090

VOLUME /data

ENTRYPOINT [ "/ddstats_client" ]