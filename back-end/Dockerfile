FROM golang:latest AS builder

ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE="on"
WORKDIR /go/src/schedule
COPY . .
RUN CGO_ENABLED=0 go build -a -o main .

FROM alpine:latest
WORKDIR /go/src/schedule
COPY --from=builder /go/src/schedule/ .
EXPOSE 8080
CMD ["/go/src/schedule/main"]