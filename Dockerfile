FROM golang:1.14-alpine as builder

RUN mkdir -p /data-receiver
ADD . /data-receiver
WORKDIR /data-receiver
RUN go build -o data-receiver .

FROM alpine:edge

USER nobody
COPY --from=builder /data-receiver/data-receiver /app/
WORKDIR /app

ENV GIN_MODE=release

EXPOSE 8080
CMD ["/app/data-receiver"]